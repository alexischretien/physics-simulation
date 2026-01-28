package model

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Simulation struct {
	Id               int64
	Title            string
	Duration         float64
	Delta_t          float64
	WritingRate      int
	IsDirty          bool
	CelestialObjects []CelestialObject
}

type SimulationRow struct {
	// Simulation table
	Id          int64
	Title       string
	Duration    float64
	Delta_t     float64
	WritingRate int
	IsDirty     bool

	// Celestial Object table
	CelestialObjectId *int64
	Name              *string
	Mass              *float64
	X_position        *float64
	Y_position        *float64
	Z_position        *float64
	X_velocity        *float64
	Y_velocity        *float64
	Z_velocity        *float64

	// Position History table
	PositionHistoryId *int64
	Time              *float64
	X                 *float64
	Y                 *float64
	Z                 *float64
}

func (s Simulation) Execute(celestialObjects []CelestialObject) error {
	if s.IsDirty == true {
		for i, t := 0, 0.0; t < s.Duration; t += s.Delta_t {
			if t == 0.0 {
				initDataFiles(s, celestialObjects)
			}
			for j := 0; j < len(celestialObjects); j++ {
				celestialObjects[j].UpdateAccelerationVelocityPosition(celestialObjects, j, s.Delta_t)

				if i != 0 && i%s.WritingRate == 0 {
					celestialObjects[j].AppendCurrentPositionToDataFile()
					celestialObjects[j].PositionHistory = append(celestialObjects[j].PositionHistory, PositionHistory{CelestialObjectId: celestialObjects[j].Id, Time: t, Position: celestialObjects[j].Position})
				}
			}
			i++
		}
		defer closeDataFiles(celestialObjects)
	}
	err := runGnuPlot(s, celestialObjects)
	return err
}

func initDataFiles(s Simulation, celestialObjects []CelestialObject) error {
	dataFolder, err := filepath.Abs(fmt.Sprintf("./data/%v", strconv.FormatInt(s.Id, 10)))
	if err != nil {
		return fmt.Errorf("initDataFiles for simulation %d: %v", s.Id, err)
	}

	err = os.RemoveAll(dataFolder)
	if err != nil {
		return fmt.Errorf("initDataFiles for simulation %d: %v", s.Id, err)
	}

	err = os.MkdirAll(dataFolder, 0700)
	if err != nil {
		return fmt.Errorf("initDataFiles for simulation %d: %v", s.Id, err)
	}

	gnuPlotFile, err := os.OpenFile(fmt.Sprintf("%s/gnuplot.gpi", dataFolder), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("initDataFiles for simulation %d: %v", s.Id, err)
	}

	strTemplate1 := "\"data/%d/%d.dat\" index 0 smooth path title \"%v\" with lines"
	strTemplate2 := `set term png
set size 1,1
set terminal png size 840,840
set xtics rotate by -45
set xlabel "X (meters)"
set ylabel "Y (meters)"
set output "data/%d/graph.png"
set title "Simulation trace, duration = %v seconds, delta-t = %v seconds"
plot `
	gnuPlotFile.Write([]byte(fmt.Sprintf(strTemplate2, s.Id, s.Duration, s.Delta_t)))

	for i := range celestialObjects {
		id := celestialObjects[i].Id
		dataFile, err := os.OpenFile(fmt.Sprintf("%s/%d.dat", dataFolder, id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("initDataFiles for simulation %d: %v", s.Id, err)
		}

		celestialObjects[i].setDataFile(dataFile)
		celestialObjects[i].AppendCurrentPositionToDataFile()
		celestialObjects[i].PositionHistory = append(celestialObjects[i].PositionHistory, PositionHistory{CelestialObjectId: celestialObjects[i].Id, Time: 0.0, Position: celestialObjects[i].Position})
		gnuPlotFile.Write([]byte(fmt.Sprintf(strTemplate1, s.Id, id, celestialObjects[i].Name)))
		if i < len(celestialObjects)-1 {
			gnuPlotFile.Write([]byte(",\\\n"))
		}
	}
	defer gnuPlotFile.Close()
	return nil
}

func closeDataFiles(celestialObjects []CelestialObject) {
	for _, c := range celestialObjects {
		c.CloseDataFile()
	}
}

func runGnuPlot(s Simulation, celestialObject []CelestialObject) error {
	cmd := exec.Command("gnuplot", fmt.Sprintf("data/%d/gnuplot.gpi", s.Id))
	stdout, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("runGnuPlot for simulation %d: %v", s.Id, err)
	}
	if len(stdout) != 0 {
		return fmt.Errorf("runGnuPlot for simulation %d: %v", s.Id, stdout)
	}
	return nil
}
