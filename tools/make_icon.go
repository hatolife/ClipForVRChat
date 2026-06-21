package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

const (
	size  = 1024
	scale = 4
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, size*scale, size*scale))
	drawGradient(img)

	drawRoundedRect(img, 120, 120, 784, 784, 150, color.RGBA{R: 255, G: 255, B: 255, A: 34})
	drawRoundedRect(img, 210, 195, 430, 545, 64, color.RGBA{R: 255, G: 255, B: 255, A: 238})
	drawRoundedRect(img, 255, 250, 340, 255, 34, color.RGBA{R: 28, G: 136, B: 186, A: 255})
	drawMountain(img)
	drawClip(img)
	drawVR(img)

	out := downsample(img)
	if err := os.MkdirAll("build", 0755); err != nil {
		panic(err)
	}
	file, err := os.Create(filepath.Join("build", "appicon.png"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, out); err != nil {
		panic(err)
	}
}

func drawGradient(img *image.RGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			nx := float64(x) / float64(b.Dx())
			ny := float64(y) / float64(b.Dy())
			r := uint8(20 + 30*nx + 10*ny)
			g := uint8(124 + 74*nx + 25*(1-ny))
			bl := uint8(158 + 58*(1-nx) + 28*ny)
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: bl, A: 255})
		}
	}
}

func drawRoundedRect(img *image.RGBA, x, y, w, h, r int, c color.RGBA) {
	x *= scale
	y *= scale
	w *= scale
	h *= scale
	r *= scale
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			cx := clamp(px, x+r, x+w-r-1)
			cy := clamp(py, y+r, y+h-r-1)
			if dist(px, py, cx, cy) <= float64(r) {
				blend(img, px, py, c)
			}
		}
	}
}

func drawMountain(img *image.RGBA) {
	fillPolygon(img, []point{{270, 630}, {425, 430}, {545, 630}}, color.RGBA{R: 42, G: 169, B: 164, A: 255})
	fillPolygon(img, []point{{405, 630}, {610, 365}, {765, 630}}, color.RGBA{R: 23, G: 111, B: 177, A: 255})
	drawCircle(img, 650, 320, 58, color.RGBA{R: 255, G: 214, B: 103, A: 255})
}

func drawClip(img *image.RGBA) {
	strokeRoundedRect(img, 245, 140, 255, 525, 86, 44, color.RGBA{R: 245, G: 250, B: 255, A: 255})
	strokeRoundedRect(img, 325, 225, 180, 365, 58, 38, color.RGBA{R: 38, G: 128, B: 174, A: 255})
}

func drawVR(img *image.RGBA) {
	drawThickLine(img, 660, 615, 790, 375, 44, color.RGBA{R: 255, G: 255, B: 255, A: 245})
	drawThickLine(img, 790, 375, 885, 615, 44, color.RGBA{R: 255, G: 255, B: 255, A: 245})
	drawThickLine(img, 885, 375, 885, 615, 44, color.RGBA{R: 255, G: 255, B: 255, A: 245})
	drawThickLine(img, 885, 375, 960, 375, 44, color.RGBA{R: 255, G: 255, B: 255, A: 245})
	drawCircle(img, 960, 450, 72, color.RGBA{R: 255, G: 255, B: 255, A: 245})
	drawCircle(img, 960, 450, 34, color.RGBA{R: 45, G: 151, B: 184, A: 255})
}

type point struct {
	x int
	y int
}

