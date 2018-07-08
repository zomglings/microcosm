package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nkashy1/microcosm/accounts"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <subcommand> [arguments]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Subcommands: accounts")
		fmt.Fprintln(os.Stderr, "For information on any <subcommand>:")
		fmt.Fprintf(os.Stderr, "\t%s <subcommand> -h\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	var keystore, password string
	var numAccounts uint

	createAccountFlags := flag.NewFlagSet("account", flag.ExitOnError)
	createAccountFlags.StringVar(&keystore, "keystore", "./", "Directory in which to store key file")
	createAccountFlags.StringVar(&password, "password", "microcosm", "Password with which to encrypt the key")
	createAccountFlags.UintVar(&numAccounts, "numAccounts", 1, "Number of accounts to create")

	subcommand := flag.Arg(0)
	switch subcommand {
	case "accounts":
		createAccountFlags.Parse(flag.Args()[1:])
		addresses, err := accounts.CreateKeys(keystore, password, numAccounts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Created keyfiles for the following accounts:")
		for i, address := range addresses {
			fmt.Printf("\t%d. %s\n", i+1, address.String())
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid subcommand: %s\n", subcommand)
		flag.Usage()
		os.Exit(1)
	}
}
