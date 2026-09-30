-- NPC name (String.wz/Npc.img.xml): 클라라

local spots = {
	{ 104040000, 10, 20, "#e헤네시스 사냥터 #b[사냥터]" },
	{ 105050000, 20, 30, "#e던전 #b[사냥터]" },
	{ 540000100, 30, 35, "#e싱가포르 #b[사냥터]" },
	{ 800040100, 35, 45, "#e영주성 #b[사냥터]" },
	{ 550000200, 45, 55, "#e진흙 사면 #b[사냥터]" },
	{ 541010010, 55, 80, "#e유령선 #b[사냥터]" },
	{ 551020000, 80, 90, "#e환상의 테마공원#b[사냥터]" },
	{ 541020000, 90, 120, "#e버려진 도시 #b[사냥터]" },
	{ 240040510, 120, 130, "#e죽은 용의 둥지 #b[사냥터]" },
	{ 223030100, 130, 250, "#e판타스틱 테마 파크 #b[사냥터]" },
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
		local sel = me:dialog_list(npc, "레벨대에 맞는 사냥터를 추천해드립니다.", choices)
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
