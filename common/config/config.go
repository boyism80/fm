package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

func FindConfigFilePath(p string) (string, error) {
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
	Host                     string `yaml:"host"`
	Port                     int    `yaml:"port"`
	TimeoutSeconds           int    `yaml:"timeout"`
	HeartbeatIntervalSeconds int    `yaml:"heartbeat_interval_seconds"`
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
	LoginInstanceID             string           `yaml:"login_instance_id"`
	Internal                    InternalEndpoint `yaml:"internal"`
	RabbitMQ                    RabbitMQEndpoint `yaml:"rabbitmq"`
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
	Dev        bool             `yaml:"dev"`
	Lua        GameLuaConfig    `yaml:"lua"`
	Duey       GameDueyConfig   `yaml:"duey"`
}

type GameRates struct {
	Exp  int `yaml:"exp"`
	Drop int `yaml:"drop"`
	Meso int `yaml:"meso"`
}

type GameLuaConfig struct {
	AlwaysReload bool `yaml:"always_reload"`
}

type GameDueyConfig struct {
	IdentityPrompt bool `yaml:"identity_prompt"`
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
	if l.RabbitMQ.IP == "" {
		l.RabbitMQ.IP = "127.0.0.1"
	}
	if l.RabbitMQ.Port == 0 {
		l.RabbitMQ.Port = 5672
	}
	if l.RabbitMQ.UID == "" {
		l.RabbitMQ.UID = "guest"
	}
	if l.RabbitMQ.PWD == "" {
		l.RabbitMQ.PWD = "guest"
	}
	if l.RabbitMQ.VHost == "" {
		l.RabbitMQ.VHost = "fm"
	}
	return &l, nil
}

type CashShop struct {
	Host       string           `yaml:"host"`
	Port       int              `yaml:"port"`
	CashShopID int              `yaml:"cash_shop_id"`
	WorldID    int              `yaml:"world_id"`
	WzPath     string           `yaml:"wz_path"`
	Internal   InternalEndpoint `yaml:"internal"`
}

func LoadCashShop(path string) (*CashShop, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	var c CashShop
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	if c.Host == "" {
		c.Host = "0.0.0.0"
	}
	if c.Port == 0 {
		c.Port = 8596
	}
	if c.WzPath == "" {
		c.WzPath = "resources/wz"
	}
	return &c, nil
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

type BotLogin struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Bot struct {
	Login          BotLogin `yaml:"login"`
	WorldID        int      `yaml:"world_id"`
	Channel        int      `yaml:"channel"`
	Seats          int      `yaml:"seats"`
	Password       string   `yaml:"password"`
	TimeoutMs      int      `yaml:"timeout_ms"`
	SuiteTimeoutMs int      `yaml:"suite_timeout_ms"`
	ScriptDir      string   `yaml:"script_dir"`
	GameScriptDir  string   `yaml:"game_script_dir"`
	ReportDir      string   `yaml:"report_dir"`
	WzPath         string   `yaml:"wz_path"`
}

func LoadBot(path string) (*Bot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	var b Bot
	if err := yaml.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	if b.Login.Host == "" {
		b.Login.Host = "127.0.0.1"
	}
	if b.Login.Port == 0 {
		b.Login.Port = 8484
	}
	if b.Seats == 0 {
		b.Seats = 3
	}
	if b.Password == "" {
		b.Password = "bot-pass"
	}
	if b.TimeoutMs == 0 {
		b.TimeoutMs = 10000
	}
	if b.SuiteTimeoutMs == 0 {
		b.SuiteTimeoutMs = 180000
	}
	if b.ScriptDir == "" {
		b.ScriptDir = "script/integration"
	}
	if b.GameScriptDir == "" {
		b.GameScriptDir = "../game/script"
	}
	if b.ReportDir == "" {
		b.ReportDir = "out"
	}
	if b.WzPath == "" {
		b.WzPath = "resources/wz"
	}
	if err := b.Validate(); err != nil {
		return nil, fmt.Errorf("config %q: %w", path, err)
	}
	return &b, nil
}

func (b *Bot) Validate() error {
	if b.Login.Host == "" {
		return fmt.Errorf("login.host is empty")
	}
	if b.Login.Port <= 0 {
		return fmt.Errorf("login.port must be positive")
	}
	if b.WorldID < 0 {
		return fmt.Errorf("world_id must be >= 0")
	}
	if b.Channel < 0 {
		return fmt.Errorf("channel must be >= 0")
	}
	if b.Seats < 1 {
		return fmt.Errorf("seats must be >= 1")
	}
	if b.Password == "" {
		return fmt.Errorf("password is empty")
	}
	if b.TimeoutMs <= 0 {
		return fmt.Errorf("timeout_ms must be positive")
	}
	if b.SuiteTimeoutMs <= 0 {
		return fmt.Errorf("suite_timeout_ms must be positive")
	}
	if b.ScriptDir == "" {
		return fmt.Errorf("script_dir is empty")
	}
	if b.WzPath == "" {
		return fmt.Errorf("wz_path is empty")
	}
	return nil
}
