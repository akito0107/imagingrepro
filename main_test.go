package main

import (
	"testing"
	"os"
	"image/jpeg"
	"github.com/disintegration/imaging"
)

func Benchmark_process(b *testing.B) {
	f := "test.jpg"
	for i := 0; i < b.N; i++ {
		process(f)
	}
}

func Benchmark_resize1(b *testing.B) {
	f, _ := os.Open("test.jpg")
	img, err := jpeg.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		imaging.Resize(img, 200, 200, imaging.Lanczos)
	}
}

func Benchmark_resize2(b *testing.B) {
	f, _ := os.Open("test.jpg")
	img, err := jpeg.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		imaging.Resize(img, 400, 400, imaging.Lanczos)
	}
}


func Benchmark_decode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := os.Open("test.jpg")
		decode(f)
	}
}