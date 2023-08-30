package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"gocv.io/x/gocv"
)

type LocalStorage struct {
	log *log.Logger
}

func NewLocalStorage(logger *log.Logger) *LocalStorage {
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err := os.Mkdir("./data/", os.ModePerm)
		if err != nil {
			logger.Fatal(err)
		}
	}

	return &LocalStorage{
		log: logger,
	}
}

func (ls *LocalStorage) SaveImage(img gocv.Mat) {
	now := time.Now()
	name := fmt.Sprintf("./data/%d.jpg", now.Unix())
	gocv.IMWrite(name, img)
}

func (ls *LocalStorage) OpenVideoWriter() (*gocv.VideoWriter, error) {
	name := fmt.Sprintf("./data/%d.mp4", time.Now().Unix())
	return gocv.VideoWriterFile(name, "mp4v", 25, 640, 480, true)
}
