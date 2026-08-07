local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		local map = me:map()
		if map == nil then
			return
		end
		if pq.mob_count(map) == 0 then
			me:play_portal_sound()
			me:map(925100400)
		else
			me:notice("아직 이 포탈을 이용할 수 없습니다.", Msg.PinkText)
		end
	end
}
