package frameprocess

import "gocv.io/x/gocv"

type FrameProcessor interface {
	ProcessFrame(frame gocv.Mat)
	Close()
}
