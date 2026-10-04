local M = {}

local GROUP_NAME = "henesys_party_quest"
local MOON_BUNNY_ID = 9300061
local SPAWN_POS = { -183, -187 }

function M.plant_seed(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:message("씨앗이 하나 심어졌습니다.")
	local player = reactor:trigger()
	if player == nil then
		log("plant_seed: trigger is nil")
		return
	end
	local sm = player:state_machine()
	if sm == nil then
		log("plant_seed: state_machine is nil for char", player:id())
		return
	end
	local group = sm:group()
	if group == nil then
		log("plant_seed: sm group is nil")
		return
	end
	if group:name() ~= GROUP_NAME then
		log("plant_seed: unexpected group", group:name(), "expected", GROUP_NAME)
		return
	end
	local stage = (tonumber(group:get_property("stage")) or 0) + 1
	group:set_property("stage", tostring(stage))
	local moon = map:find_reactor_name("fullmoon")
	if moon ~= nil then
		moon:hit(moon:state() + 1)
	end
	if stage == 6 then
		map:message("월묘를 보호하세요!")
		map:set_respawn(true)
		map:respawn({ include_one_time = true })
		map:spawn_mob(MOON_BUNNY_ID, SPAWN_POS[1], SPAWN_POS[2])
	end
end

return M
