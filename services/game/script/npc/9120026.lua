-- NPC name (String.wz/Npc.img.xml): 크리스탈

local GATES = {
	{ quest = 50392, text = "2100년 오다이바", map = 802000200 },
	{ quest = 50398, text = "2095년 공원", map = 802000300 },
	{ quest = 50402, text = "2102년 아키하바라", map = 802000400 },
	{ quest = 50411, text = "2102년 플래그쉽", map = 802000600 },
	{ quest = 50418, text = "2102년 시부야", map = 802000700 },
	{ quest = 50430, text = "2102년 롯폰기몰", map = 802000800 },
}

return {
	on_click = function(me, npc)
		local options = {}
		local maps = {}
		for _, gate in ipairs(GATES) do
			if me:quest(gate.quest):completed() then
				table.insert(options, gate.text)
				table.insert(maps, gate.map)
			end
		end
		if #options == 0 then
			me:notice("퀘스트를 진행할 경우에만 타임 게이트를 이용할 수 있습니다.")
			return
		end
		local count = 0
		for _, it in pairs(me:item(4001393)) do
			count = count + it:count()
		end
		if count < 1 then
			me:notice("시간여행자의 회중시계를 소지할 경우에만 이동할 수 있습니다.")
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
