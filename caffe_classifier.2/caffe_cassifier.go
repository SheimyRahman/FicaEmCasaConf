package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"os"

	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Para executar:\ncaffe-classifier [camera ID] [caffe_model] [arq_config.prototxt] [arq_palavras.txt]")
		return
	}

	// parse args
	deviceID := os.Args[1]
	model := os.Args[2]
	config := os.Args[3]
	descr := os.Args[4]
	descriptions, err := readDescriptions(descr)
	if err != nil {
		fmt.Printf("Error reading descriptions file: %v\n", descr)
		return
	}

	backend := gocv.NetBackendDefault
	if len(os.Args) > 5 {
		backend = gocv.ParseNetBackend(os.Args[5])
	}

	target := gocv.NetTargetCPU
	if len(os.Args) > 6 {
		target = gocv.ParseNetTarget(os.Args[6])
	}

	// open capture device
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Caffe Classifier")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// open DNN classifier
	net := gocv.ReadNet(model, config)
	if net.Empty() {
		fmt.Printf("Error reading network model from : %v %v\n", model, config)
		return
	}
	defer net.Close()
	net.SetPreferableBackend(gocv.NetBackendType(backend))
	net.SetPreferableTarget(gocv.NetTargetType(target))

	status := "Ready"
	statusColor := color.RGBA{100, 0, 255, 0}
	fmt.Printf("Start reading device: %v\n", deviceID)

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// convert image Mat to 224x224 blob that the classifier can analyze
		blob := gocv.BlobFromImage(img, 1.0, image.Pt(224, 224), gocv.NewScalar(104, 117, 123, 0), false, false)

		// feed the blob into the classifier
		net.SetInput(blob, "")

		// run a forward pass thru the network
		prob := net.Forward("")

		// reshape the results into a 1x1000 matrix
		probMat := prob.Reshape(1, 1)

		// determine the most probable classification
		_, maxVal, _, maxLoc := gocv.MinMaxLoc(probMat)

		// display classification
		status = fmt.Sprintf("description: %v, maxVal: %v\n", descriptions[maxLoc.X], maxVal)
		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		blob.Close()
		prob.Close()
		probMat.Close()

		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

// readDescriptions reads the descriptions from a file
// and returns a slice of its lines.
func readDescriptions(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Traducao para pt-br do original[1]

// O que essa aplicacao faz:
//
// Esse exemplo utiliza o deep learning framework Caffe (http://caffe.berkeleyvision.org/)
// para classificar o que estiver na frente da camera.
//
// Faca o Download do Caffe model atraves da seguinte URL:
// http://dl.caffe.berkeleyvision.org/bvlc_googlenet.caffemodel
//
// Voce tambem deve fazer o download do arquivo ".prototxt", que pode ser encontrado aqui:
// https://raw.githubusercontent.com/opencv/opencv_extra/master/testdata/dnn/bvlc_googlenet.prototxt
//
// Para executar:
//
// 		go run ./cmd/caffe-classifier/main.go 0 ./cmd/bvlc_googlenet.caffemodel ./cmd/bvlc_googlenet.prototxt ./cmd/Downloads/classification_classes.txt
//
// [1] Fonte Original: https://github.com/hybridgroup/gocv
