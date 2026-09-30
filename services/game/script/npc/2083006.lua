-- NPC name (String.wz/Npc.img.xml): 타임 게이트

local GATES = {
	{ quest = 3719, text = "2021년 평범한 마을", map = 240070100 },
	{ quest = 3724, text = "2099년 한밤의 항만", map = 240070200 },
	{ quest = 3730, text = "2215년 폭격을 맞은 도심", map = 240070300 },
	{ quest = 3736, text = "2216년 폐허가 된 도심", map = 240070400 },
	{ quest = 3742, text = "2230년 위험한 타워", map = 240070500 },
	{ quest = 3748, text = "2253년 천공전함 헤르메스호", map = 240070600 },
}

return {
	on_click = function(me, npc)
		local options = {}
		local maps = {}
		for _, gate in ipairs(GATES) do
			local q = me:quest(gate.quest)
			if q ~= nil and q:completed() then
				table.insert(options, gate.text)
				table.insert(maps, gate.map)
			end
		end
		if #options == 0 then
			me:message("퀘스트를 진행할 경우에만 타임 게이트를 이용할 수 있습니다.")
			return
		end
		local count = 0
		for _, it in pairs(me:item(4001393)) do
			count = count + it:count()
		end
		if count < 1 then
			me:message("시간여행자의 회중시계를 소지할 경우에만 이동할 수 있습니다.")
			return
		end

		local sel = me:dialog_list(npc, "", options)
		if sel == nil then
			return
		end
		local map = maps[sel]
		if map == nil then
			me:dialog(npc, "무언가 버그가 발생했습니다.")
			return
		end
		me:play_portal_sound()
		me:map(map, 1)
	end
}
