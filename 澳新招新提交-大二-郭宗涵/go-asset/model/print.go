package model

import (
	"bytes"
	"github.com/pbnjay/pixfont"
	"image"
	"image/color"
	//"image/draw"
	"github.com/llgcode/draw2d/draw2dimg"
	"image/png"
)

type SeatsPrinter struct {
	Seats  []Seat
	Path   []Seat
	Weight int
	Height int

	Space int
}

func (s *SeatsPrinter) PrintBestPath() ([]byte, error) {
	img := s.printBackground()
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color.Black)
	gc.SetStrokeColor(color.Black)
	gc.SetLineWidth(1)

	x := 20
	y := 3
	for k := range s.Path {
		gc.BeginPath()
		gc.MoveTo(float64(x), float64(y))
		x = s.Path[k].X * s.Space
		y = s.Path[k].Y * s.Space
		gc.LineTo(float64(x), float64(y))
		gc.FillStroke()
		gc.Close()
	}

	buf := bytes.Buffer{}
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *SeatsPrinter) printBackground() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, s.Weight, s.Height))
	pixer := pixfont.DefaultFont
	pixer.DrawString(img, 0, 0, "door", color.Black)

	for k := range s.Seats {
		pixfont.DrawString(img, s.Seats[k].X*s.Space, s.Seats[k].Y*s.Space, s.Seats[k].SeatId, color.Black)
	}
	return img
}
