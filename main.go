package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
	"path/filepath"
	"physics.simulation/repository"
	"strconv"
)

func main() {
	repository.ConnectToDb()

	mux := http.NewServeMux()
	mux.HandleFunc("/simulations", getSimulations)
	mux.HandleFunc("/simulations/{id}", getSimulationByID)
	mux.HandleFunc("/simulations/{id}/nested", getSimulationByIdNested)
	mux.HandleFunc("/simulations/{id}/celestialobjects", getCelestialObjectsBySimulationByID)
	mux.HandleFunc("/simulations/{id}/graph", getSimulationGraphBySimulationID)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})
	handler := corsMiddleware.Handler(mux)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func getSimulations(w http.ResponseWriter, r *http.Request) {
	sims, err, httpStatus := repository.Simulations()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), httpStatus)
		return
	}
	json.NewEncoder(w).Encode(sims)
}

func getSimulationByIdNested(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	sim, err, httpStatus := repository.SimulationByIdLeftJoinChildrenTables(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), httpStatus)
		return
	}
	json.NewEncoder(w).Encode(sim)
}

func getSimulationByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	sim, err, httpStatus := repository.SimulationById(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), httpStatus)
		return
	}
	json.NewEncoder(w).Encode(sim)
}

func getCelestialObjectsBySimulationByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	celestialObjects, err, httpStatus := repository.CelestialObjectsBySimulationByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), httpStatus)
		return
	}
	json.NewEncoder(w).Encode(celestialObjects)
}

func getSimulationGraphBySimulationID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	sim, err, httpStatus := repository.SimulationById(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), httpStatus)
		return
	}
	celestialObjects, err, httpStatus := repository.CelestialObjectsBySimulationByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), httpStatus)
		return
	}
	err = sim.Execute(celestialObjects)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	filePath, err := filepath.Abs(fmt.Sprintf("./data/%d/%v", id, "graph.png"))
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	err = repository.SaveHistoryPositions(celestialObjects)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	err = repository.UpdateSimulationIsDirty(id, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}
