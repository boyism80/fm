package entity

import "github.com/asynkron/protoactor-go/actor"

type ControllerTable struct {
	controllers        map[*actor.PID]struct{}
	mobs               map[*actor.PID]struct{}
	controller2mob     map[*actor.PID]map[*actor.PID]struct{}
	mob2controller     map[*actor.PID]*actor.PID
	onControllerChange func(ctx actor.Context, mob *actor.PID, controller *actor.PID)
}

func NewControllerTable(onControllerChange func(ctx actor.Context, mob *actor.PID, controller *actor.PID)) *ControllerTable {
	return &ControllerTable{
		onControllerChange: onControllerChange,
		controllers:        make(map[*actor.PID]struct{}),
		mobs:               make(map[*actor.PID]struct{}),
		controller2mob:     make(map[*actor.PID]map[*actor.PID]struct{}),
		mob2controller:     make(map[*actor.PID]*actor.PID),
	}
}

func (t *ControllerTable) EnterPlayer(ctx actor.Context, controller *actor.PID) {
	t.controllers[controller] = struct{}{}
	t.controller2mob[controller] = make(map[*actor.PID]struct{})

	// orphan 몹을 빠르게 조회해 바로 할당
	for mob := range t.mobs {
		if t.mob2controller[mob] == nil {
			t.assign(ctx, mob, controller)
		}
	}
}

func (t *ControllerTable) LeavePlayer(ctx actor.Context, controller *actor.PID) {
	// 다음 플레이어 찾기
	var next *actor.PID
	for pid := range t.controllers {
		if pid != controller {
			next = pid
			break
		}
	}

	// controller가 관리하던 모든 몹을 재할당
	for mob := range t.controller2mob[controller] {
		t.assign(ctx, mob, next)
	}

	delete(t.controllers, controller)
	delete(t.controller2mob, controller)
}

func (t *ControllerTable) EnterMob(ctx actor.Context, mob *actor.PID) {
	t.mobs[mob] = struct{}{}
	t.mob2controller[mob] = nil
}

func (t *ControllerTable) LeaveMob(ctx actor.Context, mob *actor.PID) {
	if cur := t.mob2controller[mob]; cur != nil {
		// O(1) 삭제
		delete(t.controller2mob[cur], mob)
		// 콜백: controller=nil
		if t.onControllerChange != nil {
			t.onControllerChange(ctx, mob, nil)
		}
	}
	delete(t.mobs, mob)
	delete(t.mob2controller, mob)
}

// assign은 몬스터→컨트롤러 매핑과 controller2mob 업데이트, 콜백 호출을 담당
func (t *ControllerTable) assign(ctx actor.Context, mob *actor.PID, controller *actor.PID) {
	// 이전 컨트롤러가 있으면 삭제
	if prev := t.mob2controller[mob]; prev != nil {
		delete(t.controller2mob[prev], mob)
	}
	// 새 컨트롤러 설정
	t.mob2controller[mob] = controller
	if controller != nil {
		if _, ok := t.controller2mob[controller]; !ok {
			t.controller2mob[controller] = map[*actor.PID]struct{}{}
		}
		t.controller2mob[controller][mob] = struct{}{}
	}
	// 콜백
	if t.onControllerChange != nil {
		t.onControllerChange(ctx, mob, controller)
	}
}
