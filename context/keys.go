package contextkeys

import "context"

type ContextKey string

const CallerKey ContextKey = "caller"

// Setter function
func SetCaller(ctx context.Context, caller string) context.Context {
	return context.WithValue(ctx, CallerKey, caller)
}

// Getter function
func GetCaller(ctx context.Context) string {
	caller, ok := ctx.Value(CallerKey).(string)
	if !ok {
		return "unknown"
	}
	return caller
}
