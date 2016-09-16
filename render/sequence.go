package render

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"
)

type Sequence struct {
	Point
	Layered
	rs         []Modifiable
	lastChange time.Time
	playing    bool
	sheetPos   int
	frameTime  int64
	cID        event.CID
}

func NewSequence(mods []Modifiable, fps float64) *Sequence {
	return &Sequence{
		Point: Point{
			X: 0.0,
			Y: 0.0,
		},
		sheetPos:   0,
		frameTime:  int64(math.Pow(10, 9) / fps),
		rs:         mods,
		lastChange: time.Now(),
		playing:    true,
	}
}

func (sq *Sequence) Copy() Modifiable {
	newSq := new(Sequence)
	*newSq = *sq

	newRs := make([]Modifiable, len(sq.rs))
	for i, r := range sq.rs {
		newRs[i] = r.Copy()
	}

	newSq.rs = newRs
	return newSq
}

func (sq *Sequence) SetTriggerID(id event.CID) {
	sq.cID = id
}

func (sq *Sequence) update() {
	if sq.playing && time.Since(sq.lastChange).Nanoseconds() > sq.frameTime {
		sq.lastChange = time.Now()
		sq.sheetPos = (sq.sheetPos + 1) % len(sq.rs)
		if sq.sheetPos == 0 && sq.cID != 0 {
			sq.cID.Trigger("AnimationEnd", nil)
		}
	}
}

func (sq *Sequence) Draw(buff draw.Image) {
	sq.update()
	img := sq.GetRGBA()
	ShinyDraw(buff, img, int(sq.X), int(sq.Y))
}

func (sq *Sequence) GetRGBA() *image.RGBA {
	return sq.rs[sq.sheetPos].GetRGBA()
}

func (sq *Sequence) ApplyColor(c color.Color) Modifiable {
	for _, r := range sq.rs {
		r.ApplyColor(c)
	}
	return sq
}

func (sq *Sequence) FillMask(img image.RGBA) Modifiable {
	for _, r := range sq.rs {
		r.FillMask(img)
	}
	return sq
}

func (sq *Sequence) ApplyMask(img image.RGBA) Modifiable {
	for _, r := range sq.rs {
		r.ApplyMask(img)
	}
	return sq
}

func (sq *Sequence) Rotate(degrees int) Modifiable {
	for _, r := range sq.rs {
		r.Rotate(degrees)
	}
	return sq
}

func (sq *Sequence) Scale(xRatio float64, yRatio float64) Modifiable {
	for _, r := range sq.rs {
		r.Scale(xRatio, yRatio)
	}
	return sq
}

func (sq *Sequence) FlipX() Modifiable {
	for _, r := range sq.rs {
		r.FlipX()
	}
	return sq
}

func (sq *Sequence) FlipY() Modifiable {
	for _, r := range sq.rs {
		r.FlipY()
	}
	return sq
}
func (sq *Sequence) Fade(alpha int) Modifiable {
	for _, r := range sq.rs {
		r.Fade(alpha)
	}
	return sq
}

func (sq *Sequence) Pause() {
	sq.playing = false
}

func (sq *Sequence) Unpause() {
	sq.playing = true
}