//go:build js || wasm

package browser

import (
	"syscall/js"
)

var window = js.Global().Get("window")
var location = window.Get("location")

func Open(url string, openType OpenType) error {
	switch openType {
	case OpenSameTab:
		location.Set("href", url)
	case OpenNewTab:
		window.Call("open", url, "_blank") // Open as new tab
	// case OpenNewWindow:
	// 	window.Call("open", url, "_blank") // Open as window // TODO: Does this work?
	case OpenNewWindowBorderless:
		window.Call("open", url, "_blank", "location=yes,scrollbars=yes,status=yes") // Open as window
	}
	// if newTab {
	// 	// "location=yes,height=570,width=520,scrollbars=yes,status=yes")
	// 	// window.Call("open", url, "_blank", "location=yes,scrollbars=yes,status=yes")
	// 	window.Call("open", url, "_blank") // Open as new tab
	// } else {
	// 	window.Call("open", url, "_blank", "location=yes,scrollbars=yes,status=yes") // Open as window
	// 	// try to use just window.location?
	// 	// location.Set("href", url)
	// }

	return nil
}
