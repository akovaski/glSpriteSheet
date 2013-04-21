glSpriteSheet
=============

Utility to make it easier to work with 2D textures in OpenGL in go.

* SpriteSheets can be used to store textures and draw 2D 'sprites'.
* Sprites contain info on what part of the texture to draw and where to draw it

* ImagePacker can take png files and pack them into one texture.
* Adding a png to a packer returns a *Sprite that will represet the png in the texture
* call Pack() on your packer variable when you are done adding png's
