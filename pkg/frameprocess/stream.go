package frameprocess

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

type webStreamer struct {
	log     *log.Logger
	clients map[chan []byte]bool
	lock    sync.Mutex
	frame   []byte
}

const boundaryWord = "MJPEGBOUNDARY"
const headerf = "\r\n" + "--" + boundaryWord + "\r\n" + "Content-Type: image/jpeg\r\n" + "Content-Length: %d\r\n" + "X-Timestamp: 0.000000\r\n" + "\r\n"

func NewWebStreamer(logger *log.Logger) FrameProcessor {
	ws := &webStreamer{
		log:     logger,
		clients: make(map[chan []byte]bool),
	}

	http.HandleFunc("/", ws.serve)
	go http.ListenAndServe(":8080", nil)

	return ws
}

func (ws *webStreamer) ProcessFrame(frame gocv.Mat) {
	jpgBuffer, err := gocv.IMEncode(".jpg", frame)
	if err != nil {
		ws.log.Println("Error encoding frame")
		return
	}
	jpg := jpgBuffer.GetBytes()

	header := fmt.Sprintf(headerf, len(jpg))
	if len(ws.frame) < len(jpg)+len(header) {
		ws.frame = make([]byte, len(jpg)+len(header))
	}

	copy(ws.frame, header)
	copy(ws.frame[len(header):], jpg)

	ws.lock.Lock()
	for c := range ws.clients {
		c <- ws.frame
	}
	ws.lock.Unlock()
}

func (ws *webStreamer) Close() {
	for k := range ws.clients {
		delete(ws.clients, k)
	}
}

func (ws *webStreamer) serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundaryWord)

	c := make(chan []byte)
	ws.lock.Lock()
	ws.clients[c] = true
	ws.lock.Unlock()

	for {
		time.Sleep(50 * time.Millisecond)
		b := <-c
		_, err := w.Write(b)
		if err != nil {
			break
		}
	}

	ws.lock.Lock()
	delete(ws.clients, c)
	ws.lock.Unlock()
}
