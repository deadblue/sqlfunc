package sqlfunc

import (
	"context"
	"errors"
)

type _ContextKey struct{}

var (
	executorKey _ContextKey

	errNoExecutor = errors.New("can not get executor from context")
)

// NewContext returns a derived context from parent, with executor in it.
func NewContext(parent context.Context, executor Executor) context.Context {
	return context.WithValue(parent, executorKey, executor)
}

func fromContext(ctx context.Context) (executor Executor, err error) {
	val := ctx.Value(executorKey)
	var ok bool
	if executor, ok = val.(Executor); !ok {
		err = errNoExecutor
	}
	return
}
