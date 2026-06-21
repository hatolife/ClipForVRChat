package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

const canvas = 1024

var sizes = []int{16, 24, 32, 48, 64, 128, 256, 512, 1024}

func main() {
	src := drawIcon(canvas)
	if err := os.MkdirAll("build/icons", 0755); err != nil {
		panic(err)
	}
	mustWrite("build/appicon.png", src)
	for _, size := range sizes {
		mustWrite(filepath.Join("build/icons", "appicon-"+itoa(size)+".png"), resize(src, size))
	}
}

func drawIcon(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	drawRoundedRect(img, 74, 74, 876, 876, 168, color.RGBA{210, 228, 252, 255})
	drawSoftSpot(img, 210, 180, 420, color.RGBA{245, 250, 255, 90}, image.Rect(74, 74, 950, 950), 168)
	drawSoftSpot(img, 835, 790, 380, color.RGBA{177, 204, 241, 90}, image.Rect(74, 74, 950, 950), 168)

	white := color.RGBA{255, 255, 255, 250}
	lavender := color.RGBA{180, 154, 218, 255}

	strokeRoundedRect(img, 230, 265, 540, 455, 48, 40, white)
	drawCircle(img, 360, 410, 56, white)
	fillPolygon(img, []point{{250, 700}, {390, 555}, {465, 630}, {590, 465}, {752, 705}}, white)

	strokeRoundedRect(img, 590, 560, 275, 270, 35, 34, lavender)
	drawRoundedRect(img, 650, 540, 160, 76, 32, lavender)
	drawCircle(img, 730, 540, 52, lavender)
	drawCircle(img, 730, 540, 20, white)
	drawRoundedRect(img, 630, 610, 195, 170, 12, white)

	return img
}

func drawRoundedRect(img *image.RGBA, x, y, w, h, r int, c color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if insideRoundedRect(px, py, x, y, w, h, r) {
				blend(img, px, py, c)
			}
		}
	}
}

func strokeRoundedRect(img *image.RGBA, x, y, w, h, r, stroke int, c color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if insideRoundedRect(px, py, x, y, w, h, r) &&
				!insideRoundedRect(px, py, x+stroke, y+stroke, w-stroke*2, h-stroke*2, max(1, r-stroke)) {
				blend(img, px, py, c)
			}
		}
	}
}

func drawSoftSpot(img *image.RGBA, cx, cy, radius int, c color.RGBA, clip image.Rectangle, clipRadius int) {
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			if !insideRoundedRect(x, y, clip.Min.X, clip.Min.Y, clip.Dx(), clip.Dy(), clipRadius) {
				continue
			}
			d := math.Hypot(float64(x-cx), float64(y-cy))
			if d > float64(radius) {
				continue
			}
			cc := c
			cc.A = uint8(float64(c.A) * (1 - d/float64(radius)))
			blend(img, x, y, cc)
		}
	}
}

func drawCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if math.Hypot(float64(x-cx), float64(y-cy)) <= float64(r) {
				blend(img, x, y, c)
			}
		}
	}
}

type point struct {
	x int
	y int
}

func fillPolygon(img *image.RGBA, pts []point, c color.RGBA) {
	minY, maxY := pts[0].y, pts[0].y
	for _, p := range pts {
		minY = min(minY, p.y)
		maxY = max(maxY, p.y)
	}
	for y := minY; y <= maxY; y++ {
		for x := 0; x < canvas; x++ {
			if insidePolygon(float64(x), float64(y), pts) {
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

func resize(src image.Image, size int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	sb := src.Bounds()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sx0 := sb.Min.X + x*sb.Dx()/size
			sx1 := sb.Min.X + (x+1)*sb.Dx()/size
			sy0 := sb.Min.Y + y*sb.Dy()/size
			sy1 := sb.Min.Y + (y+1)*sb.Dy()/size
			var r, g, b, a uint32
			var count uint32
			for sy := sy0; sy < max(sy1, sy0+1); sy++ {
				for sx := sx0; sx < max(sx1, sx0+1); sx++ {
					cr, cg, cb, ca := src.At(sx, sy).RGBA()
					r += cr
					g += cg
					b += cb
					a += ca
					count++
				}
			}
			dst.SetRGBA(x, y, color.RGBA{
				R: uint8((r / count) >> 8),
				G: uint8((g / count) >> 8),
				B: uint8((b / count) >> 8),
				A: uint8((a / count) >> 8),
			})
		}
	}
	return dst
}

func mustWrite(path string, img image.Image) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
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
		A: uint8(float64(c.A) + float64(dst.A)*(1-a)),
	})
}

func insideRoundedRect(px, py, x, y, w, h, r int) bool {
	if w <= 0 || h <= 0 || px < x || px >= x+w || py < y || py >= y+h {
		return false
	}
	cx := clamp(px, x+r, x+w-r-1)
	cy := clamp(py, y+r, y+h-r-1)
	return math.Hypot(float64(px-cx), float64(py-cy)) <= float64(r)
}

func clamp(v, lo, hi int) int {
	return max(lo, min(hi, v))
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	var buf [16]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
