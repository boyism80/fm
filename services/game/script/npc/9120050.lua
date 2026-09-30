-- NPC name (String.wz/Npc.img.xml): 입실 제어 장치

local GATES = {
	{ quest = 50433, text = "검은 천사", map = 802000821 },
	{ quest = 50434, text = "결투후의 옥상", map = 802000823 },
	{ quest = 50400, text = "롯폰기몰", map = 802000800 },
	{ quest = 50434, text = "파괴의 천사", map = 802000824 },
	{ quest = 50436, text = "차원의 겹친 장소", map = 802000825 },
	{ quest = 50400, text = "카무나", map = 802000101 },
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
