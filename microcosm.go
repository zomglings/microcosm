package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nkashy1/microcosm/accounts"
	"github.com/nkashy1/microcosm/genesis"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <subcommand> [arguments]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Subcommands: accounts, addresses, genesis")
		fmt.Fprintln(os.Stderr, "For information on any <subcommand>:")
		fmt.Fprintf(os.Stderr, "\t%s <subcommand> -h\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	var keystore, password, genesisFile string
	var numAccounts uint
	var chainID, difficulty, balance int64
	var gasLimit uint64

	createAccountFlags := flag.NewFlagSet("accounts", flag.ExitOnError)
	createAccountFlags.StringVar(&keystore, "keystore", "./", "Directory in which to store key file")
	createAccountFlags.StringVar(&password, "password", "microcosm", "Password with which to encrypt the key")
	createAccountFlags.UintVar(&numAccounts, "numAccounts", 1, "Number of accounts to create")

	getAddressesFlags := flag.NewFlagSet("addresses", flag.ExitOnError)
	getAddressesFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s addresses [keyfiles...]\n", os.Args[0])
		getAddressesFlags.PrintDefaults()
	}

	genesisFlags := flag.NewFlagSet("genesis", flag.ExitOnError)
	genesisFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s, genesis -genesisFile <output path> -chainID <chain id> -difficulty <difficulty> -balance <balance> -gasLimit <gas limit> [addresses...]\n", os.Args[0])
		genesisFlags.PrintDefaults()
	}
	genesisFlags.StringVar(&genesisFile, "genesisFile", "genesis.json", "Path at which to output genesis JSON file")
	genesisFlags.Int64Var(&chainID, "chainID", 7001337, "Chain ID for private net")
	genesisFlags.Int64Var(&difficulty, "difficulty", 100, "Chain difficulty")
	genesisFlags.Int64Var(&balance, "balance", 1000000000000000000, "Balance that each of the specified accounts should start with")
	genesisFlags.Uint64Var(&gasLimit, "gasLimit", 10000000, "Gas limit for private net")

	subcommand := flag.Arg(0)
	subcommandArgs := flag.Args()[1:]
	switch subcommand {
	case "accounts":
		createAccountFlags.Parse(subcommandArgs)
		addresses, err := accounts.CreateKeys(keystore, password, numAccounts)
		if err != nil {
			log.Fatal(err)
		}
		for _, address := range addresses {
			fmt.Printf("%s ", address.String())
		}
	case "addresses":
		getAddressesFlags.Parse(subcommandArgs)
		// If filenames have not been provided directly, check to see if they are being piped in
		var keyFiles []string
		keyFiles = getAddressesFlags.Args()
		if len(keyFiles) == 0 {
			reader := bufio.NewReader(os.Stdin)
			pipeInput, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			keyFiles = strings.Fields(pipeInput)
		}

		for _, keyFile := range keyFiles {
			address, err := accounts.GetAddress(keyFile)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s ", address.String())
		}
	case "genesis":
		genesisFlags.Parse(subcommandArgs)
		var rawAddresses []string
		rawAddresses = genesisFlags.Args()
		// If addresses have not been provided directly, check to see if they are being piped in
		if len(rawAddresses) == 0 {
			reader := bufio.NewReader(os.Stdin)
			pipeInput, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			rawAddresses = strings.Fields(pipeInput)
		}

		addresses := make([]common.Address, len(rawAddresses))
		for i := 0; i < len(rawAddresses); i++ {
			addresses[i] = common.BytesToAddress([]byte(common.FromHex(rawAddresses[i])))
		}

		genesisBlock, err := genesis.CreateGenesisBlock(
			chainID,
			difficulty,
			balance,
			gasLimit,
			addresses,
		)
		if err != nil {
			log.Fatal(err)
		}

		genesisBlockJSON, err := json.Marshal(genesisBlock)
		if err != nil {
			log.Fatal(err)
		}

		ioutil.WriteFile(genesisFile, genesisBlockJSON, 0644)
	default:
		fmt.Fprintf(os.Stderr, "Invalid subcommand: %s\n", subcommand)
		flag.Usage()
		os.Exit(1)
	}
}
