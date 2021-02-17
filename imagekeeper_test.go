package main

import (
	"fmt"
	"image"
	"testing"
)

func TestEmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	v := img.At(10, 10)
	fmt.Println(v)
}
