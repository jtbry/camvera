package main

import (
	"gocv.io/x/gocv"
)

func main() {
	cam, err := gocv.OpenVideoCapture("http://pendelcam.kip.uni-heidelberg.de/mjpg/video.mjpg")
	if err != nil {
		panic(err)
	}
	defer cam.Close()

	window := gocv.NewWindow("camvera")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()
	for {
		cam.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}
}
