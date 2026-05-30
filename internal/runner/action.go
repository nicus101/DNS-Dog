package runner

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/nicus101/godyndns-ovh/internal/config"
)

type ActionRunner interface {
	Run(context.Context, []config.ActionConfig) error
}

type commandActionRunner struct{}

func (runner commandActionRunner) Run(ctx context.Context, actions []config.ActionConfig) error {
	for _, action := range actions {
		timeout, err := action.Duration()
		if err != nil {
			return HardError{Err: fmt.Errorf("action %q has invalid timeout: %w", action.Name, err)}
		}

		actionCtx, cancel := context.WithTimeout(ctx, timeout)
		cmd := exec.CommandContext(actionCtx, action.Command, action.Args...)
		output, err := cmd.CombinedOutput()
		cancel()
		if err != nil {
			return fmt.Errorf("action %q failed: %w: %s", action.Name, err, string(output))
		}
	}
	return nil
}
