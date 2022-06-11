package av

import (
	"fmt"
	"io"
	"os/exec"
)

type FFmpeg struct {
	Cmd *exec.Cmd
	Stdin io.WriteCloser
	Stdout, Stderr io.ReadCloser
	Buffer []byte
}

func NewFFmpeg(width, height int, filename string) (*FFmpeg, error) {
	cmdWxH := fmt.Sprintf("%dx%d", width, height)
	cmd := exec.Command("ffmpeg", "-r", "60", "-f", "rawvideo", "-pix_fmt", "rgba", "-s", cmdWxH, "-i", "-", "-threads", "0", "-preset", "fast", "-y", "-pix_fmt", "yuv420p", "-crf", "21", "-vf", "vflip", filename)

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
