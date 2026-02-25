package hooks

import "server/internal/core/shared/types"

type hooks struct {
	onStop    []func() error
	onStart   []func() error
	onStopped []func() error
	onStarted []func() error
}

func New(ch chan types.Hook) types.IHooks {
	h := &hooks{}
	go func() {
		for {
			select {
			case res := <-ch:
				if res.Name == types.StartHook {
					res.Error <- h.callStart()
				}
				if res.Name == types.StopHook {
					res.Error <- h.callStop()
				}
				if res.Name == types.StartedHook {
					res.Error <- h.callStarted()
				}
				if res.Name == types.StoppedHook {
					res.Error <- h.callStopped()
				}
			}
		}
	}()
	return h
}

func (h *hooks) OnStart(fn func() error) {
	h.onStart = append(h.onStart, fn)
}

func (h *hooks) OnStop(fn func() error) {
	h.onStop = append(h.onStop, fn)
}

func (h *hooks) OnStarted(fn func() error) {
	h.onStarted = append(h.onStarted, fn)
}

func (h *hooks) OnStopped(fn func() error) {
	h.onStopped = append(h.onStopped, fn)
}

func (h *hooks) callStart() error {
	for _, val := range h.onStart {
		fn := val
		if fn != nil {
			if err := fn(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *hooks) callStop() error {
	for _, val := range h.onStop {
		fn := val
		if fn != nil {
			if err := fn(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *hooks) callStarted() error {
	for _, val := range h.onStarted {
		fn := val
		if fn != nil {
			if err := fn(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *hooks) callStopped() error {
	for _, val := range h.onStopped {
		fn := val
		if fn != nil {
			if err := fn(); err != nil {
				return err
			}
		}
	}
	return nil
}
