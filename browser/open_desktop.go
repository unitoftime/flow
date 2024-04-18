//go:build !js

package browser

import (
	"github.com/pkg/browser"
)

func Open(url string, _ OpenType) error {
	return browser.OpenURL(url)
}
