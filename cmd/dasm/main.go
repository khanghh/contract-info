package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/khanghh/contract-info/dasm"
	"github.com/urfave/cli/v2"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	// The date of the release (set via linker flags)
	gitDate = ""
	// The version of the release (set via linker flags)
	gitTag = ""
	// The app that holds all commands and flags.
	app *cli.App
)

var (
	rpcUrlFlag = &cli.StringFlag{
		Name:     "rpcurl",
		Required: true,
		EnvVars:  []string{"DASM_RPC_URL"},
		Usage:    "ethereum JSON-RPC URLs to fetch the blockchain data",
	}
	outputDirFlag = &cli.StringFlag{
		Name:    "outdir",
		Aliases: []string{"o"},
		Usage:   "Ouput directory to save disassembled bytecode",
	}
	verbosityFlag = &cli.IntFlag{
		Name:    "verbosity",
		Usage:   "Log verbosity level (0-5)",
		Value:   3,
		EnvVars: []string{"VERBOSITY"},
	}
)

func init() {
	app = cli.NewApp()
	app.Action = run
	app.Name = filepath.Base(os.Args[0])
	app.Usage = fmt.Sprintf("Ethereum bytecode disassembler %s", gitTag)
	app.Version = fmt.Sprintf("%s - %s ", gitCommit, gitDate)
	app.Flags = []cli.Flag{
		rpcUrlFlag,
		outputDirFlag,
		verbosityFlag,
	}
}

func mustInitRpcClient(cli *cli.Context) *rpc.Client {
	rpcUrl := cli.String(rpcUrlFlag.Name)
	client, err := rpc.Dial(rpcUrl)
	if err != nil {
		panic(fmt.Errorf("could not dial RPC: %w", err))
	}
	return client
}

func ethGetCode(client *rpc.Client, addr common.Address) ([]byte, error) {
	var result hexutil.Bytes
	err := client.Call(&result, "eth_getCode", addr, "latest")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func createDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking directory: %w", err)
	}
	return nil
}

func run(cli *cli.Context) error {
	// ethcore.InitDefaultLogger(cli.Int(verbosityFlag.Name))
	addrStr := cli.Args().Get(0)
	if addrStr == "" {
		return errors.New("must provide contract address")
	}

	client := mustInitRpcClient(cli)
	defer client.Close()

	addr := common.HexToAddress(addrStr)
	fmt.Printf("Fetching contract bytecode for address %v\n", addr)
	bytecode, err := ethGetCode(client, addr)
	if err != nil {
		return fmt.Errorf("could not get contract bytecode from rpc: %w", err)
	}

	instructions, err := dasm.Disassemble(bytecode)
	if err != nil {
		panic(fmt.Sprintf("Failed to disaasemble contract bytecode %v", err))
	}

	var dasmCode string
	for _, ins := range instructions {
		fmt.Print(ins)
		dasmCode += ins
	}

	outputDir := cli.String(outputDirFlag.Name)
	if err := createDirIfNotExist(outputDir); err != nil {
		panic(err)
	}

	dasmOutFile := path.Join(outputDir, fmt.Sprintf("%s.dasm", hexutil.Encode(addr.Bytes())))
	binOutFile := path.Join(outputDir, fmt.Sprintf("%s.bin", hexutil.Encode(addr.Bytes())))
	if err := os.WriteFile(dasmOutFile, []byte(dasmCode), 0644); err != nil {
		panic(fmt.Sprintf("Failed to write disassembled code to file: %v", err))
	}
	if err := os.WriteFile(binOutFile, []byte(bytecode), 0644); err != nil {
		panic(fmt.Sprintf("Failed to write bytecode to file: %v", err))
	}
	fmt.Printf("Disaaembled code saved to %s\n", dasmOutFile)
	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
