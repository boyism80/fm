package main

import (
	"fmt"
	"os"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/luax"
	"github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/game/actorx"
	"github.com/boyism80/fm/game/data"
	"github.com/boyism80/fm/game/entity"
	lua "github.com/yuin/gopher-lua"
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

func init() {
	luax.RegisterOnCreateHook(func(L *lua.LState) {
		luax.RegisterLuaType[*entity.Object](L)
		luax.RegisterLuaDerivedType[*entity.Life, *entity.Object](L)
		luax.RegisterLuaDerivedType[*entity.Character, *entity.Life](L)
	})
}

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return actorx.NewGameServerActor() })
	pid := system.Root.Spawn(props)

	serverCtx := context.NewServerContext(data.NewResources(), map[uint32]*actor.PID{}, pid)
	system.Root.Send(pid, &msg.StartListening{Port: config.Server.Game.Port, ServerCtx: serverCtx})

	select {}
}
