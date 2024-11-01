package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// TODO:

type Tuple struct {
	X, Y int
}

const (
	width     float64 = 2000
	height    float64 = 1500
	text_size float64 = 2
)

var (
	cell_size   float64 = 20
	update_rate int     = 20
)

func run() { // run Pixel
	var (
		fps          time.Duration = time.Second / time.Duration(update_rate)
		pause_fps    time.Duration = time.Second / 100
		last_time    time.Time
		current_time time.Time
		running      bool    = false
		cells                = make(map[Tuple]struct{})
		debug        bool    = true
		imd                  = imdraw.New(nil)
		step_size    float64 = 2
	)
	camera_offset := pixel.V(0, 0)

	ticker := time.NewTicker(pause_fps) // manage max frame delay
	defer ticker.Stop()

	for i := 0; len(cells) < 20; i++ { //spawns a bunch of random cells
		cell := Tuple{}
		cell.X, cell.Y = rand.Intn(1400)+200, rand.Intn(700)+200
		cells[cell] = struct{}{}
	}
	cfg := pixelgl.WindowConfig{
		Title:  "Game of Life!",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	add_rect := func(x float64, y float64, cell_size float64) {

		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(pixel.V(x, y))
		imd.Color = pixel.RGB(0, 255, 0)
		imd.Push(pixel.V(x+float64(cell_size), y+float64(cell_size)))
		imd.Rectangle(0)
	}

	draw_cells := func(cells map[Tuple]struct{}) {
		imd.Clear()
		for cell := range cells {
			var x float64 = (float64(cell.X) - camera_offset.X) * cell_size
			var y float64 = (float64(cell.Y) - camera_offset.Y) * cell_size
			add_rect(x, y, float64(cell_size))

		}

	}

	play_life := func() {
		draw_cells(cells)
		if running {
			cells = next_evolution(cells)

		}
	}

	user_input := func() {

		switch {

		case win.Pressed(pixelgl.MouseButtonLeft) && !running:
			mousePos := win.MousePosition()
			x := int(mousePos.X + camera_offset.X*cell_size)
			y := int(mousePos.Y + camera_offset.Y*cell_size)
			scale := int(cell_size)
			new_cell := Tuple{x / scale, y / scale}

			_, exists := cells[new_cell]

			if !exists {
				cells[new_cell] = struct{}{}
			}

		case win.Pressed(pixelgl.MouseButtonRight) && !running:
			mousePos := win.MousePosition()
			x := int(mousePos.X + camera_offset.X*cell_size)
			y := int(mousePos.Y + camera_offset.Y*cell_size)
			scale := int(cell_size)
			rem_cell := Tuple{x / scale, y / scale}
			delete(cells, rem_cell)

		case win.Pressed(pixelgl.KeyUp):
			update_rate++
			fps = time.Second / time.Duration(update_rate)
			ticker.Reset(fps)

		case win.Pressed(pixelgl.KeyDown) && update_rate > 1:
			update_rate--
			fps = time.Second / time.Duration(update_rate)
			ticker.Reset(fps)

		case win.JustPressed(pixelgl.KeySpace):
			running = !running
			if running {
				ticker.Reset(fps)
			} else {
				ticker.Reset(pause_fps)
			}

		case win.JustPressed(pixelgl.KeyQ):
			os.Exit(0)

		case win.JustPressed(pixelgl.KeyC):
			cells = make(map[Tuple]struct{})

		case win.JustPressed(pixelgl.KeyH):
			debug = !debug

		case win.Pressed(pixelgl.KeyW):
			camera_offset.Y += step_size

		case win.Pressed(pixelgl.KeyS):
			camera_offset.Y -= step_size

		case win.Pressed(pixelgl.KeyA):
			camera_offset.X -= step_size

		case win.Pressed(pixelgl.KeyD):
			camera_offset.X += step_size
		}
	}
	draw_instructions := func() {
		atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		txt := text.New(pixel.V(10, win.Bounds().H()-20), atlas)

		fmt.Fprintln(txt, "Press SPACE to start/pause.")
		fmt.Fprintln(txt, "Click to draw new cells.")
		fmt.Fprintln(txt, "Press 'C' to reset the grid.")
		fmt.Fprintln(txt, "Press 'Q' to quit.")
		fmt.Fprintln(txt, "Press 'Up' / 'Down' to change the frame limit.")
		fmt.Fprintln(txt, "Press 'H' to toggle info.")

		txt.Draw(win, pixel.IM.Scaled(txt.Orig, text_size))

	}

	draw_info := func() {
		mouse_pos := win.MousePosition()
		mouse_x, mouse_y := int(mouse_pos.X), int(mouse_pos.Y)

		elapsed := current_time.Sub(last_time)
		last_time = current_time
		current_time = time.Now()

		atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		txt := text.New(pixel.V(0, 0), atlas)

		fmt.Fprintf(txt, "Mouse Position in Game: %d, %d\n", mouse_x, mouse_y)
		fmt.Fprintf(txt, "Camera Pos: %d, %d\n", int(camera_offset.X), int(camera_offset.Y))
		fmt.Fprintf(txt, "Max frame time: %d ms \n", fps.Milliseconds())
		fmt.Fprintf(txt, "Frame time: %d ms \n", elapsed.Milliseconds())
		fmt.Fprintf(txt, "Alive cells: %d \n", len(cells))

		text_width := txt.Bounds().W()
		pos := pixel.V(win.Bounds().W()-2*text_width-20, win.Bounds().H()-20)

		txt.Draw(win, pixel.IM.Moved(pos).Scaled(pos, text_size))
	}

	for !win.Closed() { //update loop

		win.Clear(colornames.Black)
		if debug {
			draw_info()
			draw_instructions()
		}

		user_input()
		play_life()

		imd.Draw(win)
		win.Update()
		<-ticker.C
	}
}

func main() {
	pixelgl.Run(run)
}

func neighbours(cells map[Tuple]struct{}) map[Tuple]int { // count the neighbour cells

	frequency := make(map[Tuple]int)

	directions := []Tuple{
		{-1, 1}, {0, 1}, {1, 1},
		{-1, 0}, {1, 0},
		{-1, -1}, {0, -1}, {1, -1},
	} // is the 8 adjacent cells of 0,0.

	for cell := range cells {
		for _, d := range directions {
			adjacent := Tuple{cell.X + d.X, cell.Y + d.Y}
			frequency[adjacent]++
		}
	}
	return frequency
}

func next_evolution(alive_cells map[Tuple]struct{}) map[Tuple]struct{} {
	neighbour_counts := neighbours(alive_cells)
	evolved_cells := make(map[Tuple]struct{})

	for cell, count := range neighbour_counts {
		_, isAlive := alive_cells[cell]
		if count == 3 || count == 2 && isAlive {
			evolved_cells[cell] = struct{}{}
		}
	}
	return evolved_cells
}
