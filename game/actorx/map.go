package actorx

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/data"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"

	common_msg "github.com/boyism80/fm/common/msg"
)

type mobSpawn struct {
	Spec          *data.MobSpawnSpec
	NextSpawnTime time.Time
	Spawned       bool
}

type MapActor struct {
	sequence    uint32
	handler     *handler.MessageHandler
	Spec        *data.MapSpec
	ctx         *context.ServerContext
	objects     map[constant.ObjectType]map[uint32]*actor.PID
	mobSpawners map[uint32]*mobSpawn
	scheduler   *scheduler.TimerScheduler
}

func NewMapActorProps(ctx actor.Context, serverCtx *context.ServerContext, spec *data.MapSpec) *actor.Props {
	return actor.PropsFromProducer(func() actor.Actor {

		actor := &MapActor{
			handler: handler.NewMessageHandler(),
			Spec:    spec,
			ctx:     serverCtx,
			objects: map[constant.ObjectType]map[uint32]*actor.PID{
				constant.ObjectTypeItem:      make(map[uint32]*actor.PID),
				constant.ObjectTypeNpc:       make(map[uint32]*actor.PID),
				constant.ObjectTypeMob:       make(map[uint32]*actor.PID),
				constant.ObjectTypeCharacter: make(map[uint32]*actor.PID),
			},
			mobSpawners: map[uint32]*mobSpawn{},
			scheduler:   scheduler.NewTimerScheduler(ctx),
		}

		for _, npcSpec := range spec.NpcSpawns {
			actor.sequence++
			props := NewNpcActorProps(ctx, &npcSpec, actor.sequence)
			pid := ctx.Spawn(props)
			actor.objects[constant.ObjectTypeNpc][actor.sequence] = pid
		}

		for _, mobSpawnSpec := range spec.MobSpawns {
			actor.sequence++
			actor.mobSpawners[actor.sequence] = &mobSpawn{
				Spec:          &mobSpawnSpec,
				NextSpawnTime: time.Time{},
				Spawned:       false,
			}
		}

		handler.RegisterHandler(ctx, actor, actor.handler, onMapStarted)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapEnter)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapLeave)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapPIDList)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapBroadcastRange)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSpawnItem)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSpawnMeso)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapItemLoot)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapRemoveItem)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapChange)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSpawnNpc)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSendMessage)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSpawnMob)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapDieMob)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapClearMobs)

		return actor
	})
}

func (state *MapActor) Receive(ctx actor.Context) {
	state.handler.Handle(ctx)
}

func (state *MapActor) SpawnMobs(ctx actor.Context) {

	if len(state.objects[constant.ObjectTypeCharacter]) == 0 {
		return
	}

	for spawnId, mobSpawner := range state.mobSpawners {
		if mobSpawner.Spawned {
			continue
		}

		if mobSpawner.NextSpawnTime.After(time.Now()) {
			continue
		}

		mobSpec, ok := state.ctx.Resources.Monsters[mobSpawner.Spec.ID]
		if !ok {
			continue
		}

		props := actor.PropsFromProducer(func() actor.Actor {
			return NewMobActor(ctx, *state.ctx, entity.Mob{
				Life: entity.Life{
					Object: entity.Object{
						Position: state.Spec.DropPoint(mobSpawner.Spec.Position),
					},
					Hp:     uint16(mobSpec.MaxHP),
					Mp:     uint16(mobSpec.MaxMP),
					MaxHp:  uint16(mobSpec.MaxHP),
					MaxMp:  uint16(mobSpec.MaxMP),
					Stance: 5,
				},
				Spec:     mobSpec,
				ID:       spawnId,
				Foothold: mobSpawner.Spec.Foothold,
			}, ctx.Self())
		})
		pid := ctx.Spawn(props)
		state.objects[constant.ObjectTypeMob][spawnId] = pid
		mobSpawner.Spawned = true
		ctx.Send(pid, &msg.Spawn{})
	}
}

func onMapStarted(ctx actor.Context, state *MapActor, m *actor.Started) {
	state.scheduler.SendRepeatedly(time.Second*8, time.Second*8, ctx.Self(), &msg.MapSpawnMob{})
}

func onMapEnter(ctx actor.Context, state *MapActor, m *msg.EnterMap) {

	// 맵에 플레이어를 추가
	state.objects[constant.ObjectTypeCharacter][m.ID] = m.PID
	if len(state.objects[constant.ObjectTypeCharacter]) == 1 {
		state.SpawnMobs(ctx)
	}

	ctx.Send(m.PID, &msg.CharacterMapChanged{
		MID:        state.Spec.ID,
		Map:        ctx.Self(),
		Init:       m.Init,
		SpawnPoint: m.SpawnPoint,
	})
}

