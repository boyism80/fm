-- NPC name (String.wz/Npc.img.xml): 가루다 사령관

local spots = {
	{ 802000101, 120, 250, "나가기" },
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
		if #choices == 0 then
			return
		end
		local sel = me:dialog_list(npc, "가루다가 당신을 마중나왔다.", choices)
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
