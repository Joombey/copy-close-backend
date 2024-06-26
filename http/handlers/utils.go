package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

type ChatMessageRequest struct {
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	AuthToken string `json:"auth_token"`
}

type onChunk = func(chunk []byte, breakFunc func())

func chunkedReader(reader io.Reader, onChunk onChunk) {
	chunk := make([]byte, 8192)
	eof := false
	for !eof {
		r, err := reader.Read(chunk)
		if err != nil && err != io.EOF {
			return
		} else if err == io.EOF || r == 0 {
			eof = true
		}

		onChunk(chunk[:r], func() {
			eof = true
		})
	}
}

// func chunkedFile(file *os.File, onChunk onChunk) {
// 	chunk := make([]byte, 8192)
// 	eof := false
// 	for !eof {
// 		r, err := file.Read(chunk)
// 		if err != nil && err != io.EOF {
// 			return
// 		} else if err == io.EOF {
// 			eof = true
// 		}

// 		onChunk(chunk[:r], func() {
// 			eof = true
// 		})
// 	}
// }

func fromString(value string, receiver any) error {
	err := json.Unmarshal([]byte(value), receiver)
	if err != nil {
		return err
	}
	return nil
}

func getPathForJPEG(path string) string {
	absolutePath, _ := os.Getwd()
	return fmt.Sprintf("%s/files/%s.jpeg", absolutePath, path)
}

func getPathForName(path string) string {
	absolutePath, _ := os.Getwd()
	return fmt.Sprintf("%s/files/%s", absolutePath, path)
}

func stringToInt64(str string) (int64, error) {
	n, e := strconv.Atoi(str)
	return int64(n), e
}

type WorkerPool struct {
	mutex     sync.Mutex
	workerMap map[string][]*Worker
}

type Worker struct {
	work    work
	trigger chan any
	quit    chan any
}

type work = func()

func (w *Worker) Stop() { w.quit <- 1 }

func (p *WorkerPool) append(workId string, worker *Worker) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.workerMap[workId]) == 0 {
		p.workerMap[workId] = make([]*Worker, 0)
	}
	p.workerMap[workId] = append(p.workerMap[workId], worker)
}

func (p *WorkerPool) removeWorker(workId string, worker *Worker) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.workerMap[workId]) == 0 {
		return
	}

	for i, v := range p.workerMap[workId] {
		if v == worker {
			p.workerMap[workId] = append(p.workerMap[workId][:i], p.workerMap[workId][i+1:]...)
			return
		}
	}
}

func (p *WorkerPool) trigger(workId string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.workerMap[workId]) == 0 {
		return
	}

	for _, worker := range p.workerMap[workId] {
		worker.trigger <- 1
	}
}

var workerPool = &WorkerPool{
	mutex:     sync.Mutex{},
	workerMap: make(map[string][]*Worker),
}
