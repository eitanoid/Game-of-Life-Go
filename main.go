package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type Tuple struct {
	X, Y int
}

var (
	InputW  = flag.Int("w", 2000, "Set the screen width")
	InputH  = flag.Int("h", 1500, "Set the screen height")
	InputS  = flag.Int("s", 20, "Set the scale of the game")
	inputF  = flag.Int("f", 20, "Set the fps limit")
	inputP  = flag.String("p", "", "Import a PNG")
	inputBS = flag.String("bs", "B3/S23", "select cellular automaton rules in neighbours to be born / neighbours to survive form")
)

func run() { // run Pixel
	var (
		width       float64 = float64(*InputW)
		height      float64 = float64(*InputH)
		cell_size   float64 = float64(*InputS)
		update_rate int     = *inputF

		last_time    time.Time
		current_time time.Time
		fps          time.Duration = time.Second / time.Duration(update_rate)
		pause_fps    time.Duration = time.Second / 100

		cells               = make(map[Tuple]struct{})
		current_gen int     = 0
		running     bool    = false
		debug       bool    = true
		step_size   float64 = 2

		imd           *imdraw.IMDraw = imdraw.New(nil)
		text_size     float64        = 0.0005 * (height + width)
		camera_offset pixel.Vec      = pixel.V(0, 0)

		B = make(map[int]bool) // neighbours to be born / neighbours to survive
		S = make(map[int]bool)
	)

	validate_rules := func(inputBS *string) {
		BS := strings.Split(*inputBS, "/")
		B_in := BS[0][1:]
		S_in := BS[1][1:]
		for _, rule := range B_in {
			r, err := strconv.Atoi(string(rule))

			if err != nil {
				B = map[int]bool{3: true}
				S = map[int]bool{2: true, 3: true}
				fmt.Println("rule invalid, defaulting to game of life")

				return
			}
			B[r] = true
		}
		for _, rule := range S_in {
			r, err := strconv.Atoi(string(rule))
			if err != nil {
				B = map[int]bool{3: true}
				S = map[int]bool{2: true, 3: true}
				fmt.Println("rule invalid, defaulting to game of life")
				return
			}
			S[r] = true
		}
	}

	load_img := func(name string, cells map[Tuple]struct{}) { // load a png into the game

		if name == "" {
			return
		}

		file, err := os.Open(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		bounds := img.Bounds()
		var val uint32 = 0 * 256
		len_x, len_y := bounds.Dx(), bounds.Dy()
		for x := 0; x < len_x; x++ {
			for y := 0; y < len_y; y++ {
				r, g, b, a := img.At(x, y).RGBA()
				if r > val && g > val && b > val && a != 0 {
					cell := Tuple{X: x, Y: len_y - y}
					cells[cell] = struct{}{}
				}
			}
		}

	}

	ticker := time.NewTicker(pause_fps) // manage max frame delay
	defer ticker.Stop()

	cfg := pixelgl.WindowConfig{
		Title:  "Game of Life!",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	add_rect := func(x float64, y float64, cell_size float64) { //adds rect to draw queue

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
			cells = next_evolution(cells, B, S)
			current_gen++

		}
	}

	user_input := func() {

		switch {

		case win.Pressed(pixelgl.MouseButtonLeft) && !running: // drawing new cells
			mousePos := win.MousePosition()
			x := int(mousePos.X + camera_offset.X*cell_size)
			y := int(mousePos.Y + camera_offset.Y*cell_size)
			scale := int(cell_size)
			new_cell := Tuple{x / scale, y / scale}

			_, exists := cells[new_cell]

			if !exists {
				cells[new_cell] = struct{}{}
			}

		case win.Pressed(pixelgl.MouseButtonRight) && !running: // erasing cells
			mousePos := win.MousePosition()
			x := int(mousePos.X + camera_offset.X*cell_size)
			y := int(mousePos.Y + camera_offset.Y*cell_size)
			scale := int(cell_size)
			rem_cell := Tuple{x / scale, y / scale}
			delete(cells, rem_cell)

		case win.Pressed(pixelgl.KeyUp): // changing game speed
			update_rate++
			fps = time.Second / time.Duration(update_rate)
			ticker.Reset(fps)

		case win.Pressed(pixelgl.KeyDown) && update_rate > 1:
			update_rate--
			fps = time.Second / time.Duration(update_rate)
			ticker.Reset(fps)

		case win.JustPressed(pixelgl.KeySpace): // pause game and change run speed for smoother drawing
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

		fmt.Fprintln(txt, "Press 'Q' to quit.")
		fmt.Fprintln(txt, "Use 'LEFT MOUSE' and 'RIGHT MOUSE' to draw and erase cells.")
		fmt.Fprintln(txt, "Press SPACE to start/pause.")
		fmt.Fprintln(txt, "Press 'C' to reset the grid.")
		fmt.Fprintln(txt, "Press 'Up' / 'Down' to change the frame limit.")
		fmt.Fprintln(txt, "Press 'H' to toggle the info overlay.")

		txt.Draw(win, pixel.IM.Scaled(txt.Orig, text_size))

	}

	draw_info := func() {
		mouse_pos := win.MousePosition()
		mouse_x, mouse_y := int(mouse_pos.X), int(mouse_pos.Y)

		elapsed := current_time.Sub(last_time) // frame rate counter
		last_time = current_time
		current_time = time.Now()

		atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		txt := text.New(pixel.V(0, 0), atlas)

		fmt.Fprintf(txt, "Mouse Position in Game: %d, %d\n", mouse_x, mouse_y)
		fmt.Fprintf(txt, "Camera Pos: %d, %d\n", int(camera_offset.X), int(camera_offset.Y))
		fmt.Fprintf(txt, "Max frame time: %d ms \n", fps.Milliseconds())
		fmt.Fprintf(txt, "Frame time: %d ms \n", elapsed.Milliseconds())
		fmt.Fprintf(txt, "Alive cells: %d \n", len(cells))
		fmt.Fprintf(txt, "Current generation: %d\n", current_gen)

		text_width := txt.Bounds().W()
		pos := pixel.V(win.Bounds().W()-2*text_width-20, win.Bounds().H()-20)

		txt.Draw(win, pixel.IM.Moved(pos).Scaled(pos, text_size))
	}

	// init functions
	load_img(*inputP, cells) // loads image
	validate_rules(inputBS)

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
	flag.Parse()
	pixelgl.Run(run)
}

func neighbours(cells map[Tuple]struct{}) map[Tuple]int { // count the neighbour cells of each cell alive

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

func next_evolution(alive_cells map[Tuple]struct{}, B, S map[int]bool) map[Tuple]struct{} { // iterates the game of life
	neighbour_counts := neighbours(alive_cells)
	evolved_cells := make(map[Tuple]struct{})

	for cell, count := range neighbour_counts {
		_, is_alive := alive_cells[cell]
		if B[count] && !is_alive || S[count] && is_alive { // && !isAlive || count == 7 && isAlive {
			evolved_cells[cell] = struct{}{}
		}
	}
	return evolved_cells
}
