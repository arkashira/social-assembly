package moderation

import (
	"os"
	"path/filepath"

	"github.com/axentx/axentx-social-assembly/pkg/config"
)

func GetConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "axentx-social-assembly")
}