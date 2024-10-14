package glitchasset

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"io/ioutil"
	"path"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/unitoftime/flow/glm"
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

func (load *Load) Data(filepath string) ([]byte, error) {
	file, err := load.filesystem.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func (load *Load) Font(filepath string, size float64) (font.Face, error) {
	file, err := load.filesystem.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fontData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}
	fontFace := truetype.NewFace(font, &truetype.Options{
		Size: size,
	})
	return fontFace, nil
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

func (load *Load) Texture(filepath string, smooth bool) (*glitch.Texture, error) {
	img, err := load.Image(filepath)
	if err != nil {
		return nil, err
	}

	texture := glitch.NewTexture(img, smooth)
	return texture, nil
}

// Loads a single sprite from a filepath of an image
func (load *Load) Sprite(filepath string, smooth bool) (*glitch.Sprite, error) {
	texture, err := load.Texture(filepath, smooth)
	if err != nil {
		return nil, err
	}
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

func (load *Load) Yaml(filepath string, dat interface{}) error {
	file, err := load.filesystem.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	yamlData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlData, dat)
}

// TODO - move Aseprite stuff to another package?
type AseSheet struct {
	Frames []AseFrame `json:frames`
	Meta   AseMeta
}
type AseFrame struct {
	Filename string `json:filename`
	//Frame todo
	Duration int `json:duration`
}
type AseMeta struct {
	FrameTags []AseFrameTag `json:frameTags`
}
type AseFrameTag struct {
	Name      string `json:name`
	From      int    `json:from`
	To        int    `json:to`
	Direction string `json:direction`
}

// Loads an aseprite spritesheet.
func (load *Load) AseSheet(filepath string) (*AseSheet, error) {
	dat := AseSheet{}
	err := load.Json(filepath, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, nil
}

// TODO - Assumes that all animations share the same spritesheet
type Animation struct {
	spritesheet *Spritesheet
	Frames      map[string][]AnimationFrame
}
type AnimationFrame struct {
	Name     string
	Sprite   *glitch.Sprite
	Duration time.Duration
	MirrorY  bool
}

// TODO - Assumptions: frame name is <filename>-<framenumber>.png (Aseprite doesn't export the file name. But you could maybe repack their spritesheet into your own)
func (load *Load) AseAnimation(spritesheet *Spritesheet, filepath string) (*Animation, error) {
	base := path.Base(filepath)
	baseNoExt := base[:len(base)-len(path.Ext(base))]

	aseSheet, err := load.AseSheet(filepath)
	if err != nil {
		return nil, err
	}

	anim := Animation{
		spritesheet: spritesheet,
		Frames:      make(map[string][]AnimationFrame),
	}

	for _, frameTag := range aseSheet.Meta.FrameTags {
		// TODO - implement other directions
		if frameTag.Direction != "forward" {
			panic("NonForward frametag not supported!")
		}

		frames := make([]AnimationFrame, 0)
		for i := frameTag.From; i <= frameTag.To; i++ {
			spriteName := fmt.Sprintf("%s-%d.png", baseNoExt, i)
			sprite, err := spritesheet.Get(spriteName)
			if err != nil {
				return nil, err
			}
			frames = append(frames, AnimationFrame{
				Name:     spriteName,
				Sprite:   sprite,
				Duration: time.Duration(aseSheet.Frames[i].Duration) * time.Millisecond,
				MirrorY:  false,
			})
		}
		anim.Frames[frameTag.Name] = frames
	}

	return &anim, nil
}

func (load *Load) Mountpoints(filepath string) (packer.MountFrames, error) {
	mountFrames := packer.MountFrames{}
	err := load.Json(filepath, &mountFrames)
	if err != nil {
		return packer.MountFrames{}, err
	}

	return mountFrames, nil
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
		rect := glm.R(
			v.Frame.X,
			v.Frame.Y,
			v.Frame.X+v.Frame.W,
			v.Frame.Y+v.Frame.H).Norm()

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
	lookup  map[string]*glitch.Sprite
}

func NewSpritesheet(tex *glitch.Texture, lookup map[string]*glitch.Sprite) *Spritesheet {
	return &Spritesheet{
		texture: tex,
		lookup:  lookup,
	}
}

func (s *Spritesheet) Get(name string) (*glitch.Sprite, error) {
	sprite, ok := s.lookup[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Invalid sprite name: %s", name))
	}
	return sprite, nil
}

// https://www.aseprite.org/docs/slices/#:~:text=With%20the%20Slice%20tool,some%20extra%20user%20defined%20information.
func (s *Spritesheet) GetNinePanel(name string, border glm.Rect) (*glitch.NinePanelSprite, error) {
	sprite, ok := s.lookup[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Invalid sprite name: %s", name))
	}
	return glitch.SpriteToNinePanel(sprite, border), nil
}

func (s *Spritesheet) Picture() *glitch.Texture {
	return s.texture
}

// // Gets multiple frames with the same prefix name. Indexing starts at 0
// func (s *Spritesheet) GetFrames(name, ext string, length int) ([]*glitch.Sprite, error) {
// 	ret := make([]*glitch.Sprite, length)
// 	for i := range names {
// 		sprite, err := s.Get(fmt.Sprintf("%s%d%s", name, i, ext))
// 		if err != nil {
// 			return nil, err
// 		}
// 		ret[i] = sprite
// 	}
// 	return ret, nil
// }
