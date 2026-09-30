-- NPC name (String.wz/Npc.img.xml): 입실 제어 장치

local spots = {
	{ 802000800, 120, 250, "롯폰기 백화점" },
	{ 802000801, 120, 250, "롯폿기 백화점 - 입구" },
	{ 802000802, 120, 250, "롯폰기 백화점 - 통로" },
	{ 802000803, 120, 250, "롯폰기 백화점 - [코어블레이즈]" },
	{ 802000810, 120, 250, "롯폰기 백화점 - 옥상 입구" },
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
		local sel = me:dialog_list(npc, "입실 장치다. 원하는 장소를 클릭 하면 이동시켜준다.", choices)
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
