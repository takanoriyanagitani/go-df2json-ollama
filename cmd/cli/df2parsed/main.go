package main

import (
	"context"
	"fmt"
	"log"
	"os"

	oa "github.com/ollama/ollama/api"
	dj "github.com/takanoriyanagitani/go-df2json-ollama"
)

func df2res(
	ctx context.Context,
	c dj.Client,
	model string,
) (oa.GenerateResponse, error) {
	return c.GetParsedDfDefault(ctx, model)
}

func env2model() string {
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		return "llama3.1:8b" // default value
	}
	return model
}

func sub() error {
	cli, e := oa.ClientFromEnvironment()
	if nil != e {
		return e
	}

	var model string = env2model()

	res, e := df2res(context.Background(), dj.Client{Client: cli}, model)
	if nil != e {
		return e
	}

	var jsonRes string = dj.ResponseToJsonString(res)

	fmt.Println(jsonRes)

	return nil
}

func main() {
	e := sub()
	if nil != e {
		log.Printf("%v\n", e)
	}
}
