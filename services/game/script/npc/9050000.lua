-- NPC name (String.wz/Npc.img.xml): 소환수 피그미

local TOWNS = {
	{ name = "헤네시스", egg = 4170000 },
	{ name = "엘리니아", egg = 4170001 },
	{ name = "커닝시티", egg = 4170002 },
	{ name = "페리온", egg = 4170003 },
	{ name = "엘나스", egg = 4170004 },
	{ name = "루디브리엄", egg = 4170005 },
	{ name = "오르비스", egg = 4170006 },
	{ name = "아쿠아리움", egg = 4170007 },
	{ name = "노틸러스", egg = 4170009 },
}

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "피그미에게 어떤 행동을 하시겠습니까?\r\n#b", { "맛좋은 사료를 준다.", "피그미가 할 말이 있는 것 같습니다." })
		if sel == nil then
			return
		end
		if sel == 1 then
			me:dialog(npc, "귀여운 피그미가 이제 알을 안준답니다. 단풍샵 에서 구매해보세요.")
			return
		end

		local town_names = {}
		for i, town in ipairs(TOWNS) do
			town_names[i] = town.name
		end
		local from_sel = me:dialog_list(npc, "피그미가 자신이 낳은 알을 다른 지역의 알과 바꾸어 주겠다고 합니다. 어느 마을 피그미 에그를 교환하시겠습니까?\r\n", town_names)
		if from_sel == nil then
			return
		end
		local from = TOWNS[from_sel]
		local held = 0
		for _, it in pairs(me:item(from.egg)) do
			held = held + it:count()
		end
		if held == 0 then
			me:dialog(npc, "#b#h #님#k은 " .. from.name .. " 피그미 에그를 가지고 있지 않습니다.", false, true)
			return
		end

		local to_sel = me:dialog_list(npc, "#b" .. from.name .. " 피그미 에그#k를 어느 마을 피그미 에그로 교환하시겠습니까?\r\n", town_names)
		if to_sel == nil then
			return
		end
		if to_sel == from_sel then
			me:dialog(npc, "같은 마을 에그로는 교환 하실 수 없습니다.", false, true)
			return
		end
		local to = TOWNS[to_sel]

		local input = me:dialog_input(npc, "#b#h #님#k은 " .. from.name .. " 피그미 에그를 #b" .. held .. "개#k 가지고 있습니다. 몇 개를 교환하시 겠습니까?\r\n#b< 예 : 3 >#k")
		if input == nil then
			return
		end
		local qty = tonumber(input)
		if qty == nil or qty < 1 or qty > 32000 or qty > held then
			me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
			return
		end

		local code = me:exchange({ item = { [from.egg] = qty } }, { item = { [to.egg] = qty } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, held .. " 이하의 숫자만 가능합니다.", false, true)
			return
		end
		me:dialog(npc, "피그미 에그를 교환했습니다.", false, true)
	end
}
