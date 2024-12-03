# contract-info
A command-line tool to help identify the implemented methods, events and interfaces of a contract on EVM blockchain

## Usage
```bash
NAME:
   impl - Ethereum contract parser a1f9f966200f91aec3fdd58d796f42c58c3f89a4

USAGE:
   impl [global options] command [command options]

VERSION:
   a1f9f966200f91aec3fdd58d796f42c58c3f89a4 - 2024-11-28T17:57:43 

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --rpcurl value     ethereum JSON-RPC URLs to fetch the blockchain data [$DASM_RPC_URL]
   --abis value       ABIs directory to load the contract interfaces (default: "abis")
   --verbosity value  Log verbosity level (0-5) (default: 3) [$VERBOSITY]
   --help, -h         show help
   --version, -v      print the version
```

Example: 
```bash
$ ./impl --rpcurl=https://ethereum-rpc.publicnode.com 0xdac17f958d2ee523a2206206994597c13d831ec7
Loaded 5 interface ABIs
Fetching contract bytecode...
Contract information:
Address             0xdAC17F958D2ee523a2206206994597C13D831ec7
Is Proxy            false
Poissible Methods   - 313ce567
                    - 3f4ba83a
                    - 8da5cb5b
                    - 23b872dd transferFrom(address,address,uint256)
                    - dd62ed3e allowance(address,address)
                    - 26976e3f
                    - 27e235e3
                    - 893d20e8
                    - dd644f72
                    - 8456cb59
                    - 5c658165
                    - e4997dc5
                    - c0324c77
                    - f2fde38b
                    - a9059cbb transfer(address,uint256)
                    - e47d6060
                    - 0e136b19
                    - 095ea7b3 approve(address,uint256)
                    - 18160ddd totalSupply()
                    - cc872b66
                    - 0ecb93c0
                    - 0753c30c
                    - e5b5019a
                    - 06fdde03 name()
                    - f3bdc228
                    - 5c975abb
                    - 3eaaf86b
                    - 35390714
                    - 95d89b41 symbol()
                    - db006a75
                    - 59bf1abe
                    - 70a08231 balanceOf(address)
Poissible Events    ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
                    7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33
                    6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625
                    cb8241adb0c3fdb35b70c24ce35c5eb0c17af7431c99f827d44a445ca624176a
                    702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44
                    8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925
Possible Interfaces ERC20
```
```bash
$ ./impl 0x411D79b8cC43384FDE66CaBf9b6a17180c842511
Loaded 7 interface ABIs
Fetching contract bytecode...
Contract information:
Address                0x411D79b8cC43384FDE66CaBf9b6a17180c842511                       
Is Proxy Contract      true                                                             
Implementation Address 0x7a68e572eFE159753813eB86A8c84157d684bda2                       
Poissible Methods      - 5c60da1b implementation()                                      
                       - 8f283970 changeAdmin(address)                                  
                       - f851a440 admin()                                               
                       - 3659cfe6 upgradeTo(address)                                    
                       - 4f1ef286 upgradeToAndCall(address,bytes)                       
Poissible Events       7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f 
                       bc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b 
Possible Interfaces    - BaseAdminUpgradeabilityProxy                                   
                       - BaseUpgradeabilityProxy 
```
## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.