-- NPC name (String.wz/Npc.img.xml): 쉐인

return {
	on_click = function(me, npc)
		local first = me:quest(2050)
		local second = me:quest(2051)

		local price = 0
		local target = 0
		if second:completed() then
			if not me:dialog(npc, "다시 왔구나? 요금은 받지 않을게.", false, true) then
				return
			end
			me:map(101000102)
			return
		elseif second:started() then
			price = me:level() * 200
			target = 101000102
			if not me:dialog_yes_no(npc, "다시 왔구나? 사비트라마의 부탁을 받고 온거야? 그냥 들어보내줄 순 없고, 약간의 입장료만 내면 들여보내줄게. 입장료는 " .. price .. " 메소야.") then
				me:dialog_yes_no(npc, "흠. 마음이 바뀌면 다시 이야기 해.")
				return
			end
		elseif first:completed() then
			if not me:dialog(npc, "다시 왔구나? 요금은 받지 않을게.", false, true) then
				return
			end
			me:map(101000100)
			return
		elseif first:started() then
			price = me:level() * 100
			target = 101000100
			if not me:dialog_yes_no(npc, "사비트라마의 부탁을 받고 온거야? 그냥 들어보내줄 순 없고, 약간의 입장료만 내면 들여보내줄게. 입장료는 " .. price .. " 메소야.") then
				me:dialog_yes_no(npc, "흠. 마음이 바뀌면 다시 이야기 해.")
				return
			end
		else
			me:dialog(npc, "이 안에는 진귀한 약초가 있는 모양이야. 하지만 아무나 들여보내줄 순 없지.")
			return
		end

		local code = me:exchange({ meso = price }, nil)
		if code ~= ExchangeResult.OK then
			return
		end
		me:map(target)
	end
}
