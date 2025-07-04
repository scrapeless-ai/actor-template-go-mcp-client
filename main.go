package main

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	scrapeless_actor "github.com/scrapeless-ai/sdk-go/scrapeless/actor"
	"github.com/scrapeless-ai/sdk-go/scrapeless/log"
	"time"
)

type Input struct {
	Url          string         `json:"url"`
	Actor        string         `json:"actor"`
	ProxyCountry string         `json:"proxy_country"`
	Params       map[string]any `json:"params"`
}

func main() {
	actor := scrapeless_actor.New()
	var param Input
	if err := actor.Input(&param); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Create HTTP transport
	httpClient, err := transport.NewStreamableHTTP(param.Url)
	if err != nil {
		log.Errorf("Failed to create HTTP transport: %v", err)
	}

	// Create mcp client with the transport
	mcpClient := client.NewClient(httpClient)
	defer mcpClient.Close()

	if _, err = mcpClient.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
		log.Errorf("Failed to initialize: %v", err)
		return
	}
	mcpClient.OnNotification(func(n mcp.JSONRPCNotification) {

	})

	toolResult, err := mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "deepserp",
			Arguments: map[string]any{
				"params": map[string]any{
					"q":         param.Params["q"],
					"data_type": param.Params["data_type"],
					"date":      param.Params["date"],
					"hl":        param.Params["hl"],
					"tz":        param.Params["tz"],
				},
				"actor":         param.Actor,
				"proxy_country": param.ProxyCountry,
			},
		},
	})
	if err != nil {
		log.Errorf("Failed to call tool: %v", err)
		return
	}
	log.Infof("Result: %v", toolResult.Content[0].(mcp.TextContent).Text)
}
