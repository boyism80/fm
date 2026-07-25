-- Reactor name (Reactor.wz/9108000.img.xml): 달맞이꽃 씨앗

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
		local player = reactor:trigger()
		if player == nil then
			log("on_reactor: trigger is nil")
			return
		end
		local sm = player:state_machine()
		if sm == nil then
			log("on_reactor: state_machine is nil for char", player:id())
			return
		end
		local group = sm:group()
		if group == nil then
			log("on_reactor: sm group is nil")
			return
		end
		if group:name() ~= GROUP_NAME then
			log("on_reactor: unexpected group", group:name(), "expected", GROUP_NAME)
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
