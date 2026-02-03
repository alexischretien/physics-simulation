package api

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
	"physics.simulation/app"
)

type API struct {
	App    *app.App
	Config *Config
}

func New(app *app.App) (*API, error) {
	config, err := InitConfig()
	if err != nil {
		return nil, err
	}
	return &API{App: app, Config: config}, nil
}

func (api *API) Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /simulations", api.getSimulations)
	mux.HandleFunc("POST /simulations", api.createSimulation)
	mux.HandleFunc("GET /simulations/{id}", api.getSimulation)
	mux.HandleFunc("PATCH /simulations/{id}", api.updateSimulation)
	mux.HandleFunc(`DELETE /simulations/{id}`, api.deleteSimulation)
	mux.HandleFunc("GET /simulations/{id}/nested", api.getSimulationByIdNested)
	mux.HandleFunc("GET /simulations/{id}/run", api.runSimulation)
}

func (api *API) Cors(mux *http.ServeMux) http.Handler {
	if api.Config.Cors {
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins:   api.Config.AllowedHosts,
			AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
			AllowCredentials: true,
			Debug:            true,
		})
		return corsMiddleware.Handler(mux)
	}
	return mux
}

func (a *API) RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}
