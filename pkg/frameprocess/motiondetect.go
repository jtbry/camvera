package frameprocess

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

type motionDetector struct {
	delta  gocv.Mat
	thresh gocv.Mat
	window *gocv.Window
	mog2   gocv.BackgroundSubtractorMOG2
}

func NewMotionDetector() FrameProcessor {
	return &motionDetector{
		delta:  gocv.NewMat(),
		thresh: gocv.NewMat(),
		window: gocv.NewWindow("Motion Detect"),
		mog2:   gocv.NewBackgroundSubtractorMOG2(),
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

	for i := 0; i < contours.Size(); i++ {
		area := gocv.ContourArea(contours.At(i))
		if area < 3000 {
			continue
		}

		gocv.DrawContours(&frame, contours, i, color.RGBA{0, 0, 255, 0}, 2)
		gocv.Rectangle(&frame, gocv.BoundingRect(contours.At(i)), color.RGBA{255, 0, 0, 0}, 2)
	}

	contours.Close()
	md.window.IMShow(frame)
}

func (md *motionDetector) Close() {
	md.delta.Close()
	md.thresh.Close()
	md.window.Close()
	md.mog2.Close()
}
