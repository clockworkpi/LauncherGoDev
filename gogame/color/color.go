package color

//import "github.com/veandco/go-sdl2/sdl"


type Color struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

func (c *Color) ToHex() int {
	return int( c.A<< 24 | c.R << 16 | c.G << 8 | c.B )
}

func (c *Color) ToBytes() []byte {
	bytes := make([]byte,4)
	bytes[0] = byte(c.R)
	bytes[1] = byte(c.G)
	bytes[2] = byte(c.B)
	bytes[3] = byte(c.A)
	return bytes
}

func (c *Color) RGBA() (r, g, b, a uint32) {
  r = uint32(c.R)
  r |= r << 8
  g = uint32(c.G)
  g |= g << 8
  b = uint32(c.B)
  b |= b << 8
  a = uint32(c.A)
  a |= a << 8
  return
}
