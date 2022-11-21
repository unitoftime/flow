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
func Markdown(basedir, contents string) string {
	name := fmt.Sprintf("%s/%d.png",
		basedir,
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
	// TODO - pass marp.tmp via stdin
	marpFile := basedir + "/marp.tmp"
	os.WriteFile(marpFile, []byte(header + contents), 0755)

	// Note: 1.5 is 1080
	cmd := exec.Command("marp", marpFile, "--image", "--image-scale", "1.5", "-o", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))

	return name
}
