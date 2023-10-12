package terragruntinfo

import (
	"github.com/gads-citron/terragrunt/options"
	"github.com/gads-citron/terragrunt/pkg/cli"
)

const (
	CommandName = "terragrunt-info"
)

func NewCommand(opts *options.TerragruntOptions) *cli.Command {
	return &cli.Command{
		Name:   CommandName,
		Usage:  "Emits limited terragrunt state on stdout and exits.",
		Action: func(ctx *cli.Context) error { return Run(opts.OptionsFromContext(ctx)) },
	}
}
