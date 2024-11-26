package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/khanghh/contract-info/parser"
	"github.com/khanghh/ethcore"
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
		Usage:    "ethereum JSON-RPC URLs to fetch the blockchain data",
	}
	abisDirFlag = &cli.StringFlag{
		Name:  "abis",
		Value: "abis",
		Usage: "ABIs directory to load the contract interfaces",
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
	app.Usage = fmt.Sprintf("Ethereum contract event collector %s", gitTag)
	app.Version = fmt.Sprintf("%s - %s ", gitCommit, gitDate)
	app.Flags = []cli.Flag{
		rpcUrlFlag,
		abisDirFlag,
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

func run(cli *cli.Context) error {
	ethcore.InitDefaultLogger(cli.Int(verbosityFlag.Name))
	addrStr := cli.Args().Get(0)
	if addrStr == "" {
		return errors.New("must provide contract address")
	}

	client := mustInitRpcClient(cli)
	defer client.Close()

	addr := common.HexToAddress(addrStr)
	fmt.Printf("Fetching contract bytecode %v\n", addr)
	bytecode, err := ethGetCode(client, addr)
	if err != nil {
		return fmt.Errorf("could not get contract bytecode from rpc: %w", err)
	}
	filePath := fmt.Sprintf("build/bin/%s.bin", hexutil.Encode(addr.Bytes()))
	os.WriteFile(filePath, bytecode, 0644)

	interfaces, err := parser.LoadInterfaces(cli.String(abisDirFlag.Name))
	if err != nil {
		return err
	}

	methodIds := parser.ExtractMethodIds(bytecode)
	var implementeds []string
	for _, iface := range interfaces {
		if parser.IsImplement(iface, methodIds) {
			implementeds = append(implementeds, iface.Name)
		}
	}

	fmt.Println("Contract method IDs:")
	for _, id := range methodIds {
		fmt.Println(id)
	}
	fmt.Println("Possible contract interfaces:")
	for _, impl := range implementeds {
		fmt.Println(impl)
	}
	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
