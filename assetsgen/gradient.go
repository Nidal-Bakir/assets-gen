package assetsgen

import (
	"image"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

type GradientType string

const (
	LinearGradient GradientType = "linear"
	RadialGradient GradientType = "radial"
)

// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (gt GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := range len(gt) - 1 {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found:
	//
	// either we are before a key point, then use the first color
	if t <= gt[0].Pos {
		return gt[0].Col
	}
	// or we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Col
}

func createLinearGradient(colorsTable GradientTable, degree float64, w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	theta := degree * math.Pi / 180 // to radian
	ux, uy := math.Cos(theta), math.Sin(theta)

	W, H := float64(w), float64(h)
	// the four corners of the image
	corners := [][2]float64{
		{0, 0},
		{W, 0},
		{0, H},
		{W, H},
	}

	// find min/max dot‑product over the corners
	rMin, rMax := math.Inf(+1), math.Inf(-1)
	for _, c := range corners {
		r := c[0]*ux + c[1]*uy
		rMin = math.Min(rMin, r)
		rMax = math.Max(rMax, r)
	}

	for y := range h {
		for x := range w {
			// project (x,y) onto our direction vector
			r := float64(x)*ux + float64(y)*uy

			// normalize into [0…1]
			t := (r - rMin) / (rMax - rMin)

			// sample and paint
			c := colorsTable.GetInterpolatedColorFor(t)
			img.Set(x, y, c)
		}
	}

	return img
}

func createRadialGradient(colorsTable GradientTable, w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	center := image.Pt(w/2, h/2)
	r := (math.Min(float64(w), float64(h))) / 2

	// distance between two points: √( (x1-x2)^2 + (y1-y2)^2 )
	distance := func(p1, p2 image.Point) float64 {
		return math.Sqrt(math.Pow(float64(p2.X-p1.X), 2) + math.Pow(float64(p2.Y-p1.Y), 2))
	}

	for y := range h {
		for x := range w {
			distance := math.Abs(distance(center, image.Pt(x, y)))

			// normalize into [0…1]
			t := distance / r

			// sample and paint
			c := colorsTable.GetInterpolatedColorFor(t)
			img.Set(x, y, c)
		}
	}

	return img
}
