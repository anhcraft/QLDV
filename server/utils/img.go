package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
)

var pngquantPath = ""
var imagemagickPath = ""

func init() {
	p, has := os.LookupEnv("pngquant_path")
	if !has {
		log.Fatal().Msg("pngquant_path is missing")
		os.Exit(1)
		return
	}
	pngquantPath = p
	p, has = os.LookupEnv("imagemagick_path")
	if !has {
		log.Fatal().Msg("imagemagick_path is missing")
		os.Exit(1)
		return
	}
	imagemagickPath = p
}

func ResizeAndCompressImage(path string, ext string, width int, height int) bool {
	if ext == ".png" {
		_, err := exec.Command(imagemagickPath, path, "-resize", fmt.Sprintf("%dx%d>", width, height), path+".temp").Output()
		if err != nil {
			log.Error().Err(err).Str("ext", ext).Str("path", path).Msg("An error occurred while resizing image")
			return false
		}
		_, err = exec.Command(pngquantPath, "--force", "--output", path, "--quality", "65-80", "256", path+".temp").Output()
		if err != nil {
			log.Error().Err(err).Str("ext", ext).Str("path", path).Msg("An error occurred while compressing image")
			return false
		}
		err = os.Remove(path + ".temp")
		if err != nil {
			log.Error().Err(err).Str("ext", ext).Str("path", path).Msg("An error occurred while deleting temporary file")
			return false
		}
	} else if ext == ".jpeg" {
		_, err := exec.Command(imagemagickPath, path, "-resize", fmt.Sprintf("%dx%d>", width, height), "-quality", "60", path).Output()
		if err != nil {
			log.Error().Err(err).Str("ext", ext).Str("path", path).Msg("An error occurred while resizing image")
			return false
		}
	}
	return true
}
