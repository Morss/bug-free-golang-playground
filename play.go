package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	gridBlockCountX int     = 10
	gridBlockCountY int     = 18
	gridBlockSize   int     = 50
	windowSizeX     int     = gridBlockCountX * gridBlockSize
	windowSizeY     int     = gridBlockCountY * gridBlockSize
	fallSpeed       float64 = 2
)

func run() {
	t := GetRandom()
	t.Draw()

	cfg := pixelgl.WindowConfig{
		Title:  "Tetris of awesomeness!",
		Bounds: pixel.R(0, 0, float64(windowSizeX), float64(windowSizeY)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)
	win.Clear(colornames.Skyblue)

	gopherGreyPic, _ := loadPicture("hiking_bw.png")
	gopherColorPic, _ := loadPicture("hiking.png")

	gopherGreySprite := pixel.NewSprite(gopherGreyPic, gopherGreyPic.Bounds())
	gopherColorSprite := pixel.NewSprite(gopherColorPic, gopherGreyPic.Bounds())

	gridScaleFactor := float64(gridBlockSize) / math.Max(gopherGreyPic.Bounds().H(), gopherGreyPic.Bounds().W())

	stateGrid := make([][]bool, gridBlockCountX)
	for col := 0; col < gridBlockCountX; col++ {
		stateGrid[col] = make([]bool, gridBlockCountY)
	}

	stateGrid[0][0] = true
	stateGrid[1][0] = true
	stateGrid[2][0] = true
	stateGrid[3][0] = false
	stateGrid[4][0] = true

	stateGrid[2][0] = true
	stateGrid[2][1] = true
	stateGrid[2][2] = true

	stateGrid[1][1] = true

	activeGopherX := 3

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()

		win.Clear(colornames.Whitesmoke)

		// Wipe full rows, cascade rows above
		for y := 0; y < gridBlockCountY; y++ {
			allCurrentXTrue := true
			for x := 0; x < gridBlockCountX; x++ {
				if stateGrid[x][y] == false {
					allCurrentXTrue = false
				}
			}
			if allCurrentXTrue {
				for yy := y; yy < gridBlockCountY; yy++ {
					for x := 0; x < gridBlockCountX; x++ {
						if yy+1 < gridBlockCountY {
							stateGrid[x][yy] = stateGrid[x][yy+1]
						}
					}
				}
			}
		}

		// Draw sprites for grid and fallen blocks
		for x := 0; x < gridBlockCountX; x++ {
			col := stateGrid[x]
			for y := 0; y < gridBlockCountY; y++ {
				mat := pixel.IM
				mat = mat.Scaled(pixel.ZV, gridScaleFactor)
				mat = mat.Moved(pixel.V(float64(gridBlockSize/2), float64(gridBlockSize/2)))
				mat = mat.Moved(pixel.V(float64(x*gridBlockSize), float64(y*gridBlockSize)))
				if col[y] {
					gopherColorSprite.Draw(win, mat)
				} else {
					gopherGreySprite.Draw(win, mat)
				}
			}
		}

		if win.JustPressed(pixelgl.KeyLeft) {
			if activeGopherX >= 1 {
				activeGopherX -= 1
			}
		}
		if win.JustPressed(pixelgl.KeyRight) {
			if activeGopherX < gridBlockCountX-1 {
				activeGopherX += 1
			}
		}

		// Falling active blocks
		activeGopherY := float64(gridBlockCountY) - math.Floor(fallSpeed*dt)

		if int(activeGopherY)-1 < 0 || stateGrid[activeGopherX][int(activeGopherY)-1] {
			last = time.Now()
			stateGrid[activeGopherX][int(activeGopherY)] = true
			activeGopherY = float64(gridBlockCountY)
		}

		mat := pixel.IM
		mat = mat.Scaled(pixel.ZV, gridScaleFactor)
		mat = mat.Moved(pixel.V(float64(gridBlockSize/2), float64(gridBlockSize/2)))
		mat = mat.Moved(pixel.V(float64(activeGopherX*gridBlockSize), activeGopherY*float64(gridBlockSize)))
		gopherColorSprite.Draw(win, mat)

		win.Update()
	}
}

func drawState(grid *[][]bool) {
	for y := gridBlockCountY - 1; y > -1; y-- {
		for x := 0; x < gridBlockCountX; x++ {
			if (*grid)[x][y] {
				fmt.Print("x")
			} else {
				fmt.Print("o")
			}
		}
		fmt.Println()
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func main() {
	pixelgl.Run(run)
}
