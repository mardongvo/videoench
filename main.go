package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path"
	"sync"
)

func getPixel(img image.Image, x, y int) color.RGBA {
	return color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
}

func buildImage(source image.Image, arounds []image.Image) image.Image {
	var res *image.RGBA = image.NewRGBA(source.Bounds())
	bnds := source.Bounds()
	for x := bnds.Min.X; x < bnds.Max.X; x++ {
		for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
			cc := ColorCluster{}
			cc.add(getPixel(source, x, y))
			for _, img := range arounds {
				cc.add(getPixel(img, x, y))
			}
			res.SetRGBA(x, y, cc.decide(0.5))
		}
	}
	return res
}

func processImage(ik *ImageKeeper, outdir string, chInputFiles <-chan string, wg *sync.WaitGroup) {
	for filename := range chInputFiles {
		fmt.Printf("processing %s\n", filename)
		source := ik.sieze(filename)
		aroundFilenames := ik.getAround(filename, 2)
		arounds := make([]image.Image, 0)
		for _, f := range aroundFilenames {
			arounds = append(arounds, ik.sieze(f))
		}
		newimg := buildImage(source, arounds)
		ik.release(filename)
		for _, f := range aroundFilenames {
			ik.release(f)
		}
		out, err := os.Create(path.Join(outdir, filename))
		if err != nil {
			fmt.Printf("create file error (%s): %v\n", filename, err)
			continue
		}
		err = png.Encode(out, newimg)
		if err != nil {
			out.Close()
			fmt.Printf("write file error (%s): %v\n", filename, err)
			continue
		}
		out.Close()
	}
	wg.Done()
}

const WORKERS = 3

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: videoench <inputdir> <outdir>")
		return
	}
	///
	inpdir, outdir := os.Args[1], os.Args[2]
	inplist, err := os.ReadDir(inpdir)
	if err != nil {
		fmt.Printf("error listing dir (%s): %v\n", inpdir, err)
		return
	}
	outlist, err := os.ReadDir(outdir)
	if err != nil {
		fmt.Printf("error listing dir (%s): %v\n", outdir, err)
		return
	}
	filelist := make([]string, 0)
	for _, f := range inplist {
		if f.IsDir() {
			continue
		}
		doadd := true
		for _, g := range outlist {
			if g.Name() == f.Name() {
				doadd = false
				break
			}
		}
		if doadd {
			filelist = append(filelist, f.Name())
		}
	}
	///
	ik := NewImageKeeper(inpdir)
	chInputFiles := make(chan string, WORKERS)
	wg := sync.WaitGroup{}
	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go processImage(ik, outdir, chInputFiles, &wg)
	}
	for _, f := range filelist {
		chInputFiles <- f
	}
	wg.Wait()
}
