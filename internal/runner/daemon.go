package runner

import (
	"context"
	"errors"
	"time"
)

func RunDaemon(ctx context.Context, runner *Runner) error {
	interval, err := runner.Config.Interval()
	if err != nil {
		return HardError{Err: err}
	}
	initialBackoff, err := runner.Config.InitialBackoff()
	if err != nil {
		return HardError{Err: err}
	}
	maxBackoff, err := runner.Config.MaxBackoff()
	if err != nil {
		return HardError{Err: err}
	}

	backoff := initialBackoff
	for {
		cycleCtx, cancel := context.WithTimeout(ctx, interval)
		_, err := runner.RunCycle(cycleCtx, false)
		cancel()
		if err == nil {
			backoff = initialBackoff
			if err := sleepContext(ctx, interval); err != nil {
				return err
			}
			continue
		}
		if IsHardError(err) {
			return err
		}
		if runner.Logger != nil {
			runner.Logger.Printf("transient error: %v; retrying in %s", err, backoff)
		}
		if err := sleepContext(ctx, backoff); err != nil {
			return err
		}
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

func sleepContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil
		}
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
