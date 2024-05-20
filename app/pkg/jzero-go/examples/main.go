package main

import (
	"fmt"

	"github.com/jzero-io/jzero-go/examples/credential"
)

func main() {
	list, err := credential.GetCredentialList()
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
}
