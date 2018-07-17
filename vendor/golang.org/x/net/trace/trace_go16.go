





package trace

import "golang.org/x/net/context"



func NewContext(ctx context.Context, tr Trace) context.Context {
	return context.WithValue(ctx, contextKey, tr)
}


func FromContext(ctx context.Context) (tr Trace, ok bool) {
	tr, ok = ctx.Value(contextKey).(Trace)
	return
}
