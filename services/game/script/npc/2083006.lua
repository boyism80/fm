-- NPC name (String.wz/Npc.img.xml): 타임 게이트

local GATES = {
	{ text = "2021년 평범한 마을", map = 240070100 },
	{ text = "2099년 한밤의 항만", map = 240070200 },
	{ text = "2215년 폭격을 맞은 도심", map = 240070300 },
	{ text = "2216년 폐허가 된 도심", map = 240070400 },
	{ text = "2230년 위험한 타워", map = 240070500 },
	{ text = "2253년 천공전함 헤르메스호", map = 240070600 },
}

return {
	on_click = function(me, npc)
		local count = 0
		for _, it in pairs(me:item(4001393)) do
			count = count + it:count()
		end
		if count < 1 then
			me:message("시간여행자의 회중시계를 소지할 경우에만 이동할 수 있습니다.")
			return
		end

		local options = {}
		for i, gate in ipairs(GATES) do
			options[i] = gate.text
		end
		local sel = me:dialog_list(npc, "", options)
		if sel == nil then
			return
		end
		me:play_portal_sound()
		me:map(GATES[sel].map, 1)
	end
}
