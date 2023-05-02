package cmd

import (
	"log"
	"os"

	"github.com/n04ln/aoj.nvim/pkg/app/aoj_remote_plugin/handler"
	"github.com/neovim/go-client/nvim/plugin"
	"github.com/spf13/cobra"
)

var (
	remotePluginCmd = &cobra.Command{
		Use:   "remotePlugin",
		Short: "serve a RemotePlugin server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args); err != nil {
				log.Println(err)
			}
		},
	}
)

var (
	id, password string
)

func init() {
	// use stderr for logging
	log.SetOutput(os.Stderr)

	id = os.Getenv("AOJ_ID")
	if id == "" {
		log.Println("id is missing")
	}

	password = os.Getenv("AOJ_PASSWORD")
	if password == "" {
		log.Println("password is missing")
	}
}

func Execute() {
	if err := remotePluginCmd.Execute(); err != nil {
		log.Println(err)
	}
}

func run(cmd *cobra.Command, args []string) error {
	a := handler.NewAoj(id, password)
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleFunction(&plugin.FunctionOptions{Name: "AojSubmit"}, a.Submit)
		p.HandleFunction(&plugin.FunctionOptions{Name: "AojRunSample"}, a.Trial)
		p.HandleFunction(&plugin.FunctionOptions{Name: "AojDescription"}, a.Description)
		p.HandleCommand(&plugin.CommandOptions{Name: "AojSession"}, a.Session)
		return nil
	})

	return nil
}
