package main

import (
	"fmt"
	"log"
	"os"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/lua"
	"github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/game/actor"
	"github.com/boyism80/fm/game/data"
	"github.com/boyism80/fm/game/entity"
	raw_lua "github.com/yuin/gopher-lua"
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
	L := raw_lua.NewState()
	defer L.Close()

	lua.Register(L, "pid", map[string]raw_lua.LGFunction{})
	lua.RegisterLuaType[*entity.Object](L)
	lua.RegisterLuaDerivedType[*entity.Life, *entity.Object](L)    // Life ← Object
	lua.RegisterLuaDerivedType[*entity.Character, *entity.Life](L) // Monster ← Life

	ch := entity.NewDummyCharacter(1, "test", nil)
	ud := L.NewUserData()
	ud.Value = &ch
	L.SetMetatable(ud, L.GetTypeMetatable(ch.LuaTypeName()))
	L.SetGlobal("me", ud)

	const script = `
	-- Go 쪽에서 밀어넣은 ch(Userdata)를 me 라는 이름으로 사용
	local hp = me:hp()   -- Character 메타테이블에 바인딩된 hp() 메서드 호출
	print("====> me:hp() →", hp)
`
	if err := L.DoString(script); err != nil {
		log.Fatal(err)
	}

	config, err := loadConfig("config.yaml")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	system := protoactor.NewActorSystem()
	props := protoactor.PropsFromProducer(func() protoactor.Actor { return actor.NewGameServerActor() })
	pid := system.Root.Spawn(props)

	serverCtx := context.NewServerContext(data.NewResources(), map[uint32]*protoactor.PID{}, pid)
	system.Root.Send(pid, &msg.StartListening{Port: config.Server.Game.Port, ServerCtx: serverCtx})

	select {}
}
