package main

import (
	"image"
	"image/color"
	"time"

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

		// Timestamp the frame
		now := time.Now()
		gocv.PutText(&img, now.Format("2006-01-02 15:04:05"), image.Pt(10, 50), gocv.FontHersheyPlain, 2.0, color.RGBA{0, 255, 0, 0}, 2)

		md.ProcessFrame(img)
		gocv.WaitKey(1)
	}
}
