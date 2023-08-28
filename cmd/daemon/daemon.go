package main

import (
	"github.com/jtbry/camvera/pkg/frameprocess"
	"gocv.io/x/gocv"
)

func main() {
	cam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		panic(err)
	}
	defer cam.Close()

	img := gocv.NewMat()
	defer img.Close()

	md := frameprocess.NewMotionDetector()
	defer md.Close()
	for {
		if ok := cam.Read(&img); !ok {
			panic("cannot read device")
		}
		if img.Empty() {
			continue
		}

		md.ProcessFrame(img)
		gocv.WaitKey(1)
	}
}
