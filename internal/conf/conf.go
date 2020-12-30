package conf

import (
    "io/ioutil"
    "os"

    "github.com/st3iny/batteryhook/internal/util"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Hooks []Hook `yaml:"hooks"`
}

func Load() (*Config, error) {
    configPath, err := util.BuildConfigPath("config.yaml")
    if err != nil {
        return nil, err
    }

    _, err = os.Stat(configPath)
    if err != nil {
        return nil, err
    }

    configBlob, err := ioutil.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var config Config
    err = yaml.Unmarshal(configBlob, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}
