package loader

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type fileConfigLoader struct {
	reader     *koanf.Koanf
	app        string
	env        string
	configFile string
}

const (
	defaultConfigType = "toml"
)

var (
	// the default config file search pattern
	//
	//	./config/app/<app>/<app>.test.toml
	//	./common/config/app/<app>/<app>.test.toml
	//	../config/app/<app>/<app>.test.toml
	//  ../common/config/app/<app>/<app>.test.toml
	//  ...
	defaultConfigRootDirs = []string{
		".",                                      // current dir
		filepath.Join("config", "app"),           // default config root dir1
		filepath.Join("common", "config", "app"), // default config root dir2
	}
)

func NewFileConfigLoader(reader *koanf.Koanf, app, env, configFile string) Loader {
	return &fileConfigLoader{
		reader:     reader,
		app:        app,
		env:        env,
		configFile: configFile,
	}
}

func (l *fileConfigLoader) Load() error {
	var err error
	configFile := l.configFile
	if configFile == "" {
		configFile, err = l.findConfigFile()
		if err != nil {
			return err
		}
	}

	return l.reader.Load(file.Provider(configFile), toml.Parser())
}

// getDefaultConfigFilename 缺省的配置文件名: <app>.<env>.toml
func (l *fileConfigLoader) getDefaultConfigFileName() string {
	if l.env != "" {
		return strings.Join([]string{l.app, l.env, defaultConfigType}, ".")
	}
	return strings.Join([]string{l.app, defaultConfigType}, ".")
}

// findConfigDirs 缺省的配置文件名: <app>.<env>
func (l *fileConfigLoader) findConfigFile() (string, error) {
	// iter to root directory
	absStartPath, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}

	var foundDir string
	matchFile := l.getDefaultConfigFileName()
	currPath := absStartPath
LOOP:
	for {
		for _, rootDir := range defaultConfigRootDirs {
			// possible parent dir name
			dirName := filepath.Join(rootDir, l.app)
			checkDir := filepath.Join(currPath, dirName, matchFile)
			matches, err := filepath.Glob(checkDir)
			if err == nil && len(matches) > 0 {
				foundDir = filepath.Join(currPath, dirName)
				break LOOP
			}
		}

		// If we're already at the root, stop finding
		// windows has the driver name, so it need use TrimRight to test
		abs, _ := filepath.Abs(currPath)
		if abs == string(filepath.Separator) || len(strings.TrimRight(currPath, string(filepath.Separator))) <= 3 {
			break
		}

		// else, get parent dir
		currPath = filepath.Dir(currPath)
	}

	if foundDir == "" {
		return "", fmt.Errorf(`config file "%s" not found`, matchFile)
	}

	return filepath.Join(foundDir, matchFile), nil
}
