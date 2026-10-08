-- NPC name (String.wz/Npc.img.xml): 용감한 어린양

local MATERIAL = 3980002

return {
	on_click = function(me, npc)
		local count = me:records():get("gacha.draws")
		me:dialog(npc, "#i" .. MATERIAL .. "# #b#z" .. MATERIAL .. "##k으로 캐시아이템을 뽑아보세요!\r\n" .. (20 - count) .. "번 더 뽑으면 리버스코인을 추가로 지급합니다.\r\n")
	end
}
