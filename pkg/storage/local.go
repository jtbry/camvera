package storage

import (
	"fmt"
	"time"

	"gocv.io/x/gocv"
)

type LocalStorage struct {
	activeWriter *gocv.VideoWriter
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (ls *LocalStorage) SaveFrame(frame gocv.Mat) bool {
	path := fmt.Sprintf("./data/%d.jpg", time.Now().Unix())
	return gocv.IMWrite(path, frame)
}

func (ls *LocalStorage) OpenVideoWriter() error {
	if ls.activeWriter != nil {
		return fmt.Errorf("VideoWriter already active")
	}
	path := fmt.Sprintf("./data/%d.mp4", time.Now().Unix())
	writer, err := gocv.VideoWriterFile(path, "mp4v", 25, 640, 480, true)
	if err != nil {
		return err
	}
	ls.activeWriter = writer
	return nil
}

func (ls *LocalStorage) WriteVideoFrame(frame gocv.Mat) error {
	if ls.activeWriter == nil {
		return fmt.Errorf("VideoWriter not active")
	}
	return ls.activeWriter.Write(frame)
}

func (ls *LocalStorage) CloseVideoWriter() error {
	if ls.activeWriter == nil {
		return fmt.Errorf("VideoWriter not active")
	}
	return ls.activeWriter.Close()
}

func (ls *LocalStorage) Close() {
	if ls.activeWriter != nil {
		ls.activeWriter.Close()
	}
}
