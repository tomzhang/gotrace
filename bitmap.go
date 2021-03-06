package gotrace

import (
	"unsafe"
)

var (
	wordSize = Word(unsafe.Sizeof(Word(0)))
	wordBits = wordSize * 8
	hiBit    = Word(1) << (wordBits - 1)
	allBits  = ^Word(0)
)

type Word uint

func NewBitmap(w, h int) *Bitmap {
	dy := 0
	if w != 0 {
		dy = (w-1)/int(wordBits) + 1
	}
	return &Bitmap{
		W: w, H: h,
		Map: make([]Word, dy*h), Dy: dy,
	}
}

//	Internal bitmap format. The n-th scanline starts at scanline(n) =
//	(map + n*dy). Raster data is stored as a sequence of potrace_words
//	(NOT bytes). The leftmost bit of scanline n is the most significant
//	bit of scanline(n)[0].
type Bitmap struct {
	W, H int    // width and height, in pixels
	Dy   int    // words per scanline (not bytes)
	Map  []Word // raw data, dy*h words
}

func (bm *Bitmap) scanline(y int) []Word { return bm.Map[y*bm.Dy : (y+1)*bm.Dy] }
func (bm *Bitmap) index(x, y int) *Word  { return &(bm.Map[Word(y*bm.Dy)+Word(x)/wordBits]) }
func (bm *Bitmap) mask(x int) Word       { return hiBit >> (Word(x) & (wordBits - 1)) }
func (bm *Bitmap) Get(x, y int) bool {
	if x >= 0 && x < bm.W && y >= 0 && y < bm.H {
		return ((*bm.index(x, y)) & bm.mask(x)) != 0
	}
	return false
}
func (bm *Bitmap) Set(x, y int, v bool) {
	if x >= 0 && x < bm.W && y >= 0 && y < bm.H {
		if v {
			*bm.index(x, y) |= bm.mask(x)
		} else {
			*bm.index(x, y) &= ^bm.mask(x)
		}
	}
}

// clear the given bitmap. Set all bits to c.
func (bm *Bitmap) Clear(c bool) {
	if !c {
		for i := range bm.Map {
			bm.Map[i] = 0
		}
	} else {
		// TODO: change to set allBits in Map?
		for y := 0; y < bm.H; y++ {
			for x := 0; x < bm.W; x++ {
				bm.Set(x, y, c)
			}
		}
	}
}

// duplicate the given bitmap
func (b *Bitmap) Clone() *Bitmap {
	b2 := NewBitmap(b.W, b.H)
	copy(b2.Map, b.Map)
	return b2
}
