## Table of Contents
- [Description](Description)
- [Features](#features)
- [Installation](#installation)

# What is the Game of Life?

![Screenshot from 2024-11-06 16-47-29](https://github.com/user-attachments/assets/ff659f00-f321-4097-bdad-9e62617377bb)

[The Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) (GoL) is a Cellular Automaton (CA) / 0 player game. That is a game that evolves from an initial state with no further interaction. 
GoL runs on a 2D grid of cells with 2 states, alive or dead. The cells evolve via the following rules:

- A cell with fewer than 2 neighbours dies from under population.
- A living cell with 2 or 3 neighbours lives on to the next generation.
- A cell with more than 3 neighbours dies from over population.
- A dead cell with exactly 3 neighbours becoems alive in the next generation.

These rules can be simplified into:

- A living cell with 2 or 3 neighbours remains alive in the next generations.
- A dead cell with exactly 3 neighbours becomes alive in the next generation.

Similar CAs with evolutions relating to the the number of neighbours of a cell are called [Life-like cellular automaton](https://en.wikipedia.org/wiki/Life-like_cellular_automaton) and are described by their rule which is denoted by a string "B/S".

For example GoL is described by the rule "B3/S23" which states a cell is Born if it has 3 neighbours and Survives if it has 2 or 3, altering those numbers yields other CAs with different behaviour.
These are some examples taken from the wikipedia page for life-like CAs:

- "B3/S23" - The Game of Life
- "B36/23" - Highlife
- "B36/S125"	- 2x2
- "B1357/S1357"	- Replicator.

## Description

An interactive solution to the game of life using Golang and the `Pixel` 2D game library. Users can draw and erase cells, move the camera throughout the 2D world, import images and change the rules to play different life-like games.


## Features
- **Interactive controls** for starting, pausing, drawing and erasing cells, and resetting the game.
- **Pan** to view different parts of the grid.
- **Toggleable grid information** (frame rate, cell count, current generation).
- **Add images** turns pixels in PNG files into cells whenever the pixel color is not black or fully transparent.
- **Change rules** change the game rules to similar life-like games with "B/S" notation. 

## Usage

You can pass the following flags to the executable:

```bash

usage: Game-Of-Life-Go [<flags>]

Flags:
-w    The width of the grid  (default 2000)
-h    The height of the grid (default 1500)
-s    The scale of the cells (default 20)
-f    Set the game FPS limit (default 20)
-p    Add an image as "image_name.png" (default "")
-bs   Select cellular automaton rules using the "B/S" notation (default "B3/S23")
```


## Installation
### Prerequisites
- **Go**: This project requires the Go programming language. You can download it from [https://golang.org/dl/](https://golang.org/dl/).
- **Pixel Library**: This project relies on the Pixel library for graphics. Install [Pixel](https://github.com/faiface/pixel) and its dependencies:
  ```bash
  go get github.com/faiface/pixel
  go get github.com/faiface/pixel/pixelgl
  go get golang.org/x/image/colornames
  go get github.com/faiface/pixel/imdraw
  go get github.com/faiface/pixel/text
  go get golang.org/x/image/font/basicfont
  ```
