package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\t go run counter.go [filename] [line] [axis (x/y)] [width]")
		return
	}

	// parse args
	file := os.Args[1]
	line, _ := strconv.Atoi(os.Args[2])
	axis := os.Args[3]
	width, _ := strconv.Atoi(os.Args[4])

	video, err := gocv.VideoCaptureFile(file)
	if err != nil {
		fmt.Printf("Error opening video capture file: %s\n", file)
		return
	}
	defer video.Close()

	window := gocv.NewWindow("Track Window")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	imgFG := gocv.NewMat()
	defer imgFG.Close()

	imgCleaned := gocv.NewMat()
	defer imgCleaned.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()

	count := 0
	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", file)
			return
		}
		if img.Empty() {
			continue
		}

		// clean frame by removing background & eroding to eliminate artifacts
		mog2.Apply(img, &imgFG)
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		gocv.Erode(imgFG, &imgCleaned, kernel)
		kernel.Close()

		// calculate the image moment based on the cleaned frame
		moments := gocv.Moments(imgCleaned, true)
		area := moments["m00"]
		if area >= 1 {
			x := int(moments["m10"] / area)
			y := int(moments["m01"] / area)

			if axis == "y" {
				if x > 0 && x < img.Cols() && y > line && y < line+width {
					count++
				}
				gocv.Line(&img, image.Pt(0, line), image.Pt(img.Cols(), line), color.RGBA{255, 0, 0, 0}, 2)
			}
			if axis == "x" {
				if y > 0 && y < img.Rows() && x > line && x < line+width {
					count++
				}
				gocv.Line(&img, image.Pt(line, 0), image.Pt(line, img.Rows()), color.RGBA{255, 0, 0, 0}, 2)
			}
		}

		gocv.PutText(&img, fmt.Sprintf("Count: %d", count), image.Pt(10, 20),
			gocv.FontHersheyPlain, 1.2, color.RGBA{0, 255, 0, 0}, 2)

		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

// O que essa aplicacao faz::
//
// Este exemplo rastreia objetos como carros ou pessoas em transito podendo-se usar
// uma linha horizontal ou vertical, atraves do Moments method.
// O algoritmo do Moments method não é tão preciso para a contagem de múltiplos objetos,
// no entanto, é eficiente na execução.
//
// Como executar:
//
// 		Casos Gerais: go run counter.go /path/to/video.avi 200 y 20
//		Usando o Video desse Tutorial: go run counter.go froggerHighway.mp4 600 x 30(the best case for the cars video; mp4 file is attached in documentation)
//
// Este exemplo contem algumas modificacoes, mas teve como base o exemplo do tutorial existente em:
// [1] https://gocv.io/writing-code/
