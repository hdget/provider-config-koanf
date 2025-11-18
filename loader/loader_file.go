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
		configFile, err = l.guessConfigFile()
		if err != nil {
			return err
		}
	}

	return l.reader.Load(file.Provider(configFile), toml.Parser())
}

func (l *fileConfigLoader) guessConfigFile() (string, error) {
	filename := l.getDefaultConfigFileName()
	dir := l.findConfigDir()
	if dir == "" {
		return "", fmt.Errorf("config dir not found, app: %s, env: %s", l.app, l.env)
	}

	return filepath.Join(dir, filename), nil
}

// getDefaultConfigFilename 缺省的配置文件名: <app>.<env>.toml
func (l *fileConfigLoader) getDefaultConfigFileName() string {
	return strings.Join([]string{l.app, l.env, defaultConfigType}, ".")
}

// findConfigDirs 缺省的配置文件名: <app>.<env>
func (l *fileConfigLoader) findConfigDir() string {
	// iter to root directory
	absStartPath, err := filepath.Abs(".")
	if err != nil {
		return ""
	}

	var found string
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
				found = filepath.Join(currPath, dirName)
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

	return found
}
