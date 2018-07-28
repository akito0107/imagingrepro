package main

import (
	"github.com/disintegration/imaging"
	"os"
	"image"
	_ "net/http/pprof"
	"io"
	"image/jpeg"
	"io/ioutil"
	"flag"
	"net/http"
	"log"
	"sync"
	time2 "time"
	"runtime"
)

var count = flag.Int("count", 100, "")
var concurrency = flag.Int("concurency", 10, "")
var throughput = flag.Int("throughput", 5, "")

func main() {
	runtime.GOMAXPROCS(4)
	start := time2.Now()
	defer func() {
		finish := time2.Now()
		log.Printf("start %v, finish %v \n", start, finish)
	}()

	flag.Parse()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	filepath := "test.jpg"
	var wg sync.WaitGroup

	queue := make(chan int)
	lock := make(chan struct{}, *concurrency)

	go func() {
		for c := range queue {
			lock <- struct{}{}
			wg.Add(1)
			go func(c int) {
				log.Printf("start %d \n", c)
				process(filepath)
				<-lock
				wg.Done()
			}(c)
		}
	}()
	lim := *count
	for c := 0; c < lim; c++ {
		queue <- c
	}
	close(queue)
	wg.Wait()
}

func process(path string) error {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}
	i, err := decode(f)

	if err != nil {
		return err
	}
	img := imaging.Resize(i, 200, 200, imaging.Lanczos)
	return write(ioutil.Discard, img)
}

func decode(r io.Reader) (image.Image, error) {
	i, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return i, err
}

func write(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
}