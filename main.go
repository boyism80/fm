package main

import (
	"fmt"
	"os"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/model"
	"github.com/boyism80/fm/msg"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
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

	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return model.NewMainActor() })
	pid := system.Root.Spawn(props)
	system.Root.Send(pid, &msg.StartListening{Port: config.Server.Port})

	select {}
}
