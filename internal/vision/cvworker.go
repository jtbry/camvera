package vision

import (
	"context"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/jtbry/camvera/pkg/frameprocess"
	"gocv.io/x/gocv"
)

type cvWorker struct {
	ctx context.Context
	log *log.Logger
}

func NewCvWorker(c context.Context, logger *log.Logger) *cvWorker {
	return &cvWorker{
		ctx: c,
		log: logger,
	}
}

func (w *cvWorker) Start() {
	// Open video capture
	cam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		w.log.Fatal(err)
	}
	defer cam.Close()

	// Create frame processors
	md := frameprocess.NewMotionDetector(w.log)
	defer md.Close()
	streamer := frameprocess.NewWebStreamer(w.log)
	defer streamer.Close()

	// Read frames
	frame := gocv.NewMat()
	defer frame.Close()
	for w.ctx.Err() == nil {
		if ok := cam.Read(&frame); !ok {
			w.log.Println("cannot read device")
			continue
		}
		if frame.Empty() {
			continue
		}

		// Timestamp the frame
		now := time.Now()
		gocv.PutText(&frame, now.Format(time.DateTime), image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, color.RGBA{255, 255, 255, 0}, 2)

		// Apply frame processors
		md.ProcessFrame(frame)
		streamer.ProcessFrame(frame)

		// WaitKey
		gocv.WaitKey(1)
	}
}
