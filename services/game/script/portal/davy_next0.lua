local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		local map = me:map()
		if map == nil then
			return
		end
		if pq.mob_count(map) == 0 then
			me:play_portal_sound()
			me:map(925100100)
		else
			me:notice("이 포탈은 잠겨 있습니다.", Msg.PinkText)
		end
	end
}
