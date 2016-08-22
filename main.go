package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
}

func Main() error {
	ds, err := GetDisplays()
	if err != nil {
		return err
	}

	sizeLevel := 0
	flag.IntVar(&sizeLevel, "size", 0, "size level")
	flag.Parse()

	size, err := getSize(ds, sizeLevel)
	if err != nil {
		return err
	}

	for i, d := range ds {
		args := []string{"--output", d.Name, "--mode", size}
		if i != 0 {
			args = append(args, "--same-as", ds[0].Name)
		}
		err := ExecWithMsg("xrandr", args...)
		if err != nil {
			return err
		}
	}
	return nil
}

// returns size
func getSize(displays []Display, sizeLevel int) (string, error) {
	s := make([][]string, 0, len(displays))
	for _, d := range displays {
		s = append(s, d.Sizes)
	}
	sizes := Intersection(s...)
	if len(sizes) < sizeLevel {
		return "", fmt.Errorf("Size level too deep. Please specify 0..%d", len(sizes))
	}
	return sizes[sizeLevel], nil
}

type Display struct {
	Name  string
	Sizes []string
}

func GetDisplays() ([]Display, error) {
	b, err := exec.Command("xrandr").Output()
	if err != nil {
		return nil, err
	}

	ds := make([]Display, 0)

	sc := bufio.NewScanner(bytes.NewBuffer(b))
	for sc.Scan() {
		name, err := submatch(sc.Text(), `^(\w+) connected`, 0)
		if _, ok := err.(*notMatchErr); ok {
			size, err := submatch(sc.Text(), `^\s+(\d+x\w+)\s+`, 0)
			if _, ok := err.(*notMatchErr); ok {
				continue
			}
			if err != nil {
				return nil, err
			}
			ds[len(ds)-1].Sizes = append(ds[len(ds)-1].Sizes, size)
			continue
		}

		if err != nil {
			return nil, err
		}

		d := Display{
			Name:  name,
			Sizes: make([]string, 0),
		}
		ds = append(ds, d)
	}

	return ds, nil
}

func ExecWithMsg(cmd string, args ...string) error {
	fmt.Printf("%s %s\n", cmd, strings.Join(args, " "))
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

type notMatchErr struct {
	re string
	s  string
}

func (e *notMatchErr) Error() string {
	return fmt.Sprintf("%s doesnot match %q", e.re, e.s)
}

var _ error = &notMatchErr{}

func submatch(s, re string, n int) (string, error) {
	n++
	reg, err := regexp.Compile(re)
	if err != nil {
		return "", err
	}
	ma := reg.FindStringSubmatch(s)
	if len(ma) <= n {
		return "", &notMatchErr{re: re, s: s}
	}
	return ma[n], nil
}

func Intersection(s ...[]string) []string {
	if len(s) == 1 {
		return s[0]
	}

	res := make([]string, 0)
	for _, v := range s[0] {
		for _, w := range s[1] {
			if v == w {
				res = append(res, v)
				break
			}
		}
	}
	return Intersection(append(s[2:], res)...)
}
