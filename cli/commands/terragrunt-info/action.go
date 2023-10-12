package terragruntinfo

import (
	"encoding/json"
	"fmt"

	"github.com/gads-citron/terragrunt/cli/commands/terraform"
	"github.com/gads-citron/terragrunt/config"
	"github.com/gads-citron/terragrunt/options"
)

func Run(opts *options.TerragruntOptions) error {
	target := terraform.NewTarget(terraform.TargetPointDownloadSource, runTerragruntInfo)

	return terraform.RunWithTarget(opts, target)
}

// Struct is output as JSON by 'terragrunt-info':
type TerragruntInfoGroup struct {
	ConfigPath       string
	DownloadDir      string
	IamRole          string
	TerraformBinary  string
	TerraformCommand string
	WorkingDir       string
}

func runTerragruntInfo(opts *options.TerragruntOptions, cfg *config.TerragruntConfig) error {
	group := TerragruntInfoGroup{
		ConfigPath:       opts.TerragruntConfigPath,
		DownloadDir:      opts.DownloadDir,
		IamRole:          opts.IAMRoleOptions.RoleARN,
		TerraformBinary:  opts.TerraformPath,
		TerraformCommand: opts.TerraformCommand,
		WorkingDir:       opts.WorkingDir,
	}

	b, err := json.MarshalIndent(group, "", "  ")
	if err != nil {
		opts.Logger.Errorf("JSON error marshalling terragrunt-info")
		return err
	}
	fmt.Fprintf(opts.Writer, "%s\n", b)

	return nil

}
