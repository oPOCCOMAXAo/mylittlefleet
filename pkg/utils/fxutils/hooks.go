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

type StartableCancellable interface {
	OnStart(context.Context, context.CancelCauseFunc) error
}

func setStartableHook[S any](
	service S,
	hook *fx.Hook,
) {
	if startable, ok := any(service).(Startable); ok {
		hook.OnStart = startable.OnStart
	}
}

func setStoppableHook[S any](
	service S,
	hook *fx.Hook,
) {
	if stoppable, ok := any(service).(Stoppable); ok {
		hook.OnStop = stoppable.OnStop
	}
}

func setStartableCancellableHook[S any](
	service S,
	hook *fx.Hook,
	cancel context.CancelCauseFunc,
) {
	if cancellable, ok := any(service).(StartableCancellable); ok {
		hook.OnStart = func(ctx context.Context) error {
			return cancellable.OnStart(ctx, cancel)
		}
	}
}

// ProvideWithHooks is a helper function for providing a service with lifecycle hooks.
//
// It adds OnStart and OnStop methods of the service to the lifecycle hooks.
//
// The service must implement [Startable], [StartableCancellable] and/or [Stoppable] interfaces.
func ProvideWithHooks[T any](constructor any) fx.Option {
	if !IsTypeWithHooks[T]() {
		return fx.Provide(constructor)
	}

	const tags = `name:"originalWithHooks"`

	return fx.Options(
		fx.Provide(
			fx.Annotate(constructor, fx.ResultTags(tags)),
			fx.Private,
		),
		fx.Provide(
			fx.Annotate(func(
				lc fx.Lifecycle,
				cancel context.CancelCauseFunc,
				svc T,
			) T {
				hook := fx.Hook{}

				setStartableHook(svc, &hook)
				setStoppableHook(svc, &hook)
				setStartableCancellableHook(svc, &hook, cancel)

				if hook.OnStart != nil || hook.OnStop != nil {
					lc.Append(hook)
				}

				return svc
			}, fx.ParamTags(``, ``, tags)),
		),
	)
}

func IsTypeWithHooks[T any]() bool {
	var svc T

	_, t1 := any(svc).(Startable)
	_, t2 := any(svc).(Stoppable)
	_, t3 := any(svc).(StartableCancellable)

	return t1 || t2 || t3
}
