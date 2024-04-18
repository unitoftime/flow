package browser

type OpenType uint8
const (
	OpenSameTab OpenType = iota
	OpenNewTab
	// OpenNewWindow
	OpenNewWindowBorderless
)
