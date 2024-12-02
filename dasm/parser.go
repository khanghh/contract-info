package dasm

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"golang.org/x/exp/maps"
)

func ParseEventTopics(bytecode []byte) []string {
	topics := make(map[string]bool)
	it := NewInstructionIterator(bytecode)
	findPush32 := func(ins []instruction, maxBackward int) *instruction {
		maxLoop := minInt(len(ins), maxBackward)
		for i := 1; i <= maxLoop; i++ {
			item := ins[len(ins)-i]
			if item.op == vm.PUSH32 {
				return &item
			}
		}
		return nil
	}
	matchFn := opAnyOf(vm.LOG0, vm.LOG1, vm.LOG2, vm.LOG3, vm.LOG4)
	for it.Next() {
		if matchFn(it.Instruction()) {
			push32Ins := findPush32(it.ins, 50)
			if push32Ins != nil {
				topic := hex.EncodeToString(common.BytesToHash(push32Ins.arg).Bytes())
				topics[topic] = true
			}
		}
	}
	return maps.Keys(topics)
}

func ParseFunctionSelectors(bytecode []byte) []string {
	methodSigs := make(map[string]bool)
	it := NewInstructionIterator(bytecode)
	for it.Next() {
		ins := it.Instructions(5)
		if matchFuncSelector(ins) {
			buf := make([]byte, 4)
			if ins[1].op == vm.PUSH3 {
				copy(buf[1:], ins[1].arg)
			} else if ins[1].op == vm.PUSH4 {
				copy(buf, ins[1].arg)
			}
			methodSigs[hex.EncodeToString(buf)] = true
		}
	}
	return maps.Keys(methodSigs)
}

func GetMethodSigsByID(methodID string, interfaces []Interface) []string {
	ret := make(map[string]bool, 0)
	for _, intf := range interfaces {
		if elem, ok := intf.Elements[methodID]; ok {
			ret[elem.Identifier()] = true
		}
	}
	return maps.Keys(ret)
}
