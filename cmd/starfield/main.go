package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/starmap"
	"github.com/skycoin/cx-game/utility"
	"github.com/urfave/cli/v2"
)

//Press TAB to shuffle stars

func init() {
	// seed rand so stars will be random each program run
	rand.Seed(time.Now().UnixNano())
	//lock thread so drawing will be only in main thread, otherwise there will be errors
	runtime.LockOSThread()
}

type noiseSettings struct {
	Size     int
	Scale    float32
	Levels   uint8
	Contrast float32

	Seed        int64
	Gradmax     int
	X           int
	Xs          int
	Persistance float32
	Lacunarity  float32
	Octaves     int

	GradFile string
}
type starSettings struct {
	PixelSize int
}

type Star struct {
	// Drawable uint32
	X             float32
	Y             float32
	Size          float32
	SpriteId      int
	GradientValue float32
	GradientId    int32
	Depth         float32
}

var (
	width           = 800
	height          = 600
	stars           []*Star
	backgroundStars []*Star

	//cli options
	background int = 1 //0 is black, 1 is rgb
	starAmount int = 20

	starConfig  *starSettings  = &starSettings{}
	noiseConfig *noiseSettings = &noiseSettings{}
)

func main() {
	//parse command line arguments and flags
	initArgs()

	// initialize both glfw and gl libraries, setting up the window and shader program
	win := render.NewWindow(height, width, true)
	defer glfw.Terminate()
	spriteloader.InitSpriteloader(&win)
	window := win.Window

	window.SetKeyCallback(keyCallback)
	// program1 := win.Program
	// program2 := render.InitOpenGLCustom("./cmd/starfield/shaders/")
	shader := utility.NewShader("./cmd/starfield/shaders/vertex.glsl", "./cmd/starfield/shaders/fragment.glsl")

	if background == 1 {
		starmap.Init(&win)
		starmap.Generate(256, 0.08, 3)
	}
	//reload yaml config in a goroutine
	// go checkAndReload("config.yaml", &config)
	//randomize stars
	initStarField(&win)

	//bind gradient 1d textures
	for i := 1; i < 12; i++ {
		tex := getGradient(uint(i))
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_1D, tex)
	}

	shader.SetInt("texture_1d", 1)
	shader.Use()
	//main loop
	for !window.ShouldClose() {
		//clearing buffers
		gl.ClearColor(7.0/255.0, 8.0/255.0, 25.0/255.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		drawStarField(shader)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

//callback function to register key events
func keyCallback(w *glfw.Window, k glfw.Key, scancode int, a glfw.Action, m glfw.ModifierKey) {
	if a != glfw.Press {
		return
	}
	if k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	switch k {
	case glfw.KeyTab:
		shuffle()
	}
}

//function to parse cli flags
func initArgs() {

	app := cli.NewApp()
	app.Name = "starfield-cli"
	app.Description = "starfield example"
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:        "background",
			Aliases:     []string{"bg", "b"},
			Usage:       "background to use",
			Value:       0,
			Destination: &background,
		},
		&cli.IntFlag{
			Name:        "stars",
			Aliases:     []string{"star"},
			Usage:       "number of stars to draw",
			Value:       15,
			Destination: &starAmount,
		},
		&cli.IntFlag{
			Name:        "width",
			Usage:       "Resolution width",
			Value:       800,
			Destination: &width,
		},
		&cli.IntFlag{
			Name:        "height",
			Usage:       "Resolution height",
			Value:       600,
			Destination: &height,
		},
	}
	app.After = func(c *cli.Context) error {
		command := c.Args().First()
		if command == "help" {
			os.Exit(0)
		}
		return nil
	}
	app.Action = func(c *cli.Context) error {
		return nil
	}
	app.Run(os.Args)
}

//function to shuffle stars on the background
func shuffle() {
	for _, star := range backgroundStars {
		star.SpriteId = spriteloader.GetSpriteIdByName(fmt.Sprintf("background-stars-%d", rand.Intn(15)))
		star.Size = getSize()
	}
	for _, star := range stars {
		star.X, star.Y = getStarPosition()
		star.SpriteId = spriteloader.GetSpriteIdByName(fmt.Sprintf("stars-%d", rand.Intn(16)))
		star.Size = getSize()
	}
}

//create random stars
func initStarField(win *render.Window) {
	//spriteloader init
	spriteloader.InitSpriteloader(win)
	backgroundStarsheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/starfield_test_16x16_tiles_8x8_tile_grid_128x128.png")
	// galaxySheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/galaxy_256x256.png")
	planetsSheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/planets.png")

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			spriteloader.LoadSprite(backgroundStarsheetId,
				fmt.Sprintf("background-stars-%d", y*4+x),
				x, y)
		}
	}
	//load all sprites from spritesheet
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			//for stars
			spriteloader.LoadSprite(planetsSheetId,
				fmt.Sprintf("stars-%d", y*4+x),
				x, y,
			)
		}
	}

	for x := 0; x < win.Width/60; x++ {
		for y := 0; y < win.Height/60; y++ {
			backgroundStars = append(backgroundStars, &Star{
				X:    float32(x - win.Width/120),
				Y:    float32(y - win.Height/120),
				Size: getSize(),
				// Size:     1,
				SpriteId:      spriteloader.GetSpriteIdByName(fmt.Sprintf("background-stars-%d", rand.Intn(16))),
				GradientValue: rand.Float32(),
				GradientId:    int32(rand.Intn(10) + 1),
			})
		}
	}

	for i := 0; i < starAmount; i++ {
		star := &Star{
			//bad generation position, TODO
			Size:          1,
			SpriteId:      spriteloader.GetSpriteIdByName(fmt.Sprintf("stars-%d", rand.Intn(16))),
			GradientValue: rand.Float32(),
		}
		star.X, star.Y = getStarPosition()
		stars = append(stars, star)
	}

}

