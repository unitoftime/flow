package main

import (
	"image"
	_ "image/png"
	"bytes"
	"errors"

	"github.com/unitoftime/glitch"
	"github.com/unitoftime/flow/asset"
)

type SpriteAssetLoader struct {
}
func (l SpriteAssetLoader) Ext() []string {
	return []string{".png"}
}
func (l SpriteAssetLoader) Load(server *asset.Server, data []byte) (*glitch.Sprite, error) {
	smooth := true

	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	texture := glitch.NewTexture(img, smooth)

	return glitch.NewSprite(texture, texture.Bounds()), nil
}
func (l SpriteAssetLoader) Store(server *asset.Server, sprite *glitch.Sprite) ([]byte, error) {
	return nil, errors.New("sprites do not support writeback")
}
