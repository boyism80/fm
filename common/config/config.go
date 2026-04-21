package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

func ResolvePath(p string) (string, error) {
	try := func(abs string) (string, bool) {
		st, err := os.Stat(abs)
		if err != nil || st.IsDir() {
			return "", false
		}
		out, err := filepath.Abs(abs)
		if err != nil {
			return abs, true
		}
		return out, true
	}

	if filepath.IsAbs(p) {
		if s, ok := try(p); ok {
			return s, nil
		}
		return "", fmt.Errorf("config file %q not found", p)
	}

	if wd, err := os.Getwd(); err == nil {
		if s, ok := try(filepath.Join(wd, p)); ok {
			return s, nil
		}
	}

	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		if s, ok := try(filepath.Join(exeDir, p)); ok {
			return s, nil
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("config file %q: getwd: %w", p, err)
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		if s, ok := try(filepath.Join(dir, p)); ok {
			return s, nil
		}
		next := filepath.Dir(dir)
		if next == dir {
			break
		}
	}

	return "", fmt.Errorf("config file %q not found (cwd=%q; try -config with a full path)", p, wd)
}

type InternalEndpoint struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	TimeoutSeconds int    `yaml:"timeout"`
}

func (e InternalEndpoint) GRPCAddr() string {
	if e.Host == "" || e.Port == 0 {
		return ""
	}
	return e.Host + ":" + strconv.Itoa(e.Port)
}

type RabbitMQEndpoint struct {
	IP    string `yaml:"ip"`
	Port  int    `yaml:"port"`
	UID   string `yaml:"uid"`
	PWD   string `yaml:"pwd"`
	VHost string `yaml:"vhost"`
}

func (e RabbitMQEndpoint) Enabled() bool {
	return e.IP != "" && e.Port > 0
}

func (e RabbitMQEndpoint) AMQPURL() string {
	if !e.Enabled() {
		return ""
	}
	vhost := e.VHost
	if vhost == "" {
		vhost = "/"
	}
	return "amqp://" + e.UID + ":" + e.PWD + "@" + e.IP + ":" + strconv.Itoa(e.Port) + "/" + vhost
}

type Login struct {
	Host                        string           `yaml:"host"`
	Port                        int              `yaml:"port"`
	InitialRole                 int              `yaml:"initial_role"`
	Internal                    InternalEndpoint `yaml:"internal"`
	CatalogRetryIntervalSeconds int              `yaml:"catalog_retry_interval_seconds"`
	CatalogRetryMaxAttempts     int              `yaml:"catalog_retry_max_attempts"`
}

type Game struct {
	Host       string           `yaml:"host"`
	Port       int              `yaml:"port"`
	ChannelId  int              `yaml:"channel_id"`
	WzPath     string           `yaml:"wz_path"`
	World      string           `yaml:"world"`
	WorldId    int              `yaml:"world_id"`
	MaxPlayers int              `yaml:"max_players"`
	Rate       GameRates        `yaml:"rate"`
	Internal   InternalEndpoint `yaml:"internal"`
	RabbitMQ   RabbitMQEndpoint `yaml:"rabbitmq"`
	HighRate   bool             `yaml:"high_rate"`
	Lua        GameLuaConfig    `yaml:"lua"`
}

type GameRates struct {
	Exp  int `yaml:"exp"`
	Drop int `yaml:"drop"`
	Meso int `yaml:"meso"`
}

type GameLuaConfig struct {
	AlwaysReload bool `yaml:"always_reload"`
}

func LoadLogin(path string) (*Login, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	var l Login
	if err := yaml.Unmarshal(data, &l); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	if l.Host == "" {
		l.Host = "0.0.0.0"
	}
	if l.Port == 0 {
		l.Port = 8484
	}
	if l.CatalogRetryIntervalSeconds <= 0 {
		l.CatalogRetryIntervalSeconds = 2
	}
	if l.CatalogRetryMaxAttempts <= 0 {
		l.CatalogRetryMaxAttempts = 30
	}
	return &l, nil
}

func LoadGame(path string) (*Game, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	var g Game
	if err := yaml.Unmarshal(data, &g); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	if g.Host == "" {
		g.Host = "0.0.0.0"
	}
	if g.Port == 0 {
		g.Port = 8485
	}
	if g.WzPath == "" {
		g.WzPath = "resources/wz"
	}
	if g.World == "" {
		g.World = "Scania"
	}
	if g.MaxPlayers == 0 {
		g.MaxPlayers = 1000
	}
	if g.Rate.Exp == 0 {
		g.Rate.Exp = 100
	}
	if g.Rate.Drop == 0 {
		g.Rate.Drop = 10
	}
	if g.Rate.Meso == 0 {
		g.Rate.Meso = 1
	}
	if g.RabbitMQ.IP == "" {
		g.RabbitMQ.IP = "127.0.0.1"
	}
	if g.RabbitMQ.Port == 0 {
		g.RabbitMQ.Port = 5672
	}
	if g.RabbitMQ.UID == "" {
		g.RabbitMQ.UID = "guest"
	}
	if g.RabbitMQ.PWD == "" {
		g.RabbitMQ.PWD = "guest"
	}
	if g.RabbitMQ.VHost == "" {
		g.RabbitMQ.VHost = "fm"
	}
	return &g, nil
}
