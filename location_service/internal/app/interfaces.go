package app

import "context"

type App interface {
	Serve() error
	Shutdown(ctx context.Context) error
}
