package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/khanghh/contract-info/dasm"
	"github.com/olekukonko/tablewriter"
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
	app.Usage = fmt.Sprintf("Ethereum contract parser %s", gitTag)
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

func printContractInfo(data [][]string) {
	fmt.Println("Contract information:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false) // Disable table border
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(true)
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.AppendBulk(data)
	table.Render()
}

func renderMethodList(methodIDsMap map[string][]string) string {
	methodList := make([]string, 0)
	for methodID, methodSigs := range methodIDsMap {
		methodList = append(methodList, fmt.Sprintf("- %s %s", methodID, strings.Join(methodSigs, ",")))
	}
	return strings.Join(methodList, "\n")
}

func getContractInterfaces(intefaces []dasm.Interface, sigs []string) []string {
	ret := make([]string, 0)
	for _, intf := range intefaces {
		if dasm.IsImplement(intf, sigs) {
			ret = append(ret, intf.Name)
		}
	}
	return ret
}

func run(cli *cli.Context) error {
	addrStr := cli.Args().Get(0)
	if addrStr == "" {
		return errors.New("must provide contract address")
	}

	interfaces, err := dasm.LoadInterfaces(cli.String(abisDirFlag.Name))
	if err != nil {
		return fmt.Errorf("could not parse interface abi: %w", err)
	}
	fmt.Printf("Loaded %d interface ABIs\n", len(interfaces))

	client := mustInitRpcClient(cli)
	defer client.Close()

	addr := common.HexToAddress(addrStr)
	fmt.Println("Fetching contract bytecode...")
	bytecode, err := ethGetCode(client, addr)
	if err != nil {
		return fmt.Errorf("could not get contract bytecode from rpc: %w", err)
	}

	infos := make([][]string, 0)
	infos = append(infos, []string{"Address", addr.Hex()})
	infos = append(infos, []string{"Is Proxy", strconv.FormatBool(dasm.IsProxy(bytecode))})

	methodIDs := dasm.ParseFunctionSelectors(bytecode)
	// toPrintMethodSigs := make([]string, len(methodSigs))
	// for i, topic := range methodSigs {
	// 	toPrintMethodSigs[i] = "- " + topic
	// 	// sigs = append(sigs, topic)
	// }
	methodIDsMap := make(map[string][]string)
	for _, methodID := range methodIDs {
		methodIDsMap[methodID] = dasm.GetMethodSigsByID(methodID, interfaces)
	}
	infos = append(infos, []string{"Poissible Methods", renderMethodList(methodIDsMap)})

	topics := dasm.ParseEventTopics(bytecode)
	infos = append(infos, []string{"Poissible Events", strings.Join(topics, "\n")})
	contractInterfaces := getContractInterfaces(interfaces, append(methodIDs, topics...))
	if len(contractInterfaces) > 0 {
		infos = append(infos, []string{"Possible Interfaces", strings.Join(contractInterfaces, ",")})
	}

	// fmt.Println("Possible contract implentations:")
	// count := 1
	// for _, intf := range interfaces {
	// 	if dasm.IsImplement(intf, sigs) {
	// 		fmt.Printf("%d. %s\n", count, intf.Name)
	// 		for sig := range intf.Elements {
	// 			fmt.Printf("  - %s\n", sig)
	// 		}
	// 		count++
	// 	}
	// }

	printContractInfo(infos)
	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