func drawStarField(shader *utility.Shader) {
	//background stars
	for _, star := range backgroundStars {
		// spriteloader.DrawSpriteQuad(star.X, star.Y, star.Size, star.Size, star.SpriteId)
		// spriteloader.DrawSpriteQuad(star.X, star.Y, star.Size*(1+config.PixelSize/10), star.Size*(1+config.PixelSize/10), star.SpriteId)
		shader.SetInt("texture_1d", star.GradientId)
		shader.SetFloat("gradValue", star.GradientValue)
		spriteloader.DrawSpriteQuadCustom(star.X, star.Y, star.Size, star.Size, star.SpriteId, shader.ID)
	}

	shader.Use()
	for _, star := range stars {
		shader.SetFloat("gradValue", star.GradientValue)
		spriteloader.DrawSpriteQuadCustom(star.X, star.Y, 1, 1, star.SpriteId, shader.ID)
	}
}

func getSize() float32 {
	size := rand.Float32()/2 + 0.75
	if size > 0.5 && size < 0.75 {
		size = rand.Float32() / 4
	}
	return size
}

//get gradient file
func getGradient(gradientNumber uint) uint32 {
	var tex uint32

	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_1D, tex)

	result, img, _ := spriteloader.LoadPng(filepath.Join("./assets/starfield/gradients", fmt.Sprintf("heightmap_gradient_%02d.png", gradientNumber)))
	if result != 0 {
		log.Panic("Could not load picture!")
	}

	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage1D(gl.TEXTURE_1D, 0, gl.RGBA, int32(img.Rect.Size().X), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return tex
}

//todo
func getStarPosition() (float32, float32) {
	starGap := 0.7
	fmt.Println("abc")
	xPos, yPos := rand.Float32()*8-4, rand.Float32()*7-4
	//if too many stars
	if starAmount > 20 || starGap > 1.3 {
		return xPos, yPos
	}
	for _, star := range stars {
		if math.Abs(float64(xPos-star.X)) < float64(starGap) && math.Abs(float64(yPos-star.Y)) < float64(starGap) {
			return getStarPosition()
		}
	}
	return xPos, yPos
}

func gaussianTheta(x32, y32 float32) float32 {
	x, y := float64(x32), float64(y32)
	var sigmaX, sigmaY, x0, y0 float64
	var A float64 = 1
	var theta float64
	sigmaX = 0.1
	sigmaY = 0.3
	x0 = 0.5
	y0 = 0.5
	theta = float64(DegToRad(-45))
	// a := math.Pow(math.Cos(theta), 2)/(2*math.Pow(sigmaX, 2)) + math.Pow(math.Sin(theta), 2)/(2*math.Pow(sigmaY, 2))
	// b := -math.Sin(2*theta)/(4*math.Pow(sigmaX, 2)) + math.Sin(2*theta)/(4*math.Pow(sigmaY, 2))
	// c := math.Pow(math.Sin(theta), 2)/(2*math.Pow(sigmaX, 2)) + math.Pow(math.Cos(theta), 2)/(2*math.Pow(sigmaY, 2))

	// result := A * math.Exp(-((a * math.Pow((x-x0), 2)) + 2*b*(x-x0)*(y-y0) + c*math.Pow((y-y0), 2)))

	a := math.Pow(math.Cos(theta), 2)/(2*math.Pow(sigmaX, 2)) + math.Pow(math.Sin(theta), 2)/(2*math.Pow(sigmaY, 2))
	b := -math.Sin(2*theta)/(4*math.Pow(sigmaX, 2)) + math.Sin(2*theta)/(4*math.Pow(sigmaY, 2))
	c := math.Pow(math.Sin(theta), 2)/(2*math.Pow(sigmaX, 2)) + math.Pow(math.Cos(theta), 2)/(2*math.Pow(sigmaY, 2))

	fmt.Println(a, b, c)

	result := A * math.Exp(-(a*math.Pow(x-x0, 2) + 2*b*(x-x0)*(y-y0) + c*math.Pow(y-y0, 2)))
	// result := A * math.Exp(-float64(z))
	// if math.Abs(float64(x)) > 0.8 && math.Abs(float64(y)) > 0.8 {
	// 	fmt.Println(x, " ", y, " ", result)
	// }
	return float32(result)
}

func DegToRad(angle float32) float32 {
	return math.Pi / 180.0 * angle
}
