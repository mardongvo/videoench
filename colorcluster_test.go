package main

import (
	"image/color"
	"testing"
)

func TestSort(t *testing.T) {
	cc := ColorCluster{}
	cc.arr = append(cc.arr, ColorItem{distances: []float64{2}}) //0
	cc.arr = append(cc.arr, ColorItem{distances: []float64{1}}) //1
	cc.arr = append(cc.arr, ColorItem{distances: []float64{3}}) //2
	cc.arr = append(cc.arr, ColorItem{distances: []float64{2}}) //3
	distArr := cc.getCenterAndTail()
	if (distArr[0].index != 1) || (distArr[len(distArr)-1].index != 2) {
		t.Fatalf("getCenterAndTail: unexpected %v and %v", distArr[0].index, distArr[len(distArr)-1].index)
	}
}

func TestLAB(t *testing.T) {
	var ci ColorItem
	ci = NewColorItem(color.RGBA{0, 0, 0, 0})
	if (ci.lab.l != 0) || (ci.lab.a != 0) || (ci.lab.b != 0) {
		t.Fatalf("colorLAB: wrong %#v <-> %#v", ci.source, ci.lab)
	}
	/*
		ci = NewColorItem(color.RGBA{255, 255, 255, 0})
		if (ci.lab.l != 100) || (ci.lab.a != 0) || (ci.lab.b != 0) {
			t.Fatalf("colorLAB: wrong %#v <-> %#v", ci.source, ci.lab)
		}
		ci = NewColorItem(color.RGBA{255, 0, 0, 0})
		if (ci.lab.l != 54.29) || (ci.lab.a != 80.81) || (ci.lab.b != 69.89) {
			t.Fatalf("colorLAB: wrong %#v <-> %#v", ci.source, ci.lab)
		}
	*/
}
