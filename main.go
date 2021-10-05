package main

import (
	"math"
	"os"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type Point [2]float64
type Sphere = [10][10]Point

const (
	FRAMES_PER_SECOND = 60
	MS_PER_FRAME      = uint32(1000 / FRAMES_PER_SECOND)
)

var (
	window *sdl.Window
	render *sdl.Renderer
	right  = true
	scale  = 120.0
	phase  = 0.0
	dp     = 2.5
	dx     = 2.1
	y_ang  = 0.0
	px     = 320.0
	py     = 0.0
)

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, _ = sdl.CreateWindow("Amiga Boing Ball", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 512, sdl.WINDOW_OPENGL)
	render, _ = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	for {
		listen_for_events()
		start_ticks := sdl.GetTicks()
		run_loop()
		render.Present()
		sync_framerate(start_ticks)
	}

}

func run_loop() {
	clear_background()
	do_physics()
	var points = calc_points(math.Mod(phase, 22.5))
	transform(points)
	draw_shadow(points)
	draw_grid()
	fill_tiles(points, phase >= 22.5)
}

func listen_for_events() {
	for {
		evt := sdl.PollEvent()
		escape_pressed := sdl.GetKeyboardState()[sdl.GetScancodeFromKey(sdl.K_ESCAPE)] != 0
		if escape_pressed || (evt != nil && evt.GetType() == sdl.QUIT) {
			sdl.Quit()
			os.Exit(0)
		}
		break
	}
}

func sync_framerate(start_ticks uint32) {
	frame_ms := sdl.GetTicks() - start_ticks
	if frame_ms < MS_PER_FRAME {
		sdl.Delay(MS_PER_FRAME - frame_ms)
	}
}

func clear_background() {
	render.SetDrawColor(170, 170, 170, 255) // light gray
	render.Clear()
}

func do_physics() {
	phase_shift := dp
	if right {
		phase_shift = 45.0 - dp
	}
	phase = math.Mod(phase+phase_shift, 45.0)
	if right {
		px += dx
	} else {
		px -= dx
	}
	if px >= 505 {
		right = false
	} else if px <= 135 {
		right = true
	}
	y_ang = math.Mod((y_ang + 1.5), 360.0)
	py = 350.0 - 200.0*math.Abs(math.Cos(y_ang*math.Pi/180.0))
}

func get_lat(phase float64, i int) float64 {
	if i == 0 {
		return -90.0
	} else if i == 9 {
		return 90.0
	} else {
		return -90.0 + phase + ((float64(i) - 1.0) * 22.5)
	}
}

func calc_points(phase float64) *Sphere {
	sin_lat := make([]float64, 10)
	for i := 0; i < len(sin_lat); i++ {
		lat := get_lat(phase, i)
		sin_lat[i] = math.Sin(lat * math.Pi / 180.0)
	}

	points := Sphere{}
	for j := 0; j < len(sin_lat)-1; j++ {
		lon := -90.0 + float64(j)*22.5
		y := math.Sin(lon * math.Pi / 180.0)
		l := math.Cos(lon * math.Pi / 180.0)
		for k := 0; k < len(sin_lat); k++ {
			x := sin_lat[k] * l
			points[k][j] = Point{x, y}
		}
	}

	return &points
}

func transform(points *Sphere) {
	tilt_sphere(points, 17.0)
	scale_and_translate(points)
}

func tilt_sphere(points *Sphere, ang float64) {
	st := math.Sin(ang * math.Pi / 180.0)
	ct := math.Cos(ang * math.Pi / 180.0)
	for i := 0; i < len(points); i++ {
		for j := 0; j < len(points); j++ {
			pt := points[i][j]
			x := pt[0]*ct - pt[1]*st
			y := pt[0]*st + pt[1]*ct
			points[i][j] = Point{x, y}
		}
	}
}

func scale_and_translate(points *Sphere) {
	for i := 0; i < len(points); i++ {
		for j := 0; j < len(points); j++ {
			pt := points[i][j]
			x := pt[0]*scale + px
			y := pt[1]*scale + py
			points[i][j] = Point{x, y}
		}
	}
}

func draw_shadow(points *Sphere) {
	poly_x, poly_y := make([]int16, 16), make([]int16, 16)

	for i := 0; i <= 8; i++ {
		p := points[0][i]
		poly_x[i] = int16(p[0]) + 50
		poly_y[i] = int16(p[1])
	}
	for i := 0; i <= 8; i++ {
		p := points[9][8-i]
		poly_x[7+i] = int16(p[0]) + 50
		poly_y[7+i] = int16(p[1])
	}

	gfx.FilledPolygonRGBA(render, poly_x, poly_y, 102, 102, 102, 255) // dark gray
}

func draw_grid() {
	render.SetDrawColor(183, 45, 168, 255) // purple
	render.SetScale(1, 1)

	var is = make([]int, 13)
	for i := range is {
		y := int32(i * 36)
		render.DrawLine(50, y, 590, y)
	}

	is = make([]int, 16)
	for i := range is {
		x := int32(50 + i*36)
		render.DrawLine(x, 0, x, 432)
	}

	for i := range is {
		render.DrawLine(int32(50+i*36), 432, int32(float32(i)*42.66), 480)
	}

	ys := []int{442, 454, 468}
	is = make([]int, 3)
	for i := range is {
		y := ys[i]
		x1 := 50.0 - 50.0*(float32(y)-432.0)/(480.0-432.0)
		render.DrawLine(int32(x1), int32(y), int32(640-x1), int32(y))
	}

}

func fill_tiles(points *Sphere, alter bool) {

	poly_n := uint8(4)
	for j := 0; j < 8; j++ {
		for i := 0; i <= 8; i++ {
			p1 := points[i][j]
			p2 := points[i+1][j]
			p3 := points[i+1][j+1]
			p4 := points[i][j+1]
			poly_x, poly_y := make([]int16, poly_n), make([]int16, poly_n)
			poly_x[0] = int16(p1[0])
			poly_y[0] = int16(p1[1])
			poly_x[1] = int16(p2[0])
			poly_y[1] = int16(p2[1])
			poly_x[2] = int16(p3[0])
			poly_y[2] = int16(p3[1])
			poly_x[3] = int16(p4[0])
			poly_y[3] = int16(p4[1])

			r, g, b := uint8(255), uint8(255), uint8(255)
			if alter {
				r, g, b = uint8(255), uint8(0), uint8(0)
			}
			gfx.FilledPolygonRGBA(render, poly_x, poly_y, r, g, b, 255)

			alter = !alter
		}
	}
}
