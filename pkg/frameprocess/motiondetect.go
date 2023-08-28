package frameprocess

import (
	"image"
	"image/color"
	"time"

	"gocv.io/x/gocv"
)

type motionEvent struct {
	lastMove time.Time
}

type motionDetector struct {
	delta  gocv.Mat
	thresh gocv.Mat
	mog2   gocv.BackgroundSubtractorMOG2
	event  *motionEvent
}

func NewMotionDetector() FrameProcessor {
	return &motionDetector{
		delta:  gocv.NewMat(),
		thresh: gocv.NewMat(),
		mog2:   gocv.NewBackgroundSubtractorMOG2(),
		event:  nil,
	}
}

func (md *motionDetector) ProcessFrame(frame gocv.Mat) {
	// obtain foreground only
	md.mog2.Apply(frame, &md.delta)

	// threshold to ignore background
	gocv.Threshold(md.delta, &md.thresh, 25, 255, gocv.ThresholdBinary)

	// dilate
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	gocv.Dilate(md.thresh, &md.thresh, kernel)
	kernel.Close()

	// find contours
	contours := gocv.FindContours(md.thresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	// check contour sizes
	motion := false
	for i := 0; i < contours.Size(); i++ {
		area := gocv.ContourArea(contours.At(i))
		if area < 3000 {
			// contour area too small to be motion
			continue
		}

		// motion found
		motion = true
		if md.event == nil {
			println("Motion Started")
		} else {
			md.event.lastMove = time.Now()
		}
		gocv.Rectangle(&frame, gocv.BoundingRect(contours.At(i)), color.RGBA{255, 0, 0, 0}, 2)
	}

	// If no motion continues for 5 seconds, consider it ended
	if !motion && md.event != nil && time.Since(md.event.lastMove) > 5*time.Second {
		println("Motion Ended")
		md.event = nil
	}

	contours.Close()
}

func (md *motionDetector) Close() {
	md.delta.Close()
	md.thresh.Close()
	md.mog2.Close()
}
