package main

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/Drakirus/go-flutter-desktop-embedder"

	"path/filepath"

	"encoding/json"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const assetPath = "build/flutter_assets"
const configurationFilename = "config.json"
const icuDataPath = "bin/cache/artifacts/engine/linux-x64/icudtl.dat"

type Configuration struct {
	FlutterPath        string
	FlutterProjectPath string
	IconPath           string
	ScreenHeight       int
	ScreenWidth        int
}

func getConfig() (Configuration, error) {
	var configuration Configuration
	var err error
	var file *os.File
	file, err = os.Open(configurationFilename)
	if err != nil {
		return configuration, err
	} else {
		var decoder = json.NewDecoder(file)
		err = decoder.Decode(&configuration)
		if err != nil {
			return configuration, err
		} else {
			return configuration, nil
		}
	}
}

func handleError(err error) {
	log.Fatalln(err)
}

func main() {
	var (
		configuration Configuration
		err           error
	)
	configuration, err = getConfig()
	if err != nil {
		handleError(err)
	} else {
		var setIcon = func(window *glfw.Window) error {
			var (
				imgFile *os.File
				err     error
			)
			imgFile, err = os.Open(configuration.IconPath)
			if err != nil {
				return err
			} else {
				var img image.Image
				img, _, err = image.Decode(imgFile)
				if err != nil {
					return err
				} else {
					window.SetIcon([]image.Image{img})
					return nil
				}
			}
		}
		var options = []gutter.Option{
			gutter.OptionAssetPath(filepath.Join(configuration.FlutterProjectPath, assetPath)),
			gutter.OptionICUDataPath(filepath.Join(configuration.FlutterPath, icuDataPath)),
			gutter.OptionWindowInitializer(setIcon),
			gutter.OptionWindowDimension(configuration.ScreenWidth, configuration.ScreenHeight),
			gutter.OptionWindowInitializer(setIcon),
			gutter.OptionPixelRatio(1.2),
			gutter.OptionVMArguments([]string{"--dart-non-checked-mode", "--observatory-port=50300"}),
		}
		err = gutter.Run(options...)
		if err != nil {
			handleError(err)
		}
	}
}
