package asset

import (
	"fmt"
	"errors"
	"io/fs"
	"image"
	_ "image/png"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/unitoftime/glitch"
	"github.com/unitoftime/packer"
)

type Load struct {
	filesystem fs.FS
}

func NewLoad(filesystem fs.FS) *Load {
	return &Load{filesystem}
}

func (load *Load) Open(filepath string) (fs.File, error) {
	return load.filesystem.Open(filepath)
}

func (load *Load) Image(filepath string) (image.Image, error) {
	file, err := load.filesystem.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Loads a single sprite from a filepath of an image
func (load *Load) Sprite(filepath string, smooth bool) (*glitch.Sprite, error) {
	img, err := load.Image(filepath)
	if err != nil {
		return nil, err
	}
	// pic := pixel.PictureDataFromImage(img)
	// return pixel.NewSprite(pic, pic.Bounds()), nil
	texture := glitch.NewTexture(img, smooth)
	return glitch.NewSprite(texture, texture.Bounds()), nil
}

func (load *Load) Json(filepath string, dat interface{}) error {
	file, err := load.filesystem.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, dat)
}

func (load *Load) Spritesheet(filepath string, smooth bool) (*Spritesheet, error) {
	//Load the Json
	serializedSpritesheet := packer.SerializedSpritesheet{}
	err := load.Json(filepath, &serializedSpritesheet)
	if err != nil {
		return nil, err
	}

	imageFilepath := path.Join(path.Dir(filepath), serializedSpritesheet.ImageName)

	// Load the image
	img, err := load.Image(imageFilepath)
	if err != nil {
		return nil, err
	}
	// pic := pixel.PictureDataFromImage(img)
	texture := glitch.NewTexture(img, smooth)

	// Create the spritesheet object
	// bounds := texture.Bounds()
	lookup := make(map[string]*glitch.Sprite)
	for k, v := range serializedSpritesheet.Frames {
		rect := glitch.R(
			float32(v.Frame.X),
			float32(v.Frame.Y),
			float32(v.Frame.X + v.Frame.W),
			float32(v.Frame.Y + v.Frame.H)).Norm()

		// rect := glitch.R(
		// 	float32(v.Frame.X),
		// 	float32(float64(bounds.H()) - v.Frame.Y),
		// 	float32(v.Frame.X + v.Frame.W),
		// 	float32(float64(bounds.W()) - (v.Frame.Y + v.Frame.H))).Norm()

		lookup[k] = glitch.NewSprite(texture, rect)
	}

	return NewSpritesheet(texture, lookup), nil
}

type Spritesheet struct {
	texture *glitch.Texture
	lookup map[string]*glitch.Sprite
}

func NewSpritesheet(tex *glitch.Texture, lookup map[string]*glitch.Sprite) *Spritesheet {
	return &Spritesheet{
		texture: tex,
		lookup: lookup,
	}
}

func (s *Spritesheet) Get(name string) (*glitch.Sprite, error) {
	sprite, ok := s.lookup[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Invalid sprite name: %s", name))
	}
	return sprite, nil
}

func (s *Spritesheet) GetNinePanel(name string, border glitch.Rect) (*glitch.NinePanelSprite, error) {
	sprite, ok := s.lookup[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Invalid sprite name: %s", name))
	}
	return glitch.SpriteToNinePanel(sprite, border), nil
}

func (s *Spritesheet) Picture() *glitch.Texture {
	return s.texture
}
