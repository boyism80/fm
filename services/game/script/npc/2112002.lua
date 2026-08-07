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
			me:dialog(npc, "흐음.. 자네들의 설득이 조금은 와닿는군. 내 연구를 인정받도록 도와주겠다고 했지? 일단 따라와 보게.")
			if sm:get_property("persuade_urete") == "1" then
				pq.party_exp(sm, 10500)
			end
			local dest = map_id + 100
			pq.party_warp(sm, dest)
		end
	end
}
