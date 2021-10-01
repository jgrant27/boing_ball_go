package main

import (
	"math"
	"os"

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
	evt    *sdl.Event
	scale  = 120.0
	phase  = 0.0
	dp     = 2.5
	x      = 320.0
	dx     = 2.1
	right  = true
	y_ang  = 0.0
	y      = 0.0
)

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, _ = sdl.CreateWindow("Amiga Boing Ball", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 512, sdl.WINDOW_OPENGL)
	render, _ = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC|sdl.RENDERER_TARGETTEXTURE)

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
	transform(points, scale, x, y)
	draw_shadow(points)
	draw_wireframe()
	fill_tiles(points, phase >= 22.5)

}

func listen_for_events() {
	for {
		evt := sdl.PollEvent()
		keys := sdl.GetKeyboardState()
		if (evt != nil && evt.GetType() == sdl.QUIT) || keys[sdl.K_ESCAPE] != 0 {
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
		x += dx
	} else {
		x -= dx
	}
	if x >= 505 {
		right = false
	} else if x <= 135 {
		right = true
	}
	y_ang = math.Mod((y_ang + 1.5), 360.0)
	y = 350.0 - 200.0*math.Abs(math.Cos(y_ang*math.Pi/180.0))
}

func get_lat(phase float64, i int) float64 {
	if i == 0 {
		return -90.0
	} else if i == 9 {
		return 90.0
	} else {
		return -90.0 + phase + (float64(i)-1.0)*22.5
	}
}

func calc_points(phase float64) Sphere {
	points := Sphere{}
	sin_lat := make([]float64, 10)[:10]
	for i := 0; i < len(sin_lat); i++ {
		lat := get_lat(phase, i)
		sin_lat[i] = math.Sin(lat * math.Pi / 180.0)
	}
	for j := 0; j < len(sin_lat)-1; j++ {
		lon := -90.0 + float64(j)*22.5
		y := math.Sin(lon * math.Pi / 180.0)
		l := math.Cos(lon * math.Pi / 180.0)
		for i := 0; i < len(sin_lat); i++ {
			x := sin_lat[i] * l
			points[i][j] = Point{x, y}
		}
	}

	return points
}

func transform(points Sphere, s float64, tx float64, ty float64) {
	tilt_sphere(points, 17.0)
	scale_and_translate(points, s, tx, ty)
}

func tilt_sphere(points Sphere, ang float64) {

}

func scale_and_translate(points Sphere, s float64, tx float64, ty float64) {

}

func draw_shadow(points Sphere) {
	// var polyX, polyY: array[0..16, int16]

	// for i in 0..8:
	//	let p = points[0][i]
	//	polyX[i] = int16(p[0]) + 50
	//	polyY[i] = int16(p[1])
	// for i in 0..8:
	//	let p = points[9][8 - i]
	//	polyX[7 + i] = int16(p[0]) + 50
	//	polyY[7 + i] = int16(p[1])

	//	var //gray
	//	r: uint8 = 102
	//	g: uint8 = 102
	//	b: uint8 = 102
	// gfx.filledPolygonRGBA(render,
	//	  vx = addr(polyX[0]), vy = addr(polyY[0]), n = 15, r, g, b, 255)

}

func draw_wireframe() {
	render.SetDrawColor(183, 45, 168, 255) // purple
	render.SetScale(1, 1)

	var is = make([]int, 13)[:13]
	for i := range is {
		y := int32(i * 36)
		render.DrawLine(50, y, 590, y)
	}

	is = make([]int, 16)[:16]
	for i := range is {
		x := int32(50 + i*36)
		render.DrawLine(x, 0, x, 432)
	}

	for i := range is {
		render.DrawLine(int32(50+i*36), 432, int32(float32(i)*42.66), 480)
	}

	ys := []int{442, 454, 468}
	is = make([]int, 3)[:3]
	for i := range is {
		y := ys[i]
		x1 := 50.0 - 50.0*(float32(y)-432.0)/(480.0-432.0)
		render.DrawLine(int32(x1), int32(y), int32(640-x1), int32(y))
	}

}

func fill_tiles(points Sphere, alter bool) {

}
