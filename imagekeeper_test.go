package main

import (
	"fmt"
	"image"
	"reflect"
	"testing"
)

func TestEmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	v := img.At(10, 10)
	fmt.Println(v)
}

func TestAround(t *testing.T) {
	ik := NewImageKeeper("")
	width := 1
	ik.filelist = []string{"0", "1", "2", "3", "4", "5", "6"}
	for _, test := range []struct {
		inp string
		out []string
	}{
		{
			"0",
			[]string{"1", "2"},
		},
		{
			"1",
			[]string{"0", "2"},
		},
		{
			"6",
			[]string{"4", "5"},
		},
	} {
		out := ik.getAround(test.inp, width)
		if !reflect.DeepEqual(test.out, out) {
			t.Fatalf("not equal:: inp: >%v< out>%v< expect>%v<", test.inp, out, test.out)
		}
	}
}

func TestAround2(t *testing.T) {
	ik := NewImageKeeper("")
	width := 2
	ik.filelist = []string{"0", "1", "2", "3", "4", "5", "6"}
	for _, test := range []struct {
		inp string
		out []string
	}{
		{
			"0",
			[]string{"1", "2", "3", "4"},
		},
		{
			"1",
			[]string{"0", "2", "3", "4"},
		},
		{
			"6",
			[]string{"2", "3", "4", "5"},
		},
	} {
		out := ik.getAround(test.inp, width)
		if !reflect.DeepEqual(test.out, out) {
			t.Fatalf("not equal:: inp: >%v< out>%v< expect>%v<", test.inp, out, test.out)
		}
	}
}
