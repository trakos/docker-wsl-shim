package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

func runDockerCommandInWsl(dockerCommand string) {
	args := os.Args[1:]
	for i, v := range args {
		args[i] = convertArgument(v)
	}
	args = append([]string{"-d", getDefaultDistro(), dockerCommand}, args...)

	cmd := exec.Command("wsl", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		os.Exit(0)
	}
	var exitErr, ok = err.(*exec.ExitError)
	if ok {
		os.Exit(exitErr.ExitCode())
	}
	fmt.Printf("shim command failed %v", err)
}

var defaultDistro string = ""

func getDefaultDistro() string {
	if defaultDistro != "" {
		return defaultDistro
	}
	reg := regexp.MustCompile(`(.*) \(Default\)`)

	out, err2 := exec.Command("wsl", "-l").Output()
	if err2 != nil {
		println("wsl -l failed")
		os.Exit(5)
	}
	output, err := decodeUTF16(out)
	if err != nil {
		println("wsl -l not utf16")
		os.Exit(5)
	}

	if !reg.MatchString(output) {
		println("wsl -l has not returned default wsl instance")
		os.Exit(5)
	}

	rs := reg.FindStringSubmatch(output)

	defaultDistro = rs[1]

	return defaultDistro
}

func decodeUTF16(b []byte) (string, error) {
	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}
	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)
	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}

func convertArgument(argument string) string {
	reg := regexp.MustCompile(`(?i)^\\\\wsl\$\\` + getDefaultDistro())

	if reg.MatchString(argument) {
		argument = reg.ReplaceAllString(argument, "")
		argument = strings.ReplaceAll(argument, "\\", "/")
	}

	reg = regexp.MustCompile(`^([a-zA-Z]):[/\\]`)

	if reg.MatchString(argument) {
		argument = strings.ReplaceAll(argument, "\\", "/")
		out, err2 := exec.Command("wsl", "wslpath", argument).Output()
		if err2 != nil {
			println("wslpath failed")
			os.Exit(6)
		}
		argument = string(out)
		argument = strings.TrimSuffix(argument, "\n")
		argument = strings.TrimSuffix(argument, "\r")
	}

	return argument
}