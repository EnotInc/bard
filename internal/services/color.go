package services

import (
	"fmt"
	"strings"
)

type Apply bool

const (
	Foreground Apply = true
	Background Apply = false
)

func HexToAscii(c string, a Apply) (string, error) {
	hex := strings.TrimPrefix(c, "#")
	if len(hex) != 6 {
		return "", fmt.Errorf("Invalid hex len in string: %s", c)
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%2x%2x%2x", &r, &g, &b)
	if err != nil {
		return "", err
	}

	var escape string
	switch a {
	case Foreground:
		escape = fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	case Background:
		escape = fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
	}

	return escape, nil
}
