local pq = require("script/lib/party_quest")

local LIFE_GRASS = 4001055

return {
	on_enter = function(me)
		if not pq.is_leader(me) then
			me:notice("파티장이 생명의 풀을 갖고 이 포탈을 사용할 수 있습니다.", Msg.PinkText)
			return
		end
		if not pq.has_item(me, LIFE_GRASS, 1) then
			me:notice("생명을 풀을 얻어야 나갈 수 있습니다.", Msg.PinkText)
			return
		end
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		me:play_portal_sound()
		pq.party_warp(sm, 920010100)
	end
}
