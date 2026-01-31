# Physics Simulation REST-API Service

## Description

REST API service exposing endpoints to manage and run gravitational physics simulations.

Simulation results are determined by the mass, initial position and velocity of the celestial objects that make up its data set.

The simulation itself runs for an amount of simulated time defined by the simulation entity (simulation.duration), updating position, velocity and acceleration vectors of the associated celestial objects every interval of time (simulation.delta_t). Celestial object positions are persisted after a given amount of updates (simulation.writing_rate).

Uses Newton's law of universal gravitation for its update algorithm.

![alt text](newton.png)

## Author

Alexis Chrétien

## Dependencies

* [MySQL](https://dev.mysql.com/downloads/installer/)
* [GnuPlot](http://www.gnuplot.info/)

## Utilisation

```
$ go run . serve
```

## Features
* [X] Rest Api
    * [X] GET /simulations
        * Returns the information of all simulation instances
    * [X] GET /simulations/{id}
        * Returns the information of simulation instance {id}
    * [X] GET /simulations/{id}/nested
        * Returns the information of simulation instance {id} alongside its corresponding child entities (celestial objects, position history (calculated only after running the simulation))
    * [X] GET /simulations/{id}/run
        * Runs simulation {id}, returns a png trace of its celestial objects' spatial evolution
    * [X] POST /simulations
        * Creates a new simulation, with or without celestial objects
    * [X] PATCH /simulations/{id}
        * Modifies simulation {id} and its corresponding celestial objects. Deletes the celestial objects' position histories. Sets the simulation as being 'dirty' (simulation must be ran again to calculate new position histories)
    * [X] DELETE /simulations/{id}
        * Deletes simulation {id}, alongside its corresponding celestial objects and position histories

## To-do
* [ ] Validation on API requests
    * [ ] max length on strings
    * [ ] cap on duration / delta_t / writing_rate / number of celestial object in a given simulation
* [ ] Optimization
    * [ ] Multithreading on the simulation algorithm
* [ ] Implementing users, security