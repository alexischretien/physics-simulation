package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
)

func (api *API) getSimulations(w http.ResponseWriter, r *http.Request) {
	sims, err := api.App.Database.Simulations()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sims)
}

func (api *API) getSimulationByID(w http.ResponseWriter, r *http.Request) {
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
	sim, err := api.App.Database.SimulationById(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sim)
}

func (api *API) getSimulationByIdNested(w http.ResponseWriter, r *http.Request) {
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
	sim, err := api.App.Database.SimulationByIdLeftJoinChildrenTables(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sim)
}

func (api *API) getSimulationGraphBySimulationID(w http.ResponseWriter, r *http.Request) {
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
	sim, err := api.App.Database.SimulationByIdLeftJoinChildrenTables(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	if sim.IsDirty {
		err = api.App.Database.DeletePositionHistoryForSimulationId(sim.Id)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
	}
	if err = sim.Execute(); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	filePath, err := filepath.Abs(fmt.Sprintf("./data/%d/%v", id, "graph.png"))
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	err = api.App.Database.SavePositionHistories(sim.CelestialObjects)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	err = api.App.Database.UpdateSimulationIsDirty(id, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}
