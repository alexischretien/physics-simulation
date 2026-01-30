package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (api *API) getCelestialObjectsBySimulationByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	exists, err := api.App.Database.SimulationExists(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, fmt.Sprintf("Simulation %d Not found", id), http.StatusNotFound)
		return
	}
	celestialObjects, err := api.App.Database.CelestialObjectsBySimulationByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(celestialObjects)
}