func onMapLeave(ctx actor.Context, state *MapActor, m *msg.LeaveMap) {
	delete(state.objects[constant.ObjectTypeCharacter], m.ID)

	for _, pid := range state.objects[constant.ObjectTypeCharacter] {
		ctx.Send(pid, &common_msg.SendProtocol{
			Protocol: &resp.LeavePlayer{
				ID: m.ID,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		})
	}
}

func onMapPIDList(ctx actor.Context, state *MapActor, m *msg.MapPIDList) {
	var pids []*actor.PID
	for _, pid := range state.objects[constant.ObjectTypeCharacter] {
		pids = append(pids, pid)
	}
	ctx.Send(ctx.Sender(), &msg.MapPIDList{Targets: pids})
}

func onMapBroadcastRange(ctx actor.Context, state *MapActor, m *msg.MapBroadcastRange) {
	// 나중에 섹터 추가하고 섹터 찾아서 섹터 액터한테 던짐

	for _, pid := range state.objects[constant.ObjectTypeCharacter] {
		if m.ExceptSelf && pid == m.Sender {
			continue
		}

		ctx.Send(pid, m.Message)
	}
}

func onMapSpawnItem(ctx actor.Context, state *MapActor, m *msg.MapSpawnItem) {
	state.sequence++
	drop := m.Item.GetDrop()
	drop.Position = state.Spec.DropPoint(drop.SpawnedPoint)
	drop.ID = state.sequence
	props := actor.PropsFromProducer(func() actor.Actor {
		return NewItemActor(ctx,
			state.ctx,
			m.Item,
			ctx.Self())
	})
	pid := ctx.Spawn(props)
	state.objects[constant.ObjectTypeItem][state.sequence] = pid

	ctx.Send(pid, &msg.Spawn{})
}

func onMapSpawnMeso(ctx actor.Context, state *MapActor, m *msg.MapSpawnMeso) {
	state.sequence++
	dropPoint := state.Spec.DropPoint(m.SpawnedPoint)
	props := actor.PropsFromProducer(func() actor.Actor {
		return NewMesoActor(ctx,
			state.ctx,
			entity.Meso{
				Drop: &entity.Drop{
					Object: &entity.Object{
						Position: dropPoint,
					},
					ID:           state.sequence,
					SpawnedPoint: m.SpawnedPoint,
					DropType:     constant.DropTypeFFA,
					Owner:        m.OwnerID,
				},
				Count: m.Count,
			},
			ctx.Self(),
			dropPoint)
	})
	pid := ctx.Spawn(props)
	state.objects[constant.ObjectTypeItem][state.sequence] = pid
	ctx.Send(pid, &msg.Spawn{})
}

func onMapItemLoot(ctx actor.Context, state *MapActor, m *msg.MapItemLoot) {

	pid, ok := state.objects[constant.ObjectTypeItem][m.OID]
	if !ok {
		ctx.Send(m.Actor, &msg.CharacterLootFailed{
			OID: m.OID,
		})
	} else {
		ctx.Send(pid, &msg.ItemLooting{
			Actor:       m.Actor,
			Position:    m.Position,
			CharacterId: m.CharacterId,
		})
	}
}

func onMapRemoveItem(ctx actor.Context, state *MapActor, m *msg.MapRemoveItem) {
	for _, pid := range state.objects[constant.ObjectTypeCharacter] {
		ctx.Send(pid, &common_msg.SendProtocol{
			Protocol: &resp.RemoveItem{
				Mode:        m.Mode,
				OID:         m.OID,
				CharacterId: m.CharacterId,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		})
	}

	delete(state.objects[constant.ObjectTypeItem], m.OID)
	ctx.Stop(m.Actor)
}

func onMapChange(ctx actor.Context, state *MapActor, m *msg.MapChange) {
	onMapLeave(ctx, state, &msg.LeaveMap{
		ID: m.CharacterId,
	})

	ctx.Send(m.To, &msg.EnterMap{
		ID:         m.CharacterId,
		PID:        m.Sender,
		SpawnPoint: m.SpawnPoint,
		Init:       false,
	})
}

func onMapSpawnNpc(ctx actor.Context, state *MapActor, m *msg.MapNotifyCharacterWarped) {

	for _, v := range state.objects {
		for _, pid := range v {
			ctx.Send(pid, &msg.Warped{
				Sender: m.Sender,
			})
		}
	}
}

func onMapSendMessage(ctx actor.Context, state *MapActor, m *msg.SendMessage) {
	pid, ok := state.objects[constant.ObjectTypeNpc][m.OID]
	if !ok {
		return
	}

	ctx.Send(pid, m.Message)
}

func onMapSpawnMob(ctx actor.Context, state *MapActor, m *msg.MapSpawnMob) {
	state.SpawnMobs(ctx)
}

func onMapDieMob(ctx actor.Context, state *MapActor, m *msg.MapDieMob) {
	pid, ok := state.objects[constant.ObjectTypeMob][m.OID]
	if !ok {
		return
	}

	ctx.Stop(pid)
	delete(state.objects[constant.ObjectTypeMob], m.OID)

	spawner, ok := state.mobSpawners[m.OID]
	if ok {
		spawner.Spawned = false
		nextSpawnDuration := spawner.Spec.MobTime
		if nextSpawnDuration == 0 {
			nextSpawnDuration = constant.DefaultMobSpawnTime
		}
		spawner.NextSpawnTime = time.Now().Add(nextSpawnDuration)
	}

	ctx.Send(ctx.Self(), &msg.MapBroadcastRange{
		Sender: m.Sender,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.DieMob{
				OID:           m.OID,
				AnimationType: m.AnimationType,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
		Pivot:      m.Position,
		ExceptSelf: true,
	})
}

func onMapClearMobs(ctx actor.Context, state *MapActor, m *msg.MapClearMobs) {
	for _, v := range state.objects[constant.ObjectTypeMob] {
		ctx.Send(v, &msg.MobKill{
			AnimationType: m.AnimationType,
		})
	}
}
