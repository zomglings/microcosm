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
		fmt.Fprintln(os.Stderr, "Subcommands: account")
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

	createAccountFlags := flag.NewFlagSet("account", flag.ExitOnError)
	createAccountFlags.StringVar(&keystore, "keystore", "./", "Directory in which to store key file")
	createAccountFlags.StringVar(&password, "password", "microcosm", "Password with which to encrypt the key")

	subcommand := flag.Arg(0)
	switch subcommand {
	case "account":
		createAccountFlags.Parse(flag.Args()[1:])
		address, err := accounts.CreateKey(keystore, password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created keyfile for address: %s\n", address.String())
	default:
		fmt.Fprintf(os.Stderr, "Invalid subcommand: %s\n", subcommand)
		flag.Usage()
		os.Exit(1)
	}
}
