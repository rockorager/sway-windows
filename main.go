package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joshuarubin/go-sway"
)

type PluginSearchResult struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Exec        *string  `json:"exec"`
	Window      []int    `json:"window"`
}

type Append struct {
	Append PluginSearchResult
}

func main() {
	buf := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(buf)

	ctx := context.Background()
	client, err := sway.New(ctx)
	if err != nil {
		log.Printf("error: %s\n", err)
		os.Exit(1)
	}

	root, err := client.GetTree(ctx)
	if err != nil {
		log.Printf("error: %s\n", err)
		os.Exit(1)
	}

	cancel := context.CancelFunc(func() {})
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		var payload interface{}
		err := json.Unmarshal([]byte(line), &payload)
		if err != nil {
			log.Printf("error: %s\n", err)
			os.Exit(1)
		}
		switch v := payload.(type) {
		case string:
			switch v {
			case "Interrupt":
				cancel()
			case "Exit":
				return
			}
		case map[string]interface{}:
			for k, v := range v {
				switch k {
				case "Activate":
					cmd := fmt.Sprintf("[con_id=%v] focus", v)
					reply, err := client.RunCommand(context.Background(), cmd)
					if err != nil {
						log.Printf("error: %s", err)
					}
					log.Println(reply)
					fmt.Println("\"Close\"")
					return
				case "ActivateContext":
				case "Complete":
				case "Context":
				case "Quit":
				case "Search":
					ctx, cancel = context.WithCancel(context.Background())
					go search(ctx, root, v.(string))
				}
			}

		}
	}
}

func search(ctx context.Context, root *sway.Node, args string) {
	walk(root, ctx, args)
	fmt.Println("\"Finished\"")
}

func walk(node *sway.Node, ctx context.Context, args string) {
	if ctx.Err() != nil {
		return
	}
	switch node.Type {
	case sway.NodeRoot, sway.NodeOutput, sway.NodeWorkspace:
	case sway.NodeCon:
		if node.PID != nil {
			// This is an actual window
			result := PluginSearchResult{
				Id:          node.ID,
				Name:        "󱂬  " + node.Name,
				Description: "",
				Keywords:    []string{"sway"},
				Exec:        nil,
				Window:      nil,
			}

			if node.AppID != nil {
				result.Keywords = append(result.Keywords, *node.AppID)
				result.Description = *node.AppID
			}
			response := Append{Append: result}
			b, err := json.Marshal(response)
			if err != nil {
				log.Printf("error: %s\n", err)
				return
			}
			fmt.Println(string(b))
		}
	case sway.NodeFloatingCon:
	}
	for _, child := range node.Nodes {
		walk(child, ctx, args)
	}
}
