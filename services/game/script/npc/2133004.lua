-- NPC name (String.wz/Npc.img.xml): 스프라이트

local pq = require("script/lib/party_quest")

local PURPLE_STONE = 4001163

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		if wz == nil or wz.id ~= 930000500 then
			return
		end
		if not pq.has_item(me, PURPLE_STONE) then
			me:dialog(npc, "나에게 보라색 마력석을 찾아 오도록 해.")
			return
		end
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		map:clear_effect()
		pq.party_warp(sm, 930000600)
	end
}
