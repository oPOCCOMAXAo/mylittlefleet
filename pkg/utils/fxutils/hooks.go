package fxutils

import (
	"context"

	"go.uber.org/fx"
)

type Startable interface {
	OnStart(context.Context) error
}

type Stoppable interface {
	OnStop(context.Context) error
}

type StartableCancelleable interface {
	OnStart(context.Context, context.CancelCauseFunc) error
}

func StartableHook[S any](service S) (fx.Hook, bool) {
	if service, ok := any(service).(Startable); ok {
		return fx.Hook{
			OnStart: service.OnStart,
		}, true
	}

	return fx.Hook{}, false
}

func StoppableHook[S any](service S) (fx.Hook, bool) {
	if service, ok := any(service).(Stoppable); ok {
		return fx.Hook{
			OnStop: service.OnStop,
		}, true
	}

	return fx.Hook{}, false
}

func StartableCancelleableHook[S any](
	service S,
	cancel context.CancelCauseFunc,
) (fx.Hook, bool) {
	if svc, ok := any(service).(StartableCancelleable); ok {
		return fx.Hook{
			OnStart: func(ctx context.Context) error {
				return svc.OnStart(ctx, cancel)
			},
		}, true
	}

	return fx.Hook{}, false
}

func Decorate[S any](
	lc fx.Lifecycle,
	service S,
	cancel context.CancelCauseFunc,
) S {
	hook, ok := StartableHook(service)
	if ok {
		lc.Append(hook)
	}

	hook, ok = StoppableHook(service)
	if ok {
		lc.Append(hook)
	}

	hook, ok = StartableCancelleableHook(service, cancel)
	if ok {
		lc.Append(hook)
	}

	return service
}
