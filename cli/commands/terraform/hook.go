package terraform

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gads-citron/terragrunt/config"
	"github.com/gads-citron/terragrunt/options"
	"github.com/gads-citron/terragrunt/shell"
	"github.com/gads-citron/terragrunt/tflint"
	"github.com/gads-citron/terragrunt/util"
	"github.com/hashicorp/go-multierror"
)

func processErrorHooks(hooks []config.ErrorHook, terragruntOptions *options.TerragruntOptions, previousExecErrors *multierror.Error) error {
	if len(hooks) == 0 || previousExecErrors.ErrorOrNil() == nil {
		return nil
	}

	var errorsOccured *multierror.Error

	terragruntOptions.Logger.Debugf("Detected %d error Hooks", len(hooks))

	customMultierror := multierror.Error{
		Errors: previousExecErrors.Errors,
		ErrorFormat: func(err []error) string {
			result := ""
			for _, e := range err {
				errorMessage := e.Error()
				// Check if is process execution error and try to extract output
				// https://github.com/gads-citron/terragrunt/issues/2045
				originalError := errors.Unwrap(e)
				if originalError != nil {
					processError, cast := originalError.(shell.ProcessExecutionError)
					if cast {
						errorMessage = fmt.Sprintf("%s\n%s", processError.StdOut, processError.Stderr)
					}
				}
				result = fmt.Sprintf("%s\n%s", result, errorMessage)
			}
			return result
		},
	}
	errorMessage := customMultierror.Error()

	for _, curHook := range hooks {
		if util.MatchesAny(curHook.OnErrors, errorMessage) && util.ListContainsElement(curHook.Commands, terragruntOptions.TerraformCommand) {
			terragruntOptions.Logger.Infof("Executing hook: %s", curHook.Name)
			workingDir := ""
			if curHook.WorkingDir != nil {
				workingDir = *curHook.WorkingDir
			}

			var suppressStdout bool
			if curHook.SuppressStdout != nil && *curHook.SuppressStdout {
				suppressStdout = true
			}

			actionToExecute := curHook.Execute[0]
			actionParams := curHook.Execute[1:]

			_, possibleError := shell.RunShellCommandWithOutput(
				terragruntOptions,
				workingDir,
				suppressStdout,
				false,
				actionToExecute, actionParams...,
			)
			if possibleError != nil {
				terragruntOptions.Logger.Errorf("Error running hook %s with message: %s", curHook.Name, possibleError.Error())
				errorsOccured = multierror.Append(errorsOccured, possibleError)
			}
		}
	}
	return errorsOccured.ErrorOrNil()
}

func processHooks(hooks []config.Hook, terragruntOptions *options.TerragruntOptions, terragruntConfig *config.TerragruntConfig, previousExecErrors *multierror.Error) error {
	if len(hooks) == 0 {
		return nil
	}

	var errorsOccured *multierror.Error

	terragruntOptions.Logger.Debugf("Detected %d Hooks", len(hooks))

	for _, curHook := range hooks {
		allPreviousErrors := multierror.Append(previousExecErrors, errorsOccured)
		if shouldRunHook(curHook, terragruntOptions, allPreviousErrors) {
			err := runHook(terragruntOptions, terragruntConfig, curHook)
			if err != nil {
				errorsOccured = multierror.Append(errorsOccured, err)
			}
		}
	}

	return errorsOccured.ErrorOrNil()
}

func shouldRunHook(hook config.Hook, terragruntOptions *options.TerragruntOptions, previousExecErrors *multierror.Error) bool {
	//if there's no previous error, execute command
	//OR if a previous error DID happen AND we want to run anyways
	//then execute.
	//Skip execution if there was an error AND we care about errors

	//resolves: https://github.com/gads-citron/terragrunt/issues/459
	hasErrors := previousExecErrors.ErrorOrNil() != nil
	isCommandInHook := util.ListContainsElement(hook.Commands, terragruntOptions.TerraformCommand)

	return isCommandInHook && (!hasErrors || (hook.RunOnError != nil && *hook.RunOnError))
}

func runHook(terragruntOptions *options.TerragruntOptions, terragruntConfig *config.TerragruntConfig, curHook config.Hook) error {
	terragruntOptions.Logger.Infof("Executing hook: %s", curHook.Name)
	workingDir := ""
	if curHook.WorkingDir != nil {
		workingDir = *curHook.WorkingDir
	}

	var suppressStdout bool
	if curHook.SuppressStdout != nil && *curHook.SuppressStdout {
		suppressStdout = true
	}

	actionToExecute := curHook.Execute[0]
	actionParams := curHook.Execute[1:]

	if actionToExecute == "tflint" {
		if err := executeTFLint(terragruntOptions, terragruntConfig, curHook, workingDir); err != nil {
			return err
		}
	} else {
		_, possibleError := shell.RunShellCommandWithOutput(
			terragruntOptions,
			workingDir,
			suppressStdout,
			false,
			actionToExecute, actionParams...,
		)
		if possibleError != nil {
			terragruntOptions.Logger.Errorf("Error running hook %s with message: %s", curHook.Name, possibleError.Error())
			return possibleError
		}
	}
	return nil
}

func executeTFLint(terragruntOptions *options.TerragruntOptions, terragruntConfig *config.TerragruntConfig, curHook config.Hook, workingDir string) error {
	// fetching source code changes lock since tflint is not thread safe
	rawActualLock, _ := sourceChangeLocks.LoadOrStore(workingDir, &sync.Mutex{})
	actualLock := rawActualLock.(*sync.Mutex)
	actualLock.Lock()
	defer actualLock.Unlock()
	err := tflint.RunTflintWithOpts(terragruntOptions, terragruntConfig, curHook)
	if err != nil {
		terragruntOptions.Logger.Errorf("Error running hook %s with message: %s", curHook.Name, err.Error())
		return err
	}
	return nil
}
