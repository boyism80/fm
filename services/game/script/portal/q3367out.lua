local pq = require("script/lib/party_quest")

local BOOK = 4031797
local KEEP = 20

return {
	on_enter = function(me)
		local n = pq.item_count(me, BOOK)
		if n > KEEP then
			me:rmitem(BOOK, n - KEEP)
		end
		me:play_portal_sound()
		me:map(926130100, 5)
	end
}
