package server

import (
	"context"
	"net/http"
	"xm/pkg/handlers"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

var Module = fx.Options(fx.Invoke(Init))

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Handlers  handlers.Handlers
}

func Init(p Params) {
	mux := mux.NewRouter()

	mux.Handle("/sign-up", p.Handlers.LogRequest(http.HandlerFunc(p.Handlers.SignUp))).Methods("POST")
	mux.Handle("/sign-in", p.Handlers.LogRequest(http.HandlerFunc(p.Handlers.SignIn))).Methods("POST")

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	p.Lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go server.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		},
	)
}
