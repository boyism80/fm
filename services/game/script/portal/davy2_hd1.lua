local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		if sm:get_property("hd_202") ~= "" then
			me:notice("이번 스테이지에서는 이상 입장할 수 없습니다.", Msg.PinkText)
			return
		end
		if not pq.is_leader(me) then
			me:notice("파티장이 입장해야 합니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		pq.party_warp(sm, 925100202)
	end
}
