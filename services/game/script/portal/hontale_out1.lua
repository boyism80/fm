-- Portal (old/scripts/portal/hontale_out1.js): 생명의동굴 동굴의 사잇길

local pq = require("script/lib/party_quest")

local KEY_ITEMS = {
	4001087,
	4001088,
	4001089,
	4001090,
	4001091,
	4001092,
	4001093,
}

function on_enter(me)
	for _, id in ipairs(KEY_ITEMS) do
		pq.remove_all(id, me)
	end
	me:play_portal_sound()
	me:map(240050400)
end
