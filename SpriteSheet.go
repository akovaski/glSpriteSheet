package glSpriteSheet

import (
	"github.com/go-gl/gl"
)

type SpriteSheet struct {
	texture gl.Texture
	w, h    float32
}

type Sprite struct {
	// sprite position in texture
	top, bottom, left, right float32

	// sprite position in view
	X, Y, W, H float32
}

// texture and width and hight of texture in pixels (or whatever unit you want)
func NewSpriteSheet(texture gl.Texture, width, height int) SpriteSheet {
	sheet := SpriteSheet{}

	sheet.w = float32(width)
	sheet.h = float32(height)

	sheet.texture = texture

	return sheet
}

// returns a sprite that knows it position in the texture (params), but does not know it's
// properties of being viewed (X, Y), assume same width, height as wT, hT (may not be desired)
func (sheet SpriteSheet) GetSprite(xT, yT, wT, hT int) *Sprite {
	sprite := new(Sprite)

	sprite.top = 1.0 - float32(yT+hT)/sheet.h
	sprite.bottom = 1.0 - float32(yT)/sheet.h
	sprite.left = float32(xT) / sheet.w
	sprite.right = float32(xT+wT) / sheet.w

	sprite.W = float32(wT)
	sprite.H = float32(hT)

	return sprite
}

//Draws all the sprites in the supplied slice
func (sheet SpriteSheet) Draw(sprites []*Sprite) {
	gl.Enable(gl.TEXTURE_2D)
	sheet.texture.Bind(gl.TEXTURE_2D)

	for _, sprite := range sprites {
		gl.Begin(gl.TRIANGLE_STRIP)
		{
			gl.TexCoord2f(sprite.left, sprite.bottom)
			gl.Vertex2f(sprite.X, sprite.Y)

			gl.TexCoord2f(sprite.left, sprite.top)
			gl.Vertex2f(sprite.X, sprite.Y+sprite.H)

			gl.TexCoord2f(sprite.right, sprite.bottom)
			gl.Vertex2f(sprite.X+sprite.W, sprite.Y)

			gl.TexCoord2f(sprite.right, sprite.top)
			gl.Vertex2f(sprite.X+sprite.W, sprite.Y+sprite.H)
		}
		gl.End()
	}

	sheet.texture.Unbind(gl.TEXTURE_2D)
	gl.Disable(gl.TEXTURE_2D)
}

func (sheet SpriteSheet) MoveTextPos(sp *Sprite, xT, yT, wT, hT int) {
	sp.top = 1.0 - float32(yT+hT)/sheet.h
	sp.bottom = 1.0 - float32(yT)/sheet.h
	sp.left = float32(xT) / sheet.w
	sp.right = float32(xT+wT) / sheet.w

	sp.W = float32(wT)
	sp.H = float32(hT)
}
