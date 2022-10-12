package av

import (
	"fmt"
	"io"
	"os/exec"
	"bytes"
	"image"
	_ "image/png"
)
// TODO - Note this is heavily command line and linux specific. Ideally we should move to wrappers and use the libraries that ffmpeg usees itself (or some other custom impl)

type FFmpeg struct {
	Cmd *exec.Cmd
	Stdin io.WriteCloser
	Stdout, Stderr io.ReadCloser
	Buffer []byte
}

// This is for FFmpeg output
func NewFFmpeg(width, height, fps int, filename string) (*FFmpeg, error) {
	cmdWxH := fmt.Sprintf("%dx%d", width, height)
	fpsStr := fmt.Sprintf("%d", fps)
	// You can adjust the quality by altering the -crf option in the command line (lower numbers are better quality). Don't forget to change the resolution (-s) and framerate (-r) to match what you want, as well.
	cmd := exec.Command("ffmpeg", "-r", fpsStr, "-f", "rawvideo", "-pix_fmt", "rgba", "-s", cmdWxH, "-i", "-", "-threads", "0", "-preset", "fast", "-y", "-pix_fmt", "yuv420p", "-crf", "21", "-vf", "vflip", filename)
	fmt.Println(fpsStr, cmd)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, err

	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdin.Close()
		stdout.Close()
		return nil, err
	}

	ffmpeg := FFmpeg{
		Cmd: cmd,
		Stdin: stdin,
		Stdout: stdout,
		Stderr: stderr,
		Buffer: make([]byte, 4 * width * height, 4 * width * height), // RGBA * window size
	}
	return &ffmpeg, nil
}

func (f *FFmpeg) Start() error {
	err := f.Cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

// Sends the buffer to the underlying FFMPEG process
func (f *FFmpeg) SendBuffer() error {
	_, err := f.Stdin.Write(f.Buffer)
	if err != nil {
		return err
	}
	return nil
}

func (f *FFmpeg) Close() {
	f.Stdin.Close()
	f.Stdout.Close()
	f.Stderr.Close()

	// Note: the file pipes must be closed before cmd.Wait() is called
	f.Cmd.Wait()
}

//--------------------------------------------------------------------------------

func GetFFmpegFrame(filename string, frameNumber int) (image.Image, error) {
	selectFilter := fmt.Sprintf("select=eq(n\\,%d)", frameNumber)

	cmd := exec.Command("ffmpeg", "-i", filename, "-vf", selectFilter, "-f", "image2pipe", "-c:v", "png", "-vframes", "1", "-")

	fmt.Println(cmd)

	// Note: Use this if you want to see what is going wrong
	// b, err := cmd.CombinedOutput()
	b, err := cmd.Output()
	if err != nil {
		fmt.Println(string(b))
		return nil, err
	}

	// TODO - I could maybe use the stdout reader directly for optimization?
	imgReader := bytes.NewReader(b)

	img, _, err := image.Decode(imgReader)
	if err != nil {
		return nil, err
	}
	return img, nil
}
