package parser

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/asm"
	"github.com/ethereum/go-ethereum/core/vm"
)

type InstructionIterator interface {
	Next() bool
	Error() error
	PC() uint64
	Op() vm.OpCode
	Arg() []byte
}

type instruction struct {
	op  vm.OpCode
	arg []byte
}

func (ins instruction) OpCode() vm.OpCode {
	return ins.op
}

func (ins instruction) Arg() []byte {
	return ins.arg
}

type instructionIterator struct {
	it  InstructionIterator
	ins []instruction
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (it *instructionIterator) Instructions(n int) []instruction {
	return it.ins[len(it.ins)-minInt(n, len(it.ins)):]
}

func (it *instructionIterator) Next() bool {
	if it.it.Next() {
		it.ins = append(it.ins, instruction{it.it.Op(), it.it.Arg()})
		return true
	}
	return false
}

func (it *instructionIterator) Error() error {
	return it.it.Error()
}

func (it *instructionIterator) Op() vm.OpCode {
	return it.it.Op()
}

func (it *instructionIterator) Arg() []byte {
	return it.it.Arg()
}

func (it *instructionIterator) PC() uint64 {
	return it.it.PC()
}

func (it *instructionIterator) Instruction() *instruction {
	if len(it.ins) > 0 {
		return &it.ins[len(it.ins)-1]
	}
	return nil
}

func newInstructionIterator(bytecode []byte) *instructionIterator {
	return &instructionIterator{
		it: asm.NewInstructionIterator(bytecode),
	}
}

type matchOpCode func(op vm.OpCode) bool

func opExact(op vm.OpCode) matchOpCode {
	return func(o vm.OpCode) bool {
		return op == o
	}
}

func opAnyOf(ops ...vm.OpCode) matchOpCode {
	return func(op vm.OpCode) bool {
		for _, o := range ops {
			if op == o {
				return true
			}
		}
		return false
	}
}

func opIsPush() matchOpCode {
	return func(op vm.OpCode) bool {
		return op.IsPush()
	}
}

func matchPattern(ins []instruction, pattern []matchOpCode) bool {
	if len(ins) < len(pattern) {
		return false
	}
	for i, match := range pattern {
		if !match(ins[i].op) {
			return false
		}
	}
	return true
}

func matchFuncSelector(ins []instruction) bool {
	return matchPattern(ins, []matchOpCode{
		opExact(vm.DUP1),
		opAnyOf(vm.PUSH3, vm.PUSH4),
		opExact(vm.EQ),
		opAnyOf(vm.PUSH2, vm.PUSH3),
		opExact(vm.JUMPI),
	})
}

func matchSplitSelector(ins []instruction) bool {
	return matchPattern(ins, []matchOpCode{
		opExact(vm.DUP1),
		opAnyOf(vm.PUSH3, vm.PUSH4),
		opExact(vm.GT),
		opAnyOf(vm.PUSH2, vm.PUSH3),
		opExact(vm.JUMPI),
	})
}

func findJumptable(it *instructionIterator) bool {
	pattern := []matchOpCode{
		opExact(vm.PUSH1),
		opExact(vm.CALLDATASIZE),
		opExact(vm.LT),
		opAnyOf(vm.PUSH2, vm.PUSH3),
		opExact(vm.JUMPI),
		opExact(vm.PUSH1),
		opExact(vm.CALLDATALOAD),
		opExact(vm.PUSH1),
		opExact(vm.SHR),
	}
	for it.Next() {
		if matchPattern(it.Instructions(len(pattern)), pattern) {
			return true
		}
	}
	return false
}

func nextInstructions(it *instructionIterator, n int) []instruction {
	ins := make([]instruction, 0, n)
	for i := 0; i < n; i++ {
		if !it.Next() {
			return nil
		}
		ins = append(ins, instruction{it.Op(), it.Arg()})
	}
	return ins
}

func matchEndJumptable(ins []instruction) bool {
	return matchPattern(ins, []matchOpCode{
		opAnyOf(vm.PUSH0, vm.PUSH1),
		opExact(vm.DUP1),
		opExact(vm.REVERT),
	})
}

func matchEndFuncSelector(ins []instruction) bool {
	return matchPattern(ins, []matchOpCode{
		opExact(vm.PUSH2),
		opExact(vm.JUMP),
	})
}

// ExtractMethodIds parses the contract byte code and returns all 4-bytes method ids from jump table
// Code for selecting from n functions without split:
// SELECT[n]:
// DUP1, PUSH4 <id_i>, EQ, PUSH2/3 <tag_i>, JUMPI
// PUSH2/3 <fallback>, JUMP
//
// Code for selecting from n functions with split:
// DUP1, PUSH4 <pivot>, GT, PUSH2/3 <tag_less>, JUMPI
// -> SELECT[n/2]
// tag_less:
// -> SELECT[n/2]
func ExtractMethodIds(bytecode []byte) []string {
	it := newInstructionIterator(bytecode)
	if !findJumptable(it) {
		return nil
	}

	// The next 5 instructions must be a function selector, otherwise it is a split selector
	// so we can skip to the next 5 instructions of the split code
	methodIds := []string{}
	if matchSplitSelector(nextInstructions(it, 5)) {
		it.Next()
	}
	for {
		if matchFuncSelector(it.Instructions(5)) {
			methodIds = append(methodIds, fmt.Sprintf("%08x", it.Instructions(5)[1].arg))
		}
		if matchEndJumptable(it.Instructions(3)) {
			break
		}
		it.Next()
	}
	return methodIds
}
