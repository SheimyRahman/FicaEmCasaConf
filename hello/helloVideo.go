package main

import (
	"gocv.io/x/gocv"
)

func main() {

	// open webcam
	webcam, _ := gocv.VideoCaptureDevice(0)

	// open display window
	window := gocv.NewWindow("Hello, #FicaEmCasaConf 🏠!")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		gocv.WaitKey(1)
	}
}

// Traducao para pt-br do original[1]
// Este exemplo contem algumas modificacoes, mas teve como base o exemplo do tutorial existente em:
// [1] https://gocv.io/writing-code/
//
// Para executar:
//
// go run helloVideo.go
