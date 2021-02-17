package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path"
	"sync"
)

type ImageRecord struct {
	source   image.Image
	usecount int64
}

type ImageKeeper struct {
	path      string                 //path to directory with images
	filelist  []string               //ordered list of files to find neightbors
	resources map[string]ImageRecord //allocated resources
	mtx       sync.Mutex
}

func NewImageKeeper(path string) ImageKeeper {
	var res ImageKeeper
	res.path = path
	//TODO: list path
	res.resources = make(map[string]ImageRecord)
	return res
}

func (ik *ImageKeeper) sieze(filename string) image.Image {
	ik.mtx.Lock()
	defer ik.mtx.Unlock()
	ir, ok := ik.resources[filename]
	if !ok {
		ik.resources[filename] = NewImageRecord(path.Join(ik.path, filename))
		ir = ik.resources[filename]
	}
	ir.usecount++
	return ir.source
}

func (ik *ImageKeeper) release(filename string) {
	ik.mtx.Lock()
	defer ik.mtx.Unlock()
	ir, ok := ik.resources[filename]
	if !ok {
		//??
		return
	}
	ir.usecount--
	if ir.usecount <= 0 {
		delete(ik.resources, filename)
	}
}

func NewImageRecord(fullpath string) ImageRecord {
	var res ImageRecord
	rdr, err := os.Open(fullpath)
	if err != nil {
		fmt.Printf("open image error(%s): %v\n", fullpath, err)
		return res
	}
	defer rdr.Close()
	srcimg, _, err := image.Decode(rdr)
	if err != nil {
		fmt.Printf("read image error(%s): %v\n", fullpath, err)
		return res
	}
	res.source = srcimg
	return res
}

func (ik ImageKeeper) getAround(filename string, width int) []string {
	var pos int = -1
	var lwidth int = width
	var rwidth int = width
	var res []string
	if width*2+1 > len(ik.filelist) {
		return res
	}
	for i, v := range ik.filelist {
		if v == filename {
			pos = i
			break
		}
	}
	if pos == -1 {
		return res
	}
	ldif := 0 - (pos - lwidth)
	if ldif > 0 {
		lwidth -= ldif
		rwidth += ldif
	}
	rdif := (pos + rwidth) - (len(ik.filelist) - 1)
	if rdif > 0 {
		rwidth -= rdif
		lwidth += rdif
	}
	for i := pos - lwidth; i < pos; i++ {
		res = append(res, ik.filelist[i])
	}
	for i := pos + 1; i <= pos+rwidth; i++ {
		res = append(res, ik.filelist[i])
	}
	return res
}
