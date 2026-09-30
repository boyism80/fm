-- NPC name (String.wz/Npc.img.xml): 클라라

local spots = {
	{ 864000000, 10, 250, "#e퀘스트 의뢰와 대장간 #b[강화맵]" },
	{ 100000207, 30, 250, "#e루사 #b[퀘스트]" },
	{ 450002000, 150, 250, "#e츄츄 아일랜드 #b[퀘스트]" },
	{ 450003000, 160, 250, "#e레헬른 #b[퀘스트]" },
	{ 450005010, 170, 250, "#e아르카나 #b[퀘스트]" },
	{ 800040000, 30, 250, "#e영주성 #b[퀘스트]" },
	{ 130000000, 30, 250, "#e에레브 #b[퀘스트]" },
	{ 801000000, 70, 250, "#e쇼와 마을 #b[퀘스트]" },
	{ 101070000, 30, 250, "#e엘리넬 #b[퀘스트]" },
	{ 209000100, 10, 250, "#e크리스마스 #b[퀘스트]" },
	{ 540000000, 30, 250, "#e싱가포르 #b[퀘스트]" },
	{ 541000000, 60, 250, "#e항구마을 #b[퀘스트]" },
	{ 541020000, 90, 250, "#e버려진도시 #b[퀘스트]" },
	{ 231000000, 70, 250, "#e벚꽃성 #b[퀘스트]" },
	{ 223000000, 120, 250, "#e판타스틱 테마파크 #b[퀘스트]" },
	{ 300000000, 70, 250, "#e엘린숲 #b[퀘스트]" },
	{ 240090000, 150, 250, "#e암벽거인 #b[퀘스트]" },
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
		local sel = me:dialog_list(npc, "서버 중요맵에 이동시켜드립니다.", choices)
		if sel == nil then
			return
		end
		local map_id = maps[sel]
		if map_id == nil then
			return
		end
		if id2map(map_id) == nil then
			me:dialog(npc, "아직 갈 수 없는 곳입니다.")
			return
		end
		me:map(map_id, 0)
	end
}
