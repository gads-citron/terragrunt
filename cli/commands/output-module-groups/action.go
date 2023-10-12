package outputmodulegroups

import (
	"fmt"

	"github.com/gads-citron/terragrunt/cli/commands/terraform"
	"github.com/gads-citron/terragrunt/config"
	"github.com/gads-citron/terragrunt/configstack"
	"github.com/gads-citron/terragrunt/options"
)

func Run(opts *options.TerragruntOptions) error {
	target := terraform.NewTarget(terraform.TargetPointParseConfig, runOutputModuleGroups)

	return terraform.RunWithTarget(opts, target)
}

func runOutputModuleGroups(opts *options.TerragruntOptions, cfg *config.TerragruntConfig) error {
	stack, err := configstack.FindStackInSubfolders(opts, nil)
	if err != nil {
		return err
	}

	js, err := stack.JsonModuleDeployOrder(opts.TerraformCommand)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(opts.Writer, "%s\n", js)
	if err != nil {
		return err
	}

	return nil

}