func fillPolygon(img *image.RGBA, pts []point, c color.RGBA) {
	minY, maxY := pts[0].y, pts[0].y
	for _, p := range pts {
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	for y := minY * scale; y <= maxY*scale; y++ {
		for x := 0; x < size*scale; x++ {
			if insidePolygon(float64(x)/scale, float64(y)/scale, pts) {
				blend(img, x, y, c)
			}
		}
	}
}

func insidePolygon(x, y float64, pts []point) bool {
	inside := false
	j := len(pts) - 1
	for i := range pts {
		xi, yi := float64(pts[i].x), float64(pts[i].y)
		xj, yj := float64(pts[j].x), float64(pts[j].y)
		if (yi > y) != (yj > y) && x < (xj-xi)*(y-yi)/(yj-yi)+xi {
			inside = !inside
		}
		j = i
	}
	return inside
}

func strokeRoundedRect(img *image.RGBA, x, y, w, h, r, stroke int, c color.RGBA) {
	ox := x * scale
	oy := y * scale
	ow := w * scale
	oh := h * scale
	or := r * scale
	ix := (x + stroke) * scale
	iy := (y + stroke) * scale
	iw := (w - stroke*2) * scale
	ih := (h - stroke*2) * scale
	ir := max(1, r-stroke) * scale
	for py := oy; py < oy+oh; py++ {
		for px := ox; px < ox+ow; px++ {
			if insideRoundedRect(px, py, ox, oy, ow, oh, or) && !insideRoundedRect(px, py, ix, iy, iw, ih, ir) {
				blend(img, px, py, c)
			}
		}
	}
}

func drawCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	cx *= scale
	cy *= scale
	r *= scale
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if dist(x, y, cx, cy) <= float64(r) {
				blend(img, x, y, c)
			}
		}
	}
}

func drawThickLine(img *image.RGBA, x1, y1, x2, y2, width int, c color.RGBA) {
	x1 *= scale
	y1 *= scale
	x2 *= scale
	y2 *= scale
	width *= scale
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	length := math.Hypot(dx, dy)
	for y := min(y1, y2) - width; y <= max(y1, y2)+width; y++ {
		for x := min(x1, x2) - width; x <= max(x1, x2)+width; x++ {
			t := ((float64(x-x1) * dx) + (float64(y-y1) * dy)) / (length * length)
			t = math.Max(0, math.Min(1, t))
			px := float64(x1) + t*dx
			py := float64(y1) + t*dy
			if math.Hypot(float64(x)-px, float64(y)-py) <= float64(width)/2 {
				blend(img, x, y, c)
			}
		}
	}
}

func downsample(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			var r, g, b, a uint32
			for yy := 0; yy < scale; yy++ {
				for xx := 0; xx < scale; xx++ {
					cr, cg, cb, ca := src.At(x*scale+xx, y*scale+yy).RGBA()
					r += cr
					g += cg
					b += cb
					a += ca
				}
			}
			div := uint32(scale * scale)
			dst.SetRGBA(x, y, color.RGBA{
				R: uint8((r / div) >> 8),
				G: uint8((g / div) >> 8),
				B: uint8((b / div) >> 8),
				A: uint8((a / div) >> 8),
			})
		}
	}
	return dst
}

func blend(img *image.RGBA, x, y int, c color.RGBA) {
	if !(image.Pt(x, y).In(img.Bounds())) {
		return
	}
	dst := img.RGBAAt(x, y)
	a := float64(c.A) / 255
	img.SetRGBA(x, y, color.RGBA{
		R: uint8(float64(c.R)*a + float64(dst.R)*(1-a)),
		G: uint8(float64(c.G)*a + float64(dst.G)*(1-a)),
		B: uint8(float64(c.B)*a + float64(dst.B)*(1-a)),
		A: 255,
	})
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func dist(x1, y1, x2, y2 int) float64 {
	return math.Hypot(float64(x1-x2), float64(y1-y2))
}

func insideRoundedRect(px, py, x, y, w, h, r int) bool {
	if w <= 0 || h <= 0 {
		return false
	}
	if px < x || px >= x+w || py < y || py >= y+h {
		return false
	}
	cx := clamp(px, x+r, x+w-r-1)
	cy := clamp(py, y+r, y+h-r-1)
	return dist(px, py, cx, cy) <= float64(r)
}
