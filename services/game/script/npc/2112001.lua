-- NPC name (String.wz/Npc.img.xml): 유레테

local pq = require("script/lib/party_quest")

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz().id
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		if map_id == 926100500 or map_id == 926110500 then
			me:dialog(npc, "실험은 실패로군...")
			pq.party_warp(sm, map_id + 100)
		end
	end
}
