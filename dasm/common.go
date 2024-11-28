package dasm

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
