local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		local sm = me:state_machine()
		local map = me:map()
		if sm == nil or map == nil then
			return
		end
		if sm:get_property("stage2") == "3" and pq.mob_count(map) == 0 then
			me:play_portal_sound()
			me:map(925100200)
		else
			me:notice("아직 이 포탈을 이용할 수 없습니다.", Msg.PinkText)
		end
	end
}
