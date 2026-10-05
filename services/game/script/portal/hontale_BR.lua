local pq = require("script/lib/party_quest")

local ALTAR = 2408003

return {
	on_enter = function(me)
		local map = me:map()
		local altar = map:reactor(ALTAR)
		if pq.mob_count(map) > 0 or altar == nil or altar:state() ~= 1 then
			me:message("아직 혼테일의 머리가 남아 있습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		pq.warp(me, map:wz():id() + 100)
	end
}
