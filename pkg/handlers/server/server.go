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

	mux.Handle("/company/create", p.Handlers.LogRequest(p.Handlers.Middleware(http.HandlerFunc(p.Handlers.CreateCompany)))).Methods("POST")
	mux.Handle("/company", p.Handlers.LogRequest(p.Handlers.Middleware(http.HandlerFunc(p.Handlers.GetCompanyByID)))).Methods("GET")
	mux.Handle("/company", p.Handlers.LogRequest(p.Handlers.Middleware(http.HandlerFunc(p.Handlers.DeleteCompany)))).Methods("DELETE")
	mux.Handle("/companies", p.Handlers.LogRequest(p.Handlers.Middleware(http.HandlerFunc(p.Handlers.GetAllCompanies)))).Methods("POST")
	mux.Handle("/company/update", p.Handlers.LogRequest(p.Handlers.Middleware(http.HandlerFunc(p.Handlers.UpdateCompany)))).Methods("PATCH")

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
