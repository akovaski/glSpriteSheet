package glSpriteSheet

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type ImagePacker struct {
	images  []sizedImage
	sprites []*Sprite
}

type sizedImage struct {
	w, h  int
	image image.Image
}

func (i *ImagePacker) AddFromFile(filename string) (*Sprite, error) {
	// Open the file
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image.
	m, err := png.Decode(file)
	if err != nil {
		fmt.Println("decoding error")
		return nil, err
	}
	bounds := m.Bounds()
	size := bounds.Max.Sub(bounds.Min)

	i.images = append(i.images, sizedImage{size.X, size.Y, m})
	i.sprites = append(i.sprites, new(Sprite))
	return i.sprites[len(i.sprites)-1], nil
}

func (i *ImagePacker) Pack() (SpriteSheet, error) {
	maxTexSize := getMaxTextureSize()
	var sheet SpriteSheet
	var err error

	for size := 16; size <= maxTexSize; size *= 2 {
		sheet, err = i.pack(size)
		fmt.Println(size, err)
		if err == nil {
			break
		}
	}

	if err != nil {
		return SpriteSheet{}, err
	}

	return sheet, err
}

// Packs image into a texture of size * size dimensions
func (i *ImagePacker) pack(size int) (SpriteSheet, error) {
	rootNode := newNode(size, size)

	for v, img := range i.images {
		err := rootNode.recInsert(img.w, img.h, v)
		if err != nil {
			return SpriteSheet{}, err
		}
	}

	nodeImage := image.NewRGBA(image.Rect(0, 0, size, size))

	traverseNodes(rootNode, func(nd node) {
		draw.Draw(nodeImage, image.Rect(nd.rc.left, nd.rc.top, nd.rc.right, nd.rc.bottom), i.images[nd.id].image, image.ZP, draw.Src)

		i.sprites[nd.id].left = float32(nd.rc.left) / float32(size)
		i.sprites[nd.id].top = float32(nd.rc.bottom) / float32(size)
		i.sprites[nd.id].right = float32(nd.rc.right) / float32(size)
		i.sprites[nd.id].bottom = float32(nd.rc.top) / float32(size)
		i.sprites[nd.id].W = float32(nd.rc.right - nd.rc.left)
		i.sprites[nd.id].H = float32(nd.rc.bottom - nd.rc.top)
	})

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, size, size, 0, gl.RGBA, gl.UNSIGNED_BYTE, nodeImage.Pix)

	texture.Unbind(gl.TEXTURE_2D)

	spriteSheet := NewSpriteSheet(texture, size, size)

	return spriteSheet, nil
}

func traverseNodes(root node, do func(node)) {
	// only 'do' if the node has been assigned an id (id == -1 if not assigned)
	if root.id > -1 {
		do(root)
	}
	for _, child := range root.child {
		traverseNodes(child, do)
	}
}

func getMaxTextureSize() int {
	maxTexSize := make([]int32, 1)
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, maxTexSize)
	return int(maxTexSize[0])
}

// Node image packing taken from http://www.blackpawn.com/texts/lightmaps/
type node struct {
	child []node
	rc    rectangle
	id    int
}

func newNode(width, height int) node {
	return node{nil, rectangle{width, height, 0, 0, width, height}, -1}
}

type rectangle struct {
	w, h, left, top, right, bottom int
}

func newRectangle(left, top, right, bottom int) rectangle {
	return rectangle{right - left, bottom - top, left, top, right, bottom}
}

// recursive insert operation, returns inserted node (nil if no insertion)
func (n *node) recInsert(w, h int, id int) error {
	if len(n.child) != 0 {
		// try inserting into first child
		if err := n.child[0].recInsert(w, h, id); err == nil {
			return err
		} else { // no room, insert into second
			return n.child[1].recInsert(w, h, id)
		}
	} else {
		// if there's already a lightmap here, return (-1 is id for empty)
		if n.id > -1 {
			return errors.New("Could not insert image")
		}

		// if we're too small, return
		if n.rc.w < w || n.rc.h < h {
			return errors.New("Could not insert image")
		}

		// if we're just right, accept
		if n.rc.w == w && n.rc.h == h {
			n.id = id
			return nil
		}

		n.child = make([]node, 2, 2)
		// otherwise, gotta split this node and create some kids
		n.child[0] = node{id: -1}
		n.child[1] = node{id: -1}

		// decide which way to split
		dw := n.rc.w - w
		dh := n.rc.h - h

		if dw > dh { // divide the node into left and right segments
			n.child[0].rc = newRectangle(n.rc.left, n.rc.top, n.rc.left+w, n.rc.bottom)
			n.child[1].rc = newRectangle(n.rc.left+w, n.rc.top, n.rc.right, n.rc.bottom)
		} else { // divide the node into top and bottom segments
			n.child[0].rc = newRectangle(n.rc.left, n.rc.top, n.rc.right, n.rc.top+h)
			n.child[1].rc = newRectangle(n.rc.left, n.rc.top+h, n.rc.right, n.rc.bottom)
		}

		// insert into the first child we created
		return n.child[0].recInsert(w, h, id)
	}
}
