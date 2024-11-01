## Table of Contents
- [Description](Description)
- [Features](#features)
- [Installation](#installation)

# What is the Game of Life?

The Game of Life (GoL) is a 0 player game, that is a game that evolves from an initial state. GoL runs on a 2D grid of cells which evolve via the following rules:

- A cell with fewer than 2 neighbours dies from under population.
- A living cell with 2 or 3 neighbours lives on to the next generation.
- A cell with more than 3 neighbours dies from over population.
- A dead cell with exactly 3 neighbours becoems alive in the next generation.

These rules can be simplified into:

- A living cell with 2 or 3 neighbours remains alive in the next generations.
- A dead cell with exactly 3 neighbours becomes alive in the next generation.

## Description

An interactive solution to the game of life using Golang and the `Pixel` 2D game library, where users can draw and erase cells, and move the camera throughout the GoL world.


## Features
- **Interactive controls** for starting, pausing, and resetting the game.
- **Pan** to view different parts of the grid.
- **Toggleable grid information** (frame rate, cell count, current generation).

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
