package tilemap

// Inspired By: http://www.cr31.co.uk/stagecast/wang/1sideedge.html


// Follows the tile numbering that godot uses. Numbered in reading order, starting at 0.
// Godot: https://docs.godotengine.org/en/stable/tutorials/2d/using_tilemaps.html
func PackedBlobmapNumber(t, b, l, r, tl, tr, bl, br bool) uint8 {
	wang := WangBlobmapNumber(t, b, l, r, tl, tr, bl, br)

	// Default blank
	ret := 22

	switch wang {
	case 0:
		ret = 22

	case 1:
		ret = 24
	case 4:
		ret = 37
	case 16:
		ret = 0
	case 64:
		ret = 39

	case 5:
		ret = 25
	case 20:
		ret = 1
	case 80:
		ret = 3
	case 65:
		ret = 27

	case 7:
		ret = 44
	case 28:
		ret = 8
	case 112:
		ret = 11
	case 193:
		ret = 47

	case 17:
		ret = 12
	case 68:
		ret = 38

	case 21:
		ret = 13
	case 84:
		ret = 2
	case 81:
		ret = 15
	case 69:
		ret = 26

	case 23:
		ret = 28
	case 92:
		ret = 5
	case 113:
		ret = 19
	case 197:
		ret = 42

	case 29:
		ret = 16
	case 116:
		ret = 6
	case 209:
		ret = 31
	case 71:
		ret = 41

	case 31:
		ret = 20
	case 124:
		ret = 10
	case 241:
		ret = 35
	case 199:
		ret = 45

	case 85:
		ret = 14

	case 87:
		ret = 7
	case 93:
		ret = 43
	case 117:
		ret = 40
	case 219:
		ret = 4

	case 95:
		ret = 32
	case 125:
		ret = 9
	case 245:
		ret = 23
	case 215:
		ret = 46

	case 119:
		ret = 21
	case 221:
		ret = 34

	case 127:
		ret = 17
	case 253:
		ret = 18
	case 247:
		ret = 30
	case 223:
		ret = 29

	case 255:
		ret = 33
	}

	return uint8(ret)
}


// This function computes the wang tilenumber of a tile based on the tiles around it
func WangBlobmapNumber(t, b, l, r, tl, tr, bl, br bool) uint8 {
	// If surrounding edges aren't set, then corners must be false
	if !(t && l) { tl = false }
	if !(t && r) { tr = false }
	if !(b && l) { bl = false }
	if !(b && r) { br = false }

	total := uint8(0)
	if t	{ total	+= (1 << 0) }
	if tr { total += (1 << 1) }
	if r	{ total	+= (1 << 2) }
	if br { total += (1 << 3) }
	if b	{ total	+= (1 << 4) }
	if bl { total += (1 << 5) }
	if l	{ total	+= (1 << 6) }
	if tl { total += (1 << 7) }

	return total
}
