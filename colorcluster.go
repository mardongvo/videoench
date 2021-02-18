package main

import (
	"image/color"
	"math"
)

//incapsulate color value
type ColorItem struct {
	source color.RGBA //value from source image
	lab    struct {   //source in Lab colorspace
		l float64
		a float64
		b float64
	}
	distances []float64 //distances between colors in cluster
	distsum   float64
}

func NewColorItem(source color.RGBA) ColorItem {
	var res ColorItem
	res.source = source
	res.CalcLAB()
	return res
}

func (c *ColorItem) CalcLAB() {
	// src: https://gist.github.com/manojpandey/f5ece715132c572c80421febebaf66ae
	// http://www.easyrgb.com/en/math.php
	var f_r float64 = float64(c.source.R) / 255.0
	var f_g float64 = float64(c.source.G) / 255.0
	var f_b float64 = float64(c.source.B) / 255.0

	// 1 RGB->XYZ

	if f_r > 0.04045 {
		f_r = math.Pow((f_r+0.055)/1.055, 2.4)
	} else {
		f_r = f_r / 12.92
	}

	if f_g > 0.04045 {
		f_g = math.Pow((f_g+0.055)/1.055, 2.4)
	} else {
		f_g = f_g / 12.92
	}

	if f_b > 0.04045 {
		f_b = math.Pow((f_b+0.055)/1.055, 2.4)
	} else {
		f_b = f_b / 12.92
	}

	f_r *= 100.0
	f_g *= 100.0
	f_b *= 100.0

	X := f_r*0.4124 + f_g*0.3576 + f_b*0.1805
	Y := f_r*0.2126 + f_g*0.7152 + f_b*0.0722
	Z := f_r*0.0193 + f_g*0.1192 + f_b*0.9505

	// 2 XYZ - Lab
	X = X / 95.047
	Y = Y / 100.0
	Z = Z / 108.883

	if X > 0.008856 {
		X = math.Pow(X, 1/3.0)
	} else {
		X = (7.787 * X) + (16.0 / 116.0)
	}

	if Y > 0.008856 {
		Y = math.Pow(Y, 1/3.0)
	} else {
		Y = (7.787 * Y) + (16.0 / 116.0)
	}

	if Z > 0.008856 {
		Z = math.Pow(Z, 1/3.0)
	} else {
		Z = (7.787 * Z) + (16.0 / 116.0)
	}

	c.lab.l = (116 * Y) - 16
	c.lab.a = 500 * (X - Y)
	c.lab.b = 200 * (Y - Z)

	c.lab.l = math.Round(c.lab.l*100) / 100.0
	c.lab.a = math.Round(c.lab.a*100) / 100.0
	c.lab.b = math.Round(c.lab.b*100) / 100.0
}

func distance(c1, c2 ColorItem) float64 {
	return math.Sqrt(math.Pow(c1.lab.l-c2.lab.l, 2) + math.Pow(c1.lab.a-c2.lab.a, 2) + math.Pow(c1.lab.b-c2.lab.b, 2))
}

type ColorCluster struct {
	arr []ColorItem
}

func (cc *ColorCluster) add(rgba color.RGBA) {
	ci := NewColorItem(rgba)
	if len(cc.arr) > 0 {
		for i, _ := range cc.arr {
			d := distance(cc.arr[i], ci)
			cc.arr[i].distances = append(cc.arr[i].distances, d)
			ci.distances = append(ci.distances, d)
		}
	}
	cc.arr = append(cc.arr, ci)
}

type distItem struct {
	index int64
	dist  float64
}

func (cc *ColorCluster) sumDist() {
	for i, ci := range cc.arr {
		var d float64
		for _, s := range ci.distances {
			d += s
		}
		cc.arr[i].distsum = d
	}
}

func (cc *ColorCluster) getCenterAndTail() []distItem {
	var distArr []distItem = make([]distItem, len(cc.arr), len(cc.arr))
	cc.sumDist()
	for i, _ := range distArr {
		distArr[i].index = -1
		distArr[i].dist = math.MaxFloat64
	}
	for i, ci := range cc.arr {
		di := distItem{int64(i), ci.distsum}
		//fill-sort distance array
		for j, _ := range distArr {
			if distArr[j].dist > di.dist {
				di, distArr[j] = distArr[j], di
			}
		}
	}
	return distArr
}

func (cc *ColorCluster) decide(relation float64) color.RGBA {
	distArr := cc.getCenterAndTail()
	if distArr[0].index == 0 {
		return cc.arr[0].source
	}
	d := cc.arr[0].distsum
	d0 := distArr[0].dist
	d1 := distArr[len(distArr)-1].dist
	//if first element of cluster is near to center - do nothing
	if d < d0+(d1-d0)*relation {
		return cc.arr[0].source
	}
	//if first element too far - change to cluster center
	return cc.arr[distArr[0].index].source
}
