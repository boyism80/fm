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
	sequence        uint32
	handler         *handler.MessageHandler
	Spec            *data.MapSpec
	ctx             *context.ServerContext
	objects         map[constant.ObjectType]map[uint32]*actor.PID
	mobSpawners     map[uint32]*mobSpawn
	scheduler       *scheduler.TimerScheduler
	controllerTable *entity.ControllerTable
}

func OnMobControllerChange(ctx actor.Context, mob *actor.PID, before *actor.PID, after *actor.PID) {

	ctx.Send(mob, &msg.MobControllerChange{
		Before: before,
		After:  after,
	})
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
			mobSpawners:     map[uint32]*mobSpawn{},
			scheduler:       scheduler.NewTimerScheduler(ctx),
			controllerTable: entity.NewControllerTable(OnMobControllerChange),
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
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSpawnItems)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSpawnMeso)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapItemLoot)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapRemoveItem)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapChange)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSpawnNpc)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSendMessage)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapRepeatSpawnMobs)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapDieMob)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapClearMobs)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapMoveMob)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSpawningMob)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapSpawnedMob)
		handler.RegisterHandler(ctx, actor, actor.handler, onMapCharacterAttack)

		return actor
	})
}

func (state *MapActor) Receive(ctx actor.Context) {
	state.handler.Handle(ctx)
}

func (state *MapActor) SpawnMob(ctx actor.Context, mobId uint32, oid uint32, foothold int16, position types.Vector2[int16]) bool {
	mobSpec, ok := state.ctx.Resources.Monsters[mobId]
	if !ok {
		return false
	}

	props := actor.PropsFromProducer(func() actor.Actor {
		spawnPoint, ok := state.Spec.DropPoint(position)
		if !ok {
			spawnPoint = position
		}
		return NewMobActor(ctx, state.ctx, entity.Mob{
			Life: entity.Life{
				Object: entity.Object{
					Position: spawnPoint,
				},
				Hp:     uint16(mobSpec.MaxHP),
				Mp:     uint16(mobSpec.MaxMP),
				MaxHp:  uint16(mobSpec.MaxHP),
				MaxMp:  uint16(mobSpec.MaxMP),
				Stance: 5,
			},
			Spec:     mobSpec,
			OID:      oid,
			Foothold: foothold,
		}, ctx.Self())
	})
	pid := ctx.Spawn(props)
	state.objects[constant.ObjectTypeMob][oid] = pid
	ctx.Send(pid, &msg.Spawn{})
	return true
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

		if !state.SpawnMob(ctx, mobSpawner.Spec.ID, spawnId, mobSpawner.Spec.Foothold, mobSpawner.Spec.Position) {
			continue
		}

		mobSpawner.Spawned = true
	}
}

func onMapStarted(ctx actor.Context, state *MapActor, m *actor.Started) {
	state.scheduler.SendRepeatedly(time.Second*8, time.Second*8, ctx.Self(), &msg.MapRepeatSpawnMobs{})
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

	state.controllerTable.EnterPlayer(ctx, m.PID)
}

