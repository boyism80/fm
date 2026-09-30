-- NPC name (String.wz/Npc.img.xml): 기함 파이어 올드 폭스 지원 AI

local spots = {
	{ 802000611, 120, 250, "올드 폭스 - 갑판" },
	{ 802000610, 120, 250, "올드 폭스 - 입구" },
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
		local sel = me:dialog_list(npc, "기계를 조종 하는 지원 AI다.", choices)
		if sel == nil then
			return
		end
		local map_id = maps[sel]
		if map_id == nil or id2map(map_id) == nil then
			me:dialog(npc, "아직 갈 수 없는 곳입니다.")
			return
		end
		me:map(map_id, 0)
	end
}
