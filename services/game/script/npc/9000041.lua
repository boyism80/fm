-- NPC name (String.wz/Npc.img.xml): 기부함

return {
	on_click = function(me, npc)
		local q = me:quest(29503)
		if not q:started() then
			me:dialog(npc, "자율적인 기금을 받는 자선 냄비다.")
			return
		end
		local try_map = tonumber(q:record_ex("trymap")) or 0
		if try_map ~= me:map():wz():id() then
			me:dialog(npc, "내가 기부하려는 마을은... #b#m" .. try_map .. "##k였던가?")
			return
		end
		if me:meso() < 100000 then
			me:dialog(npc, "최소 기부 금액은 10만 메소 입니다.")
			return
		end

		local input = me:dialog_input(npc, "얼마를 기부하시겠습니까?\r\n최소 기부금액은 100,000 메소입니다.\r\n기부 가능 최대금액은 2,147,483,647 메소입니다.\r\n#b< 현재 소지금액 : " .. me:meso() .. " 메소 >#k")
		if input == nil then
			return
		end
		local amount = tonumber(input)
		if amount == nil or amount > me:meso() or amount < 100000 or amount > 2147483647 then
			me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
			return
		end
		local money = tonumber(q:record_ex("money")) or 0
		if 2147483647 - money < amount then
			me:dialog(npc, "기부하려는 액수가 너무 큽니다. " .. (2147483647 - money) .. " 메소 이하로 입력해 주세요.")
			return
		end

		local code = me:exchange({ meso = amount }, nil)
		if code ~= ExchangeResult.OK then
			return
		end
		q:record_ex("money", tostring(money + amount))
		me:dialog(npc, amount .. " 메소를 기부하셨습니다. 마을 사람들은 당신의 선행을 잊지 않을 것입니다.")
	end
}
