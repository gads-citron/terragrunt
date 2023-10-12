package outputmodulegroups

import (
	"github.com/gads-citron/terragrunt/options"
	"github.com/gads-citron/terragrunt/pkg/cli"
)

const (
	CommandName = "output-module-groups"
)

func NewCommand(opts *options.TerragruntOptions) *cli.Command {
	return &cli.Command{
		Name:   CommandName,
		Usage:  "Output groups of modules ordered for apply as a list of list in JSON (useful for CI use cases).",
		Action: func(ctx *cli.Context) error { return Run(opts.OptionsFromContext(ctx)) },
	}
}
