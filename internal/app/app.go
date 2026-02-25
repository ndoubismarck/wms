package app

import (
	"os"
	"os/signal"
	"path/filepath"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/osutil"
	"time"
)

type App struct {
	ctx       types.IContext
	hooksChan chan types.Hook
}

func New(ctx types.IContext, hooksChan chan types.Hook) *App {
	return &App{
		ctx:       ctx,
		hooksChan: hooksChan,
	}
}

func (a *App) Run() error {
	if err := a.init(); err != nil {
		return err
	}
	if err := a.start(); err != nil {
		return err
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, osutil.Signals()...)
	for val := range sig {
		if osutil.IsHangupSignal(val) {
			if err := a.restart(); err != nil {
				return err
			}
		}
		if osutil.IsInterruptSignal(val) {
			if err := a.stop(); err != nil {
				return err
			}
			os.Exit(0)
		}
	}
	return nil
}

func (a *App) init() error {
	exec, err := os.Executable()
	if err != nil {
		return err
	}
	workingDir := filepath.Dir(exec)
	if err := os.Chdir(workingDir); err != nil {
		return err
	}
	if _, err = a.ctx.Config().Get(); err != nil {
		return err
	}
	return nil
}

func (a *App) stop() error {
	hook := types.Hook{
		Name:  types.StopHook,
		Error: make(chan error),
	}
	a.hooksChan <- hook
	err := <-hook.Error
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	return a.stopped()
}

func (a *App) stopped() error {
	hook := types.Hook{
		Name:  types.StoppedHook,
		Error: make(chan error),
	}
	a.hooksChan <- hook
	err := <-hook.Error
	if err != nil {
		return err
	}
	return nil
}

func (a *App) start() error {
	hook := types.Hook{
		Name:  types.StartHook,
		Error: make(chan error),
	}
	a.hooksChan <- hook
	err := <-hook.Error
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	return a.started()
}

func (a *App) started() error {
	hook := types.Hook{
		Name:  types.StartedHook,
		Error: make(chan error),
	}
	a.hooksChan <- hook
	err := <-hook.Error
	if err != nil {
		return err
	}
	return nil
}

func (a *App) restart() error {
	if err := a.stop(); err != nil {
		return err
	}
	return a.start()
}
