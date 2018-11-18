package handler

import (
	"encoding/json"

	"golang.org/x/net/context"
	"github.com/micro/go-log"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"

	apiClient "github.com/hb-go/micro/istio/http/api/client"
	example "github.com/hb-go/micro/istio/http/srv/proto/example"
)

type Example struct{}

func extractValue(pair *api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

// Example.Call is called by the API as /http/example/call with post body {"name": "foo"}
func (e *Example) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Log("Received Example.Call request")

	// extract the client from the context
	exampleClient, ok := apiClient.ExampleFromContext(ctx)
	if !ok {
		return errors.InternalServerError("go.micro.api.sample.example.call", "example client not found")
	}

	// make request
	response, err := exampleClient.Call(ctx, &example.Request{
		Name: extractValue(req.Post["name"]),
	})
	if err != nil {
		return errors.InternalServerError("go.micro.api.sample.example.call", err.Error())
	}

	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}