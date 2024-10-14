package main

import (
	"bytes"
	"errors"
	"image"
	_ "image/png"

	"github.com/unitoftime/flow/asset"
	"github.com/unitoftime/glitch"
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
