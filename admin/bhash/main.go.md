# Admin/Hash

	<<#-->>
	package main

	import (
		"encoding/hex"
		"fmt"
		"os"

		"github.com/tokenshift/blob/admin"
	)

Utility application to hash an input parameter, using whatever hashing
algorithm and format the admin interface expects.

	func main() {
		var password, salt, hash []byte
		var err error

		if len(os.Args) == 2 {
			salt = admin.Salt()
		} else if len(os.Args) == 3 {
			salt, err = hex.DecodeString(os.Args[2])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Salt must be hexadecimal encoded.")
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Use: bhash password [salt]")
			os.Exit(1)
		}

		password = []byte(os.Args[1])
		hash = admin.Hash(password, salt)
		
		fmt.Printf("Salt: %x\n", salt)
		fmt.Printf("Hash: %x\n", hash)
	}
