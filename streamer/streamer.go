package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// go get -u github.com/hybridgroup/mjpeg
	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

var (
	deviceID int
	err      error
	webcam   *gocv.VideoCapture
	stream   *mjpeg.Stream
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\t go run mjpeg-streamer.go [camera ID] [host:port]")
		return
	}

	// parse args
	deviceID := os.Args[1]
	host := os.Args[2]

	// open webcam
	webcam, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	// create the mjpeg stream
	stream = mjpeg.NewStream()

	// start capturing
	go mjpegCapture()

	fmt.Println("Capturing. Point your browser to " + host)

	// start http server
	http.Handle("/", stream)
	log.Fatal(http.ListenAndServe(host, nil))
}

func mjpegCapture() {
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}

// O que essa aplicacao faz:
//
// Esse exemplo captura imagens do dispositivo configurado e depois faz o MJPEG stream dessa captura.
// Apos executar o programa, aponte o seu navegador para o hostname/porta que você deseja, por exemplo (http://localhost:8080)
// entao devera ver o streamming nesse host.
//
// Como executar:
//
// go run streamer.go [camera ID] [host:port]
//
//		go get -u github.com/hybridgroup/mjpeg
// 		go run streamer.go 1 0.0.0.0:8080
//
// Este exemplo contem algumas modificacoes, mas teve como base o exemplo do tutorial existente em:
// https://gocv.io/writing-code/
