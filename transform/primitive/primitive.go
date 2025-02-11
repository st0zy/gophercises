package primitive

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Mode int8

const (
	Combo Mode = iota
	Triangle
	Rect
	Ellipse
	Circle
	RotatedRect
	Beziers
	Rotatedellipse
	Psolygon
)

func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	in, err := os.CreateTemp("/tmp", fmt.Sprintf("in_*%s", ext))
	if err != nil {
		return nil, err
	}
	defer os.Remove(in.Name())

	out, err := os.CreateTemp("/tmp", fmt.Sprintf("out_*%s", ext))
	if err != nil {
		return nil, err
	}
	defer os.Remove(out.Name())

	_, err = io.Copy(in, image)
	if err != nil {
		return nil, err
	}

	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, opts...)
	if err != nil {
		return nil, fmt.Errorf("primitive: failed to run the primitive command. err=%s", err)
	}
	fmt.Println(stdCombo)

	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func primitive(inputFile, outputFile string, numShapes int, opts ...func() []string) (string, error) {
	argsString := fmt.Sprintf("-i %s -o %s -n %d", inputFile, outputFile, numShapes)
	args := strings.Fields(argsString)
	for _, opt := range opts {
		args = append(args, opt()...)
	}
	fmt.Println(args)
	cmd := exec.Command("primitive", args...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}
