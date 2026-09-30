-- NPC name (String.wz/Npc.img.xml): 네네치킨

local spots = {
	{ 864000000, 30, 250, "#e퀘스트 의뢰와 대장간 #b[강화]" },
	{ 864000100, 30, 250, "차원의균열-아인크라드" },
	{ 130000100, 10, 12, "시그너스 전직소" },
	{ 130000100, 30, 31, "시그너스 전직소" },
	{ 130000100, 70, 71, "시그너스 전직소" },
	{ 200100001, 55, 250, "크리세" },
	{ 223000000, 120, 250, "테마파크" },
	{ 229000000, 10, 250, "헌티드 맨션" },
	{ 240090000, 150, 250, "암벽 거인" },
	{ 910811000, 30, 250, "전국시대-히에이잔" },
	{ 211060010, 130, 250, "사자왕의 성 - 차디찬 벌판" },
	{ 301000000, 250, 255, "크림슨우드의 성채" },
	{ 950100000, 10, 250, "황금 사원" },
	{ 231010000, 70, 250, "벚꽃성 외곽" },
	{ 123456788, 155, 250, "도로시의 공간" },
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
		local sel = me:dialog_list(npc, "특별한 곳에 데려다주는 택시다.", choices)
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
