package config

import (
	"fmt"
	"os"
	"path/filepath"
	"to-do-list/pkg/env"

	"gopkg.in/yaml.v3"
)

func New(repoName string) (*Config, error) {
	filename := getConfigFile(repoName, env.ServiceEnv())
	return newWithFile(filename)
}

func newWithFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil

}

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
}

type Server struct {
	HTTP HTTPServer `yaml:"http"`
}

type HTTPServer struct {
	Address string `yaml:"address"`
}

type Database struct {
	Driver     string `yaml:"driver"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	DbName     string `yaml:"dbname"`
	Credential string `yaml:"credential"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
}

func getConfigFile(repoName, env string) string {
	var (
		filename = fmt.Sprintf("%s.%s.yaml", repoName, env)
	)

	// use local files in dev
	// repoPath := filepath.Join(repoName)
	return filepath.Join("files/etc", repoName, filename)
}
