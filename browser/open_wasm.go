//go:build js || wasm

package browser

import (
	"syscall/js"
)

var window = js.Global().Get("window")
var location = window.Get("location")

func Open(url string, newTab bool) error {
	if newTab {
		// "location=yes,height=570,width=520,scrollbars=yes,status=yes")
		// window.Call("open", url, "_blank", "location=yes,scrollbars=yes,status=yes")
		window.Call("open", url, "_blank") // Open as new tab
	} else {
		window.Call("open", url, "_blank", "location=yes,scrollbars=yes,status=yes") // Open as window
		// try to use just window.location?
		// location.Set("href", url)
	}

	return nil
}
