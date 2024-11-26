package parser

import (
	"github.com/khanghh/ethcore/client"
)

type ContractParser struct {
	client     client.RemoteChainReader
	interfaces []Interface
}

func (p *ContractParser) Parse(bytecode []byte) error {
	return nil
}

func (p *ContractParser) ParseContract(bytecode []byte) (*Contract, error) {
	return nil, nil
}

func (p *ContractParser) PrintInfo(bytecode []byte) {

}

func NewByteCodeParser(client client.RemoteChainReader, interfaces []Interface) *ContractParser {
	return &ContractParser{
		client:     client,
		interfaces: interfaces,
	}
}
