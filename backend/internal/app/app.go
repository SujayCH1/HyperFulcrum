package App

import (
	"context"
)

type Application struct {
	// pkg
	ctx context.Context
	// emitter *logger.LogEmitter {add later}

	// api server

	// router config

	// repositories

	// services
}

func NewApp(ctx context.Context) (*Application, error) {
	return &Application{}, nil
}

func (a *Application) Start(ctx context.Context, app *Application) (*Application, error) {
	// start all services
	return &Application{}, nil
}

func (a *Application) Stop(ctx context.Context, app *Application) (*Application, error) {
	// stop all services
	return &Application{}, nil
}
