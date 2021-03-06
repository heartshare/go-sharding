package config

import (
	"github.com/XiaoMi/Gaea/core"
	"github.com/XiaoMi/Gaea/logging"
	"github.com/XiaoMi/Gaea/util"
	"go.uber.org/config"
	"os"
	"path/filepath"
)

var logger = logging.GetLogger("config")

type BootConfig struct {
	Provider string
}

func LoadConfig() *BootConfig {
	var sources []config.YAMLOption

	files := defaultFileLocations()

	var sb = core.NewStringBuilder()
	sb.WriteLine()
	sb.WriteLine("Search configuration locations:")
	for _, f := range files {
		if util.FileExists(f) {
			sources = append(sources, config.File(f))
			sb.WriteLine("[Found]:", f)
		} else {
			sb.WriteLine("[Not Found]:", f)
		}
	}

	bootCnf := &BootConfig{
		Provider: "file",
	}

	if len(sources) > 0 {
		var err error
		var yaml *config.YAML
		if yaml, err = config.NewYAML(sources...); err == nil {
			err = yaml.Get("config").Populate(bootCnf)
		}
		if err != nil {
			logger.Warn("Load boot config file fault.", util.LineSeparator, err)
		}
	}

	return bootCnf
}

func defaultFileLocations() []string {
	files := make(map[string]bool, 3)
	if !util.IsWindows() {
		files["/etc/go-sharding/config.yaml"] = false
		files["/etc/go-sharding/config.yml"] = false
	}
	dir, err := os.Getwd()
	if err == nil {
		files[filepath.Join(dir, "config.yaml")] = false
	} else {
		files["config.yaml"] = false
	}

	result := make([]string, len(files))
	i := 0
	for k, _ := range files {
		result[i] = k
		i++
	}
	return result
}
