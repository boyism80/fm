local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		for id = 4031780, 4031784 do
			pq.remove_all(id, me)
		end
		me:play_portal_sound()
		me:map(926130100, 4)
	end
}
