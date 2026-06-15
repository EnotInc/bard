package config

import (
	"os"
	"path/filepath"
)

const spacePath = ".bard/space/"

func GetSpacePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return spacePath
	}
	return filepath.Join(home, spacePath)
}

func CreateSpace() {
	path := GetSpacePath()
	_, err := os.Stat(path)
	if err != nil {
		os.Mkdir(path, 0755)
	}
}
