-- NPC name (String.wz/Npc.img.xml): 컴퓨터

local maps = {
	190000000,
	191000000,
	192000000,
	195000000,
	196000000,
	197000000,
}

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "메이플스토리에 오신것을 환영합니다.", {
			"#m190000000#",
			"#m191000000#",
			"#m192000000#",
			"#m195000000#",
			"#m196000000#",
			"#m197000000#",
		})
		if sel == nil then
			return
		end
		local map_id = maps[sel]
		if map_id == nil then
			return
		end
		me:map(map_id)
	end
}
