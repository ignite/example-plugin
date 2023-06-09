package main

import (
	"context"
	"fmt"
	"path/filepath"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/ignite/services/chain"
	"github.com/ignite/cli/ignite/services/plugin"
)

type p struct{}

func (p) Manifest(ctx context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name: "example-plugin",
		// Add commands here
		Commands: []*plugin.Command{
			// Example of a command
			{
				Use:   "example-plugin",
				Short: "Explain what the command is doing...",
				Long:  "Long description goes here...",
				Flags: []*plugin.Flag{
					{Name: "my-flag", Type: plugin.FlagTypeString, Usage: "my flag description"},
				},
				PlaceCommandUnder: "ignite",
				// Examples of adding subcommands:
				// Commands: []*plugin.Command{
				// 	{Use: "add"},
				// 	{Use: "list"},
				// 	{Use: "delete"},
				// },
			},
		},
		// Add hooks here
		Hooks: []*plugin.Hook{},
	}, nil
}

func (p) Execute(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	// TODO: write command execution here
	fmt.Printf("Hello I'm the example-plugin plugin\n")
	fmt.Printf("My executed command: %q\n", cmd.Path)
	fmt.Printf("My args: %v\n", cmd.Args)

	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}

	myFlag, _ := flags.GetString("my-flag")
	fmt.Printf("My flags: my-flag=%q\n", myFlag)
	fmt.Printf("My config parameters: %v\n", cmd.With)

	// This is how the plugin can access the chain:
	// c, err := getChain(cmd)

	// According to the number of declared commands, you may need a switch:
	switch cmd.Use {
	case "add":
		fmt.Println("Adding stuff...")
	case "list":
		fmt.Println("Listing stuff...")
	case "delete":
		fmt.Println("Deleting stuff...")
	}
	return nil
}

func (p) ExecuteHookPre(ctx context.Context, h *plugin.ExecutedHook) error {
	fmt.Printf("Executing hook pre %q\n", h.Hook.GetName())
	return nil
}

func (p) ExecuteHookPost(ctx context.Context, h *plugin.ExecutedHook) error {
	fmt.Printf("Executing hook post %q\n", h.Hook.GetName())
	return nil
}

func (p) ExecuteHookCleanUp(ctx context.Context, h *plugin.ExecutedHook) error {
	fmt.Printf("Executing hook cleanup %q\n", h.Hook.GetName())
	return nil
}

func getChain(cmd *plugin.ExecutedCommand, chainOption ...chain.Option) (*chain.Chain, error) {
	flags, err := cmd.NewFlags()
	if err != nil {
		return nil, err
	}

	var (
		home, _ = flags.GetString("home")
		path, _ = flags.GetString("path")
	)
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return chain.New(absPath, chainOption...)
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"example-plugin": plugin.NewGRPC(&p{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
