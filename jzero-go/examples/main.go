package main

import (
	"context"
	"fmt"

	"github.com/jaronnie/jzero-go/model/pb/credentialpb"

	"github.com/jaronnie/jzero-go"
	"github.com/jaronnie/jzero-go/model/types"
	"github.com/jaronnie/jzero-go/rest"
)

func main() {
	clientset, err := jzero.NewClientWithOptions(
		rest.WithAddr("127.0.0.1"),
		rest.WithPort("8001"),
		rest.WithProtocol("http"))
	if err != nil {
		panic(err)
	}

	// api interface
	result, err := clientset.Hello().HelloPathHandler(context.Background(), &types.PathRequest{
		Name: "jzero",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Message)

	// proto gateway interface

	list, err := clientset.Credential().CredentialList(context.Background(), &credentialpb.CredentialListRequest{
		Page: 1,
		Size: 10,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
}
