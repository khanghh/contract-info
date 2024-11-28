package dasm

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
)

type matcherFn func(in instruction) bool

func opExact(op vm.OpCode) matcherFn {
	return func(in instruction) bool {
		return in.op == op
	}
}

func opIsPush(hex string) matcherFn {
	return func(in instruction) bool {
		if !in.op.IsPush() {
			return false
		}
		val, _ := hexutil.Decode(hex)
		return bytes.Equal(in.arg, val)
	}
}

func opAnyOf(ops ...vm.OpCode) matcherFn {
	return func(in instruction) bool {
		for _, op := range ops {
			if in.op == op {
				return true
			}
		}
		return false
	}
}

func matchAnyOf(fns ...matcherFn) matcherFn {
	return func(in instruction) bool {
		for _, fn := range fns {
			if fn(in) {
				return true
			}
		}
		return false
	}
}

func matchPattern(ins []instruction, pattern []matcherFn) bool {
	if len(ins) < len(pattern) {
		return false
	}
	for i, match := range pattern {
		if !match(ins[i]) {
			return false
		}
	}
	return true
}

func matchFuncSelector(ins []instruction) bool {
	return matchPattern(ins, []matcherFn{
		opExact(vm.DUP1),
		opAnyOf(vm.PUSH3, vm.PUSH4),
		opAnyOf(vm.GT, vm.EQ),
		opAnyOf(vm.PUSH2, vm.PUSH3),
		opExact(vm.JUMPI),
	})
}

func matchSplitSelector(ins []instruction) bool {
	return matchPattern(ins, []matcherFn{
		opExact(vm.DUP1),
		opAnyOf(vm.PUSH3, vm.PUSH4),
		opExact(vm.GT),
		opAnyOf(vm.PUSH2, vm.PUSH3),
		opExact(vm.JUMPI),
	})
}

func findJumptable(it *instructionIterator) bool {
	pattern := []matcherFn{
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

func matchEndJumptable(ins []instruction) bool {
	return matchPattern(ins, []matcherFn{
		opAnyOf(vm.PUSH0, vm.PUSH1),
		opExact(vm.DUP1),
		opExact(vm.REVERT),
	})
}

func matchEndFuncSelector(ins []instruction) bool {
	return matchPattern(ins, []matcherFn{
		opExact(vm.PUSH2),
		opExact(vm.JUMP),
	})
}

func matchByteCode(bytecode []byte, pattern []byte) bool {
	if len(bytecode) < len(pattern) {
		return false
	}
	for i := 0; i < len(bytecode)-len(pattern); i++ {
		if bytes.Equal(bytecode[i:i+len(pattern)], pattern) {
			return true
		}
	}
	return false
}

func IsProxy(bytecode []byte) bool {
	pattern := []matcherFn{
		opExact(vm.CALLDATASIZE),
		opIsPush("0x00"),
		opExact(vm.DUP1),
		opExact(vm.CALLDATACOPY),
		opIsPush("0x00"),
		opExact(vm.DUP1),
		opExact(vm.CALLDATASIZE),
		opIsPush("0x00"),
		opExact(vm.DUP5),
		opExact(vm.GAS),
		opExact(vm.DELEGATECALL),
		opExact(vm.RETURNDATASIZE),
		opIsPush("0x00"),
		opExact(vm.DUP1),
		opExact(vm.RETURNDATACOPY),
		opExact(vm.DUP1),
		opExact(vm.DUP1),
		opExact(vm.ISZERO),
		opExact(vm.PUSH2),
		opExact(vm.JUMPI),
		opExact(vm.RETURNDATASIZE),
		opIsPush("0x00"),
		opExact(vm.RETURN),
		opExact(vm.JUMPDEST),
		opExact(vm.RETURNDATASIZE),
		opIsPush("0x00"),
		opExact(vm.REVERT),
	}

	it := NewInstructionIterator(bytecode)
	for it.Next() {
		if matchPattern(it.Instructions(len(pattern)), pattern) {
			return true
		}
	}
	return false
}

func IsImplement(inft Interface, sigs []string) bool {
	checkMap := make(map[string]bool)
	for _, id := range sigs {
		checkMap[id] = true
	}
	for id := range inft.Elements {
		if !checkMap[id] {
			return false
		}
	}
	return true
}
