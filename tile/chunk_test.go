package tile

import (
	"fmt"
	"testing"
)

func TestHashUnhash(t *testing.T) {
	a := int16(-1)
	fmt.Println(int16(a))
	fmt.Println(uint32(int16(a))) // 4294967295
	fmt.Println(uint32(uint16(int16(a)))) // 4294967295
	fmt.Println(int16(uint32(int16(a))))

	positions := []ChunkPosition{
		ChunkPosition{0, 0},
		ChunkPosition{1, 1},
		ChunkPosition{1, -1},
		ChunkPosition{-1, 1},
		ChunkPosition{-1, -1},
	}
	for _, c := range positions {
		hashed := c.hash()
		fmt.Println("----:", c)
		fmt.Println("Hashed: ", hashed)
		fmt.Println("Compare: ", c, fromHash(hashed))
	}
}
