package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// isSubset check if `sset` is a subset of `set`
func isSubset(sset, set []string) bool {
	checkMap := make(map[string]bool)
	for _, id := range set {
		checkMap[id] = true
	}
	for _, id := range sset {
		if !checkMap[id] {
			return false
		}
	}
	return true
}

func IsERC20(bytecode []byte) bool {
	erc20Sigs := []string{
		"dd62ed3e", //allowance(address,address)
		"095ea7b3", //approve(address,uint256)
		"70a08231", //balanceOf(address)
		"18160ddd", //totalSupply()
		"a9059cbb", //transfer(address,uint256)
		"23b872dd", //transferFrom(address,address,uint256)
	}
	return isSubset(erc20Sigs, ExtractMethodIds(bytecode))
}

func IsERC721(bytecode []byte) bool {
	erc721Sigs := []string{
		"095ea7b3", // "approve(address,uint256)": "095ea7b3",
		"70a08231", // "balanceOf(address)": "70a08231",
		"081812fc", // "getApproved(uint256)": "081812fc",
		"e985e9c5", // "isApprovedForAll(address,address)": "e985e9c5",
		"6352211e", // "ownerOf(uint256)": "6352211e",
		"42842e0e", // "safeTransferFrom(address,address,uint256)": "42842e0e",
		"b88d4fde", // "safeTransferFrom(address,address,uint256,bytes)": "b88d4fde",
		"a22cb465", // "setApprovalForAll(address,bool)": "a22cb465",
		"01ffc9a7", // "supportsInterface(bytes4)": "01ffc9a7",
		"23b872dd", // "transferFrom(address,address,uint256)": "23b872dd"
	}
	return isSubset(erc721Sigs, ExtractMethodIds(bytecode))
}

func LoadInterfaces(abiDir string) ([]Interface, error) {
	interfaces := make([]Interface, 0)
	entries, err := os.ReadDir(abiDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		data, err := os.ReadFile(path.Join(abiDir, entry.Name()))
		if err != nil {
			return nil, err
		}
		elems := make([]ABIElement, 0)
		if err := json.Unmarshal(data, &elems); err != nil {
			return nil, err
		}

		fileName := filepath.Base(entry.Name())
		ifaceName := fileName[:len(fileName)-len(filepath.Ext(fileName))]
		iface, err := NewInterface(ifaceName, elems)
		if err != nil {
			return nil, fmt.Errorf("invalid contract interface abi %s", ifaceName)
		}
		interfaces = append(interfaces, iface)
	}
	return interfaces, nil
}

func IsImplement(inft Interface, methodIds []string) bool {
	checkMap := make(map[string]bool)
	for _, id := range methodIds {
		checkMap[id] = true
	}
	for _, id := range inft.MethodIDs() {
		if !checkMap[id] {
			return false
		}
	}
	return true
}
