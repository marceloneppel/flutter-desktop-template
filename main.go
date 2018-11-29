package main

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/Drakirus/go-flutter-desktop-embedder"

	"path/filepath"

	"encoding/json"

	"runtime"

	"errors"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const configurationFilename = "config.json"

type configuration struct {
	FlutterPath        string
	FlutterProjectPath string
	IconPath           string
	ScreenHeight       int
	ScreenWidth        int
}

func buildAssetPath(flutterProjectPath string, assetPath string) (string, error) {
	if flutterProjectPath == "" {
		var (
			path string
			err error
		)
		path, err = os.Executable()
		if err != nil {
			return "", err
		}
		return filepath.Join(filepath.Dir(path), "flutter_assets"), nil
	}
	return filepath.Join(flutterProjectPath, assetPath), nil
}

func buildICUDataPath(flutterPath string, icuDataPath string) (string, error) {
	if flutterPath == "" {
		var (
			path string
			err error
		)
		path, err = os.Executable()
		if err != nil {
			return "", err
		}
		return filepath.Join(filepath.Dir(path), "icudtl.dat"), nil
	}
	return filepath.Join(flutterPath, icuDataPath), nil
}

func getConfig() (configuration, error) {
	var config configuration
	var err error
	var file *os.File
	var path string
	path, err = os.Executable()
	if err != nil {
		return config, err
	}
	var configFilename = filepath.Join(filepath.Dir(path), configurationFilename)
	file, err = os.Open(configFilename)
	if err != nil {
		return config, err
	}
	var decoder = json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func getPaths() (string, string, error) {
	var assetPath string
	var icuDataPath string
	var err error
	switch runtime.GOOS {
	case "darwin":
		assetPath = "build/flutter_assets"
		icuDataPath = "bin/cache/artifacts/engine/darwin-x64/icudtl.dat"
		break
	case "linux":
		assetPath = "build/flutter_assets"
		icuDataPath = "bin/cache/artifacts/engine/linux-x64/icudtl.dat"
		break
	case "windows":
		assetPath = "build\\flutter_assets"
		icuDataPath = "bin\\cache\\artifacts\\engine\\windows-x64\\icudtl.dat"
		break
	default:
		err = errors.New("invalid operating system")
		break
	}
	return assetPath, icuDataPath, err
}

func handleError(err error) {
	log.Fatalln(err)
}

func main() {
	runtime.LockOSThread()
	var (
		config configuration
		err    error
	)
	config, err = getConfig()
	if err != nil {
		handleError(err)
	} else {
		var setIcon = func(window *glfw.Window) error {
			var (
				imgFile *os.File
				err     error
			)
			var iconPath = config.IconPath
			if string(iconPath[0]) != "/" {
				var path string
				path, err = os.Executable()
				if err != nil {
					return err
				}
				iconPath = filepath.Join(filepath.Dir(path), iconPath)
			}
			imgFile, err = os.Open(iconPath)
			if err != nil {
				return err
			}
			var img image.Image
			img, _, err = image.Decode(imgFile)
			if err != nil {
				return err
			}
			window.SetIcon([]image.Image{img})
			return nil
		}
		var (
			assetPath   string
			icuDataPath string
		)
		assetPath, icuDataPath, err = getPaths()
		if err != nil {
			handleError(err)
		} else {
			var builtAssetPath string
			builtAssetPath, err = buildAssetPath(config.FlutterProjectPath, assetPath)
			if err != nil {
				handleError(err)
			} else {
				var builtICUDataPath string
				builtICUDataPath, err = buildICUDataPath(config.FlutterPath, icuDataPath)
				if err != nil {
					handleError(err)
				} else {
					var options = []gutter.Option{
						gutter.OptionAssetPath(builtAssetPath),
						gutter.OptionICUDataPath(builtICUDataPath),
						gutter.OptionWindowInitializer(setIcon),
						gutter.OptionWindowDimension(config.ScreenWidth, config.ScreenHeight),
						gutter.OptionWindowInitializer(setIcon),
						gutter.OptionPixelRatio(1.9),
						gutter.OptionVmArguments([]string{"--dart-non-checked-mode", "--observatory-port=50300"}),
					}
					err = gutter.Run(options...)
					if err != nil {
						handleError(err)
					}
				}
			}
		}
	}
}
