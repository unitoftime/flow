package av

import (
	"fmt"
	"os"
	"os/exec"
	"hash/crc32"
)

func Codeblock(inner string) string {
	return "```"+inner+"```"
}

// Calls marp with markdown to generate an image. Everything after --- will be ignored
// Returns the image file which is uniquely named based on the passed in contents
func Markdown(contents string) string {
	name := fmt.Sprintf("%d.png",
		crc32.ChecksumIEEE([]byte(contents)))

	_, err := os.Stat(name)
	if err == nil {
		// File already exists
		return name
	}

	header := `
<!-- theme: uncover -->
<!-- class: invert -->
`
	os.WriteFile("marp.tmp", []byte(header + contents), 0755)

	cmd := exec.Command("marp", "marp.tmp", "--image", "--image-scale", "2", "-o", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))

	return name
}
