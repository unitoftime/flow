//go:build !js
package browser

import (
	"github.com/pkg/browser"
)

func Open(url string, _ bool) error {
	return browser.OpenURL(url)
}