func onMapLeave(ctx actor.Context, state *MapActor, m *msg.LeaveMap) {

	if pid, ok := state.objects[constant.ObjectTypeCharacter][m.ID]; ok {
		state.controllerTable.LeavePlayer(ctx, pid)
		delete(state.objects[constant.ObjectTypeCharacter], m.ID)
	}

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

func (state *MapActor) broadcast(ctx actor.Context, sender *actor.PID, m any, exceptSelf bool, excepts map[*actor.PID]struct{}) {

	for _, pid := range state.objects[constant.ObjectTypeCharacter] {
		if exceptSelf && pid == sender {
			continue
		}

		if excepts != nil {
			if _, ok := excepts[pid]; ok {
				continue
			}
		}

		ctx.Send(pid, m)
	}
}

func onMapBroadcastRange(ctx actor.Context, state *MapActor, m *msg.MapBroadcast) {
	// 나중에 섹터 추가하고 섹터 찾아서 섹터 액터한테 던짐

	state.broadcast(ctx, m.Sender, m.Message, m.ExceptSelf, m.Excepts)
}

func onMapSpawnItem(ctx actor.Context, state *MapActor, m *msg.MapSpawnItem) {
	state.sequence++
	drop := m.Item.GetDrop()
	dropPoint, ok := state.Spec.DropPoint(drop.Position)
	if !ok {
		dropPoint = drop.SpawnedPoint
	}
	drop.Position = dropPoint
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

func onMapSpawnItems(ctx actor.Context, state *MapActor, m *msg.MapSpawnItems) {
	spawnPoint := m.Position
	spacing := int16(15)

	for i, dropped := range m.Items {
		destPoint := spawnPoint
		if len(m.Items) > 1 {
			offset := spacing * int16(i/2+1)
			if i%2 == 0 {
				destPoint.X += offset
			} else {
				destPoint.X -= offset
			}
		}

		if dropped.IsMeso() {
			meso := dropped.(*entity.Meso)
			onMapSpawnMeso(ctx, state, &msg.MapSpawnMeso{
				Count:        dropped.GetCount32(),
				SpawnedPoint: spawnPoint,
				DestPoint:    destPoint,
				Owner:        m.Owner,
				OwnerID:      meso.GetDrop().Owner,
				DropType:     meso.GetDrop().DropType,
			})
		} else {
			item, ok := dropped.(entity.Item)
			if !ok {
				continue
			}

			drop := item.GetDrop()
			drop.SpawnedPoint = spawnPoint
			drop.Position = destPoint
			onMapSpawnItem(ctx, state, &msg.MapSpawnItem{
				Item:    item,
				Owner:   m.Owner,
				OwnerID: item.GetDrop().Owner,
			})
		}
	}
}

func onMapSpawnMeso(ctx actor.Context, state *MapActor, m *msg.MapSpawnMeso) {
	state.sequence++
	dropPoint, ok := state.Spec.DropPoint(m.DestPoint)
	if !ok {
		dropPoint = m.SpawnedPoint
	}
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
					DropType:     m.DropType,
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

func onMapRepeatSpawnMobs(ctx actor.Context, state *MapActor, m *msg.MapRepeatSpawnMobs) {
	state.SpawnMobs(ctx)
}

func onMapDieMob(ctx actor.Context, state *MapActor, m *msg.MapDieMob) {
	pid, ok := state.objects[constant.ObjectTypeMob][m.OID]
	if !ok {
		return
	}

	state.controllerTable.LeaveMob(ctx, pid)

	spawner, ok := state.mobSpawners[m.OID]
	if ok {
		spawner.Spawned = false
		nextSpawnDuration := spawner.Spec.MobTime
		if nextSpawnDuration == 0 {
			nextSpawnDuration = constant.DefaultMobSpawnTime
		}
		spawner.NextSpawnTime = time.Now().Add(nextSpawnDuration)
	}

	ctx.Send(ctx.Self(), &msg.MapBroadcast{
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

	delete(state.objects[constant.ObjectTypeMob], m.OID)
	ctx.Stop(pid)
}

func onMapClearMobs(ctx actor.Context, state *MapActor, m *msg.MapClearMobs) {
	for _, v := range state.objects[constant.ObjectTypeMob] {
		ctx.Send(v, &msg.MobKill{
			AnimationType: m.AnimationType,
		})
	}
}

func onMapMoveMob(ctx actor.Context, state *MapActor, m *msg.MapMoveMob) {
	pid, ok := state.objects[constant.ObjectTypeMob][m.OID]
	if !ok {
		return
	}

	controller, ok := state.controllerTable.GetController(pid)
	if !ok {
		return
	}

	if m.Sender != controller {
		if m.Unknown2 {
			// TODO: stopControl
		} else {
			// TODO: switchControl
		}
		return
	}

	ctx.Send(pid, &msg.MobMove{
		MoveMob: m.MoveMob,
		Sender:  m.Sender,
	})
}

func onMapSpawningMob(ctx actor.Context, state *MapActor, m *msg.MapSpawningMob) {
	state.sequence++
	foothold, ok := state.Spec.Footholds.Find(m.Position)
	if !ok {
		return
	}
	state.SpawnMob(ctx, m.MobId, state.sequence, foothold.ID, m.Position)
}

func onMapSpawnedMob(ctx actor.Context, state *MapActor, m *msg.MapSpawnedMob) {

	state.controllerTable.EnterMob(ctx, m.PID)
}

func onMapCharacterAttack(ctx actor.Context, state *MapActor, m *msg.MapCharacterAttack) {

	for _, damage := range m.AttackInfo.Damages {
		pid, ok := state.objects[constant.ObjectTypeMob][damage.OID]
		if !ok {
			continue
		}

		ctx.Send(pid, &msg.MobDamaged{
			Sender:      m.Sender,
			CharacterId: m.CharacterId,
			DamagePairs: damage.DamagePairs,
		})
	}

	state.broadcast(ctx, m.Sender, &common_msg.SendProtocol{
		Protocol: &resp.Attack{
			CharacterId: m.CharacterId,
			AttackInfo:  m.AttackInfo,
			SkillLevel:  0,
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	}, true, nil)
}
