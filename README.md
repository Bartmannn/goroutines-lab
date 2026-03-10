# Concurrent Grid Simulation in Go

This project was built as an experiment with Go and its lightweight approach to concurrent programming. It runs a simulation on a rectangular grid where travelers move independently, dangers appear over time, squatters try to escape when someone steps onto their field, and the whole world is rendered in the terminal in parallel with the simulation itself.

## What the program does

- generates a `30 x 10` grid,
- spawns travelers handled by independent concurrent actors,
- moves each traveler randomly to one of the available neighboring fields,
- periodically spawns dangers marked as `##`,
- periodically spawns squatters marked as `**`,
- refreshes the terminal view independently from the movement logic.

## Demo



## Simulation rules

- A traveler moves randomly to one of the available neighboring fields.
- If a traveler steps onto a danger `##`, the traveler is removed from the board.
- A squatter `**` tries to escape when another entity attempts to enter its field.
- If the squatter has nowhere to move, it disappears and the incoming entity takes its place.
- Dangers and squatters live only for a limited time and then disappear automatically.
- The board edges are blocked, so movement happens only inside the grid.

## Symbols shown on the board

- `00`, `01`, `02`, ...: traveler identifiers
- `##`: danger
- `**`: squatter
- `-` and `|`: traces of recent movement between fields
- empty field: no entity

## Concurrency model

The main idea of the project is to split the simulation into many goroutines that communicate through channels:

- each grid field acts as an independent node with its own state,
- movement between neighbors is handled by message passing through channels,
- a dedicated `Reviver` spawns new travelers, squatters, and dangers,
- a dedicated `Snapshotter` collects updates and renders the board in the terminal,
- simulation logic and rendering run concurrently.

## Timing configuration

The default timing values are defined in [server_solution/consts.go](/home/bartosz-bohdziewicz/University/Semestr5.1/Programowanie współbieżne/lista2/server_solution/consts.go):

- board refresh: every `2s`,
- new traveler: every `6s`,
- new squatter: every `5s`,
- new danger: every `10s`,
- squatter lifetime: `12s`,
- danger lifetime: `17s`,
- traveler movement: randomly every `3-4s`.

## Running the program

Go `1.21.2` or newer is required.

```bash
go run .
```

The simulation runs until `ENTER` is pressed or the process receives `Ctrl+C`.

## Running with Docker

Build the image:

```bash
docker build -t concurrent-grid-simulation .
```

Run the container in interactive mode:

```bash
docker run --rm -it concurrent-grid-simulation
```

The container can be stopped with `ENTER`, `Ctrl+C`, or `docker stop`.

## Project structure

- `main.go`: application entry point
- `server_solution/init.go`: grid and channel initialization
- `server_solution/vertex.go`: single-field logic, movement, and collisions
- `server_solution/reviver.go`: spawning new entities
- `server_solution/snapshotter.go`: terminal rendering

## Good moments to capture in a demo

- the empty board at startup,
- the first traveler appearing and moving across the board,
- a squatter `**` escaping from an occupied field,
- a traveler stepping on `##` and disappearing,
- parallel terminal logs and board refreshes happening at the same time.
