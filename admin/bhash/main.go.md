# Admin/Hash

	<<#-->>
	package main

	import (
		"fmt"
		"os"

		"github.com/tokenshift/blob/admin"
	)

Utility application tot hash an input parameter, using whatever hashing
algorithm and format the admin interface expects.

	func main() {
		fmt.Println(admin.Hash(os.Args[1]))
	}
