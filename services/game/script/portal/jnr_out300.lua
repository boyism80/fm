local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil or not pq.is_leader(me) then
			me:notice("파티장이 먼저 이 포탈을 사용해야 합니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		pq.party_warp(sm, 926110400)
	end
}
