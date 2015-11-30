package main

import (
	"bufio"
	"bytes"
	"fmt"
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
	fmt.Println(ds)
	return nil
}

type Display struct {
	Name  string
	Sizes []string
}

func GetDisplays() ([]Display, error) {
	b, err := ExecWithMsg("xrandr")
	if err != nil {
		return nil, err
	}

	ds := make([]Display, 0)

	sc := bufio.NewScanner(bytes.NewBuffer(b))
	for sc.Scan() {
		name, err := submatch(sc.Text(), `^(\w+) connected`, 0)
		if _, ok := err.(*notMatchErr); ok {
			size, err := submatch(sc.Text(), `^\s+(\d+x\d+)\s+`, 0)
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

func ExecWithMsg(cmd string, args ...string) ([]byte, error) {
	fmt.Printf("%s %s\n", cmd, strings.Join(args, " "))
	return exec.Command(cmd, args...).Output()
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
