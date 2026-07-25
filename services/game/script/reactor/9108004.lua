-- Reactor name (Reactor.wz/9108004.img.xml): 달맞이꽃 씨앗

local GROUP_NAME = "henesys_party_quest"
local MOON_BUNNY_ID = 9300061
local SPAWN_POS = { -183, -187 }

return {
	on_reactor = function(reactor, item)
		if item ~= nil then
			return
		end
		local map = reactor:map()
		if map == nil then
			return
		end
		map:message("씨앗이 하나 심어졌습니다.")
		local group = state_machine(GROUP_NAME)
		if group == nil then
			return
		end
		local stage = (tonumber(group:get_property("stage")) or 0) + 1
		group:set_property("stage", tostring(stage))
		local moon = map:reactor_by_name("fullmoon")
		if moon ~= nil then
			moon:hit(moon:state() + 1)
		end
		if stage == 6 then
			map:message("월묘를 보호하세요!")
			map:block_gen(true)
			map:respawn(true)
			map:spawn_mob(MOON_BUNNY_ID, SPAWN_POS[1], SPAWN_POS[2])
		end
	end
}
