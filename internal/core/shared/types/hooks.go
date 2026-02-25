package types

type Hook struct {
	Name  HookName
	Error chan error
}

type HookName string

const (
	StopHook    HookName = "stop"
	StartHook            = "start"
	StartedHook          = "started"
	StoppedHook          = "stopped"
)

type IHooks interface {
	OnStop(fn func() error)
	OnStart(fn func() error)
	OnStarted(fn func() error)
	OnStopped(fn func() error)
}
