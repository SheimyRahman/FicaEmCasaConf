// This example is all based on the tutorial existing in: https://gocv.io/writing-code/,
// with some changes to compile in Windows.

package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"gocv.io/x/gocv"
)

const MinimumArea = 3000

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\t go run hand-gestures.go [camera ID]") // camera ID, try channel 0 or 1.
		return
	}

	// parse args
	deviceID := os.Args[1]

	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("Hand Gestures")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	imgGrey := gocv.NewMat()
	defer imgGrey.Close()

	imgBlur := gocv.NewMat()
	defer imgBlur.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	hull := gocv.NewMat()
	defer hull.Close()

	defects := gocv.NewMat()
	defer defects.Close()

	purple := color.RGBA{155, 0, 255, 0} // Sets Purple as the color default for the dots

	fmt.Printf("Start reading device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// cleaning up image
		gocv.CvtColor(img, &imgGrey, gocv.ColorBGRToGray)
		gocv.GaussianBlur(imgGrey, &imgBlur, image.Pt(35, 35), 0, 0, gocv.BorderDefault)
		gocv.Threshold(imgBlur, &imgThresh, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)

		// now find biggest contour
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		c := getBiggestContour(contours)

		gocv.ConvexHull(c, &hull, true, false)
		gocv.ConvexityDefects(c, hull, &defects)

		var angle float64
		defectCount := 0
		for i := 0; i < defects.Rows(); i++ {
			start := c[defects.GetIntAt(i, 0)]
			end := c[defects.GetIntAt(i, 1)]
			far := c[defects.GetIntAt(i, 2)]

			a := math.Sqrt(math.Pow(float64(end.X-start.X), 2) + math.Pow(float64(end.Y-start.Y), 2))
			b := math.Sqrt(math.Pow(float64(far.X-start.X), 2) + math.Pow(float64(far.Y-start.Y), 2))
			c := math.Sqrt(math.Pow(float64(end.X-far.X), 2) + math.Pow(float64(end.Y-far.Y), 2))

			// apply cosine rule here
			angle = math.Acos((math.Pow(b, 2)+math.Pow(c, 2)-math.Pow(a, 2))/(2*b*c)) * 57

			// ignore angles > 90 and highlight rest with dots
			if angle <= 90 {
				defectCount++
				gocv.Circle(&img, far, 1, purple, 2)
			}
		}

		// add 1 to count because que number of fingers is number of dots (light identified areas) plus one.
		status := fmt.Sprintf("defectCount: %d", defectCount+1)
		rect := gocv.BoundingRect(c)
		gocv.Rectangle(&img, rect, color.RGBA{255, 255, 255, 0}, 2)

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, purple, 2)

		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}
	}
}

func getBiggestContour(contours [][]image.Point) []image.Point {
	var area float64
	index := 0
	for i, c := range contours {
		newArea := gocv.ContourArea(c)
		if newArea > area {
			area = newArea
			index = i
		}
	}
	return contours[index]
}

// O que a aplicacao faz:
//
// Este programa detecta o numero de dedos levantados em frente a camera.
//
// 		go run hand_gestures 0
