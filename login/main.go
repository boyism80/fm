package main

import (
	"fmt"
	"os"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/login/actor"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Login struct {
			Port int `yaml:"port"`
		} `yaml:"login"`

		Game struct {
			Port int `yaml:"port"`
		} `yaml:"game"`
	} `yaml:"server"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	system := protoactor.NewActorSystem()
	props := protoactor.PropsFromProducer(func() protoactor.Actor { return actor.NewLoginServerActor() })
	pid := system.Root.Spawn(props)

	serverCtx := context.NewServerContext(nil, map[uint32]*protoactor.PID{}, pid)
	system.Root.Send(pid, &msg.StartListening{Port: config.Server.Login.Port, ServerCtx: serverCtx})

	select {}
}
