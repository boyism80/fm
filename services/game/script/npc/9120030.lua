-- NPC name (String.wz/Npc.img.xml): 마르

local spots = {
	{ 802000101, 120, 250, "망각의숲 - 카뮤나" },
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
		local sel = me:dialog_list(npc, "돌아가실 건가요?", choices)
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
