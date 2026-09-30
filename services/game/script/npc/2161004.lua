-- NPC name (String.wz/Npc.img.xml): 머트

local spots = {
	{ 211060300, 130, 250, "사자왕의 성 - 성벽 아래2" },
	{ 211060500, 130, 250, "사자왕의 성 - 성벽 아래3" },
	{ 211060700, 130, 250, "사자왕의 성 - 성벽 아래4" },
	{ 211060900, 130, 250, "사자왕의 성 - 성벽 아래5" },
	{ 211060610, 130, 250, "사자왕의 성 - 낮은 성벽" },
	{ 211060830, 130, 250, "사자왕의 성 - 높은 성벽" },
}

return {
	on_click = function(me, npc)
		local choices = {}
		local maps = {}
		local level = me:level()
		for _, spot in ipairs(spots) do
			if level >= spot[2] and level <= spot[3] then
				table.insert(choices, spot[4] .. " | #r레벨 : " .. spot[2] .. " ~ " .. spot[3])
				table.insert(maps, spot[1])
			end
		end
		local sel = me:dialog_list(npc, "다른 곳 으로 이동하고 싶으신가요?", choices)
		if sel == nil then
			return
		end
		local map_id = maps[sel]
		if map_id == nil then
			return
		end
		me:map(map_id, 0)
	end
}
