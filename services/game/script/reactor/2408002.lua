-- Reactor name (Reactor.wz/2408002.img.xml): 혼테일 PQ 열쇠 워프

local HUB_MAP = 240050100
local DROP_POS = { 874, 5 }

local DROPS = {
	[240050101] = { item = 4001088, door_state = 1 },
	[240050102] = { item = 4001089, door_state = 0 },
	[240050103] = { item = 4001090, door_state = 1 },
	[240050104] = { item = 4001091, door_state = 0 },
}

function on_reactor_2408002(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	local wz = map:wz()
	if wz == nil then
		return
	end
	local drop = DROPS[wz.id]
	if drop == nil then
		return
	end
	local player = reactor:trigger()
	if player == nil then
		return
	end
	local sm = player:state_machine()
	if sm == nil or sm:group() == nil then
		return
	end
	local hub = sm:group():map(HUB_MAP)
	if hub == nil then
		return
	end
	hub:spawn_item(drop.item, 1, DROP_POS)
	local door = hub:reactor(2402002)
	if door ~= nil then
		door:hit(drop.door_state)
	end
	map:message("열쇠가 어디론가 사라졌습니다.")
	for _, p in ipairs(sm:players()) do
		if p ~= nil then
			p:notice("반짝이는 빛과 함께 어딘가에서 열쇠가 나타났습니다.", Msg.LightBlueText)
		end
	end
end
