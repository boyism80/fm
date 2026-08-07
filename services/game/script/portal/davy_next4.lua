local pq = require("script/lib/party_quest")

local BOSS_MAP = 925100500

return {
	on_enter = function(me)
		local map = me:map()
		if map == nil then
			return
		end
		local r1 = map:reactor_by_name("sMob1")
		local r2 = map:reactor_by_name("sMob2")
		local r3 = map:reactor_by_name("sMob3")
		local r4 = map:reactor_by_name("sMob4")
		local doors_ok = r1 ~= nil and r1:state() >= 1
			and r2 ~= nil and r2:state() >= 1
			and r3 ~= nil and r3:state() >= 1
			and r4 ~= nil and r4:state() >= 1
		if pq.mob_count(map) == 0 and doors_ok then
			if not pq.is_leader(me) then
				me:notice("파티장이 이 포탈을 이용해야 합니다.", Msg.PinkText)
				return
			end
			local sm = me:state_machine()
			if sm == nil then
				return
			end
			me:play_portal_sound()
			pq.party_warp(sm, BOSS_MAP)
		else
			me:notice("아직 이 포탈을 이용할 수 없습니다.", Msg.PinkText)
		end
	end
}
