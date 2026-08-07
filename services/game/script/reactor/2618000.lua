-- Reactor name (Reactor.wz/2618000.img.xml): 비커

local pq = require("script/lib/party_quest")
local rj = require("script/lib/romeo_juliet")

return {
	on_reactor = function(reactor)
		if reactor:state() < 7 then
			return
		end
		local map = reactor:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz().id
		local player = reactor:trigger()
		local sm = player ~= nil and player:state_machine() or nil
		if sm == nil then
			return
		end
		map:message("한 개의 비커를 가득 채웠습니다.")
		local door_name = "jnr2_door"
		if map_id == 926100100 then
			door_name = "rnj2_door"
		end
		local door = map:reactor_by_name(door_name)
		local stage3 = (tonumber(sm:get_property("stage3")) or 0) + 1
		sm:set_property("stage3", tostring(stage3))
		if door ~= nil then
			door:hit(door:state() + 1)
		end
		if stage3 == 3 then
			pq.party_exp(sm, 20000)
			map:block_gen(true)
			map:kill_all_mobs()
			rj.clear_fx(map)
		end
	end
}
