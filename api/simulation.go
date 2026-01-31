package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"physics.simulation/model"
)

func (api *API) getSim(w http.ResponseWriter, r *http.Request, withCelestialObjects bool) (*model.Simulation, error, int) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%v", err), http.StatusBadRequest
	}
	exists, err := api.App.Database.SimulationExists(id)
	if err != nil {
		return nil, fmt.Errorf("%v", err), http.StatusInternalServerError
	}
	if !exists {
		return nil, fmt.Errorf("Simulation %d Not found", id), http.StatusNotFound
	}
	var sim *model.Simulation
	if withCelestialObjects {
		sim, err = api.App.Database.SimulationByIdLeftJoinChildrenTables(id)
	} else {
		sim, err = api.App.Database.SimulationById(id)
	}
	if err != nil {
		return nil, fmt.Errorf("%v", err), http.StatusInternalServerError
	}
	return sim, nil, http.StatusOK
}

func (api *API) getSimulations(w http.ResponseWriter, r *http.Request) {
	sims, err := api.App.Database.Simulations()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sims)
}

func (api *API) getSimulation(w http.ResponseWriter, r *http.Request) {
	sim, err, status := api.getSim(w, r, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), status)
		return
	}
	json.NewEncoder(w).Encode(sim)
}

func (api *API) getSimulationByIdNested(w http.ResponseWriter, r *http.Request) {
	sim, err, status := api.getSim(w, r, true)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), status)
		return
	}
	json.NewEncoder(w).Encode(sim)
}

func (api *API) runSimulation(w http.ResponseWriter, r *http.Request) {
	sim, err, status := api.getSim(w, r, true)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), status)
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
	filePath, err := filepath.Abs(fmt.Sprintf("./data/%d/%v", sim.Id, "graph.png"))
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	err = api.App.Database.SavePositionHistories(sim.CelestialObjects)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	err = api.App.Database.UpdateSimulationIsDirty(sim.Id, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func (api *API) createSimulation(w http.ResponseWriter, r *http.Request) {
	var s *model.Simulation
	json.NewDecoder(r.Body).Decode(&s)

	id, err := api.App.Database.CreateSimulation(*s)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	if len(s.CelestialObjects) > 0 {
		for _, c := range s.CelestialObjects {
			_, err := api.App.Database.CreateCelestialObjectForSimulation(*id, c)
			if err != nil {
				http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				return
			}
		}
	}
	s, err = api.App.Database.SimulationByIdLeftJoinChildrenTables(*id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(s)
}

func (api *API) updateSimulation(w http.ResponseWriter, r *http.Request) {
	s, err, status := api.getSim(w, r, true)

	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), status)
		return
	}
	var s2 *model.Simulation
	json.NewDecoder(r.Body).Decode(&s2)

	idMap := createCelestialObjectIdMap(*s)
	idMap2 := createCelestialObjectIdMap(*s2)

	// new celestial objects shouldn't have a new Id in the request body
	for _, c := range s2.CelestialObjects {
		if c.Id != 0 && idMap[c.Id] == nil {
			http.Error(w, fmt.Sprintf("Celestial Object %d Not found", c.Id), http.StatusNotFound)
			return
		}
	}
	// creating or updating celestial objects
	for _, c := range s2.CelestialObjects {
		if c.Id == 0 {
			_, err = api.App.Database.CreateCelestialObjectForSimulation(s.Id, c)
		} else {
			err = api.App.Database.UpdateCelestialObject(c)
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), status)
			return
		}
	}

	// deleting missing celestial objects
	for _, c := range s.CelestialObjects {
		if idMap2[c.Id] == nil {
			err = api.App.Database.DeleteCelestialObject(c.Id)
			if err != nil {
				http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				return
			}
		}
	}
	// updating simulation, dirty = true
	err = api.App.Database.UpdateSimulation(*s2)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	// Delete Position Histories
	err = api.App.Database.DeletePositionHistoryForSimulationId(s2.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func createCelestialObjectIdMap(s model.Simulation) map[int64]*model.CelestialObject {
	idMap := make(map[int64]*model.CelestialObject)
	if len(s.CelestialObjects) > 0 {
		for _, c := range s.CelestialObjects {
			idMap[c.Id] = &c
		}
	}
	return idMap
}

func (api *API) deleteSimulation(w http.ResponseWriter, r *http.Request) {
	s, err, status := api.getSim(w, r, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), status)
		return
	}
	err = api.App.Database.DeleteSimulation(s.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}
