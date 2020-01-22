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
	gridBlockCountX int = 10
	gridBlockCountY int = 18
	gridBlockSize   int = 50
	windowSizeX     int = gridBlockCountX * gridBlockSize
	windowSizeY     int = gridBlockCountY * gridBlockSize
)

var (
	fallingTet  Tet
	fallingTetX int
	fallingTetY float64
	fallSpeed   float64 = 2
)

func run() {
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

	blockWhitePic, _ := loadPicture("block_white.png")
	blockColorPic, _ := loadPicture("block_blue.png")

	blockWhiteSprite := pixel.NewSprite(blockWhitePic, blockWhitePic.Bounds())
	blockColorSprite := pixel.NewSprite(blockColorPic, blockWhitePic.Bounds())

	gridScaleFactor := float64(gridBlockSize) / math.Max(blockWhitePic.Bounds().H(), blockWhitePic.Bounds().W())

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

	fallingTet = GetRandomTet()
	fallingTetX = gridBlockCountX / 2
	fallingTetY = float64(gridBlockCountY - 5)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Whitesmoke)

		// Check if row is full
		for y := 0; y < gridBlockCountY; y++ {
			allCurrentXTrue := true
			for x := 0; x < gridBlockCountX; x++ {
				if stateGrid[x][y] == false {
					allCurrentXTrue = false
				}
			}
			// If row is full, move all rows above one block down
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
					blockColorSprite.Draw(win, mat)
				} else {
					blockWhiteSprite.Draw(win, mat)
				}
			}
		}

		// Falling active blocks
		if win.JustPressed(pixelgl.KeyLeft) {
			if fallingTetX >= 1 {
				fallingTetX--
			}
		}
		if win.JustPressed(pixelgl.KeyRight) {
			if (fallingTetX + fallingTet.size) < gridBlockCountX {
				fallingTetX++
			}
		}
		if win.JustPressed(pixelgl.KeyDown) {
			fallSpeed = 20
		}
		if win.JustPressed(pixelgl.KeySpace) {
			fallingTet.Rot()
		}

		fallingTetY -= fallSpeed * dt
		fallingTetDrawY := int(fallingTetY)

		cannotFallFurther := false
		for m := 0; m < fallingTet.size; m++ {
			for n := 0; n < fallingTet.size; n++ {
				if fallingTet.mat[m][n] {
					// Check if tet can move down
					if fallingTetDrawY-1 < 0 || stateGrid[fallingTetX+n][fallingTetDrawY+m-1] {
						cannotFallFurther = true
						break
					}
				}
			}
			if cannotFallFurther {
				break
			}
		}
		if cannotFallFurther {
			for m := 0; m < fallingTet.size; m++ {
				for n := 0; n < fallingTet.size; n++ {
					if fallingTet.mat[m][n] {
						stateGrid[fallingTetX+n][fallingTetDrawY+m] = true
					}
				}
			}
			fallSpeed = 2
			fallingTet = GetRandomTet()
			fallingTetY = float64(gridBlockCountY - 5)
		}

		for m := 0; m < fallingTet.size; m++ {
			for n := 0; n < fallingTet.size; n++ {
				if fallingTet.mat[m][n] {
					mat := pixel.IM
					mat = mat.Scaled(pixel.ZV, gridScaleFactor)
					mat = mat.Moved(pixel.V(float64(gridBlockSize/2), float64(gridBlockSize/2)))
					mat = mat.Moved(pixel.V(float64((fallingTetX+n)*gridBlockSize), float64((int(fallingTetY)+m)*gridBlockSize)))
					blockColorSprite.Draw(win, mat)
				}
			}
		}

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
