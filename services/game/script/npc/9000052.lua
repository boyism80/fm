-- NPC name (String.wz/Npc.img.xml): 용감한 어린양

local MATERIAL = 3980002

return {
	on_click = function(me, npc)
		local q = me:quest(5500001)
		if not q:started() then
			q:start("0")
		end
		local count = tonumber(q:record()) or 0
		me:dialog(npc, "#i" .. MATERIAL .. "# #b#z" .. MATERIAL .. "##k으로 캐시아이템을 뽑아보세요!\r\n" .. (20 - count) .. "번 더 뽑으면 리버스코인을 추가로 지급합니다.\r\n")
	end
}
