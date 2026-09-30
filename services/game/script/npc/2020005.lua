-- NPC name (String.wz/Npc.img.xml): 알케스터

local GOODS = {
	{ item = 2050003, cost = 300, desc = "저주와 봉인 상태이상을 해제하는 아이템일세" },
	{ item = 2050004, cost = 400, desc = "모든 상태이상을 해제하는 아이템일세" },
	{ item = 4006000, cost = 3000, desc = "강력한 마법을 사용하는데 필요한 아이템일세" },
	{ item = 4006001, cost = 3000, desc = "강력한 소환 마법 아이템을 사용하는데 필요한 아이템일세" },
}
local ELIXIR = { item = 2000004, cost = 3800, max = 150, record = 110311, name = "엘릭서", desc = "HP와 MP를 최대 HP/MP의 50% 만큼 회복하는 아이템일세" }
local POWER_ELIXIR = { item = 2000005, cost = 6800, max = 40, record = 110312, name = "파워 엘릭서", desc = "HP와 MP를 모두 회복하는 아이템일세" }

return {
	on_click = function(me, npc)
		if not me:quest(3035):completed() then
			me:dialog(npc, "나는 이곳에서 30년 이상을 살아온 연금술사, 알케스터라고 한다네.", false, true)
			return
		end

		local d = datetime()
		local today = string.format("%04d%02d%02d", d.year, d.month, d.day)
		local day_q = me:quest(110310)
		local elixir_q = me:quest(ELIXIR.record)
		local power_q = me:quest(POWER_ELIXIR.record)
		local reset = not day_q:started() or day_q:record() ~= today
		if not day_q:started() then
			day_q:start(today)
		else
			day_q:record(today)
		end
		for _, limited in ipairs({ { q = elixir_q, max = ELIXIR.max }, { q = power_q, max = POWER_ELIXIR.max } }) do
			if not limited.q:started() then
				limited.q:start(tostring(limited.max))
			elseif reset then
				limited.q:record(tostring(limited.max))
			end
		end
		local elixir_left = tonumber(elixir_q:record()) or 0
		local power_left = tonumber(power_q:record()) or 0

		local options = {}
		for i, goods in ipairs(GOODS) do
			options[i] = "#b#t" .. goods.item .. "# (가격: " .. goods.cost .. " 메소)#k"
		end
		table.insert(options, "#b#t" .. ELIXIR.item .. "# (가격: " .. ELIXIR.cost .. " 메소) (오늘 남은수량 : " .. elixir_left .. ")#k")
		table.insert(options, "#b#t" .. POWER_ELIXIR.item .. "# (가격: " .. POWER_ELIXIR.cost .. " 메소) (오늘 남은수량 : " .. power_left .. ")#k")
		local sel = me:dialog_list(npc, "오오, 자네는 저번에 #b#t4031056##k을 찾아준 고마운 젊은이가 아닌가! 어떻게 찾아왔는가? 무엇인가 필요한 아이템이라도 있는가?#b\r\n", options)
		if sel == nil then
			return
		end
		local goods = GOODS[sel]
		local limited = nil
		local left = 0
		if sel == #GOODS + 1 then
			goods = ELIXIR
			limited = elixir_q
			left = elixir_left
		elseif sel == #GOODS + 2 then
			goods = POWER_ELIXIR
			limited = power_q
			left = power_left
		end

		local input = me:dialog_input(npc, "흐음.. #b#t" .. goods.item .. "##k 을 구매하고 싶단 말이지..? " .. goods.desc .. ". 몇개를 구매하고 싶은가? 가격은 #b" .. goods.cost .. " 메소#k일세.")
		if input == nil then
			return
		end
		local amount = tonumber(input)
		if amount == nil or amount < 1 or amount > 100 then
			me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
			return
		end
		local total = goods.cost * amount
		if not me:dialog_yes_no(npc, "흠.. #r" .. amount .. "개의 #t" .. goods.item .. "##k 아이템을 구매하고 싶다 이건가? 가격은 개당 " .. goods.cost .. " 메소일세. 총 가격은 #r" .. total .. " 메소#k라네. 정말 구매하고 싶은가?") then
			me:dialog(npc, "그런가? 언제든지 필요하면 다시 찾아오게나", false, true)
			return
		end
		if limited ~= nil and left < amount then
			me:dialog(npc, "흐음. 오늘 구매할 수 있는 양보다 적게 선택하게나. 오늘은 #b" .. left .. "#k 개의 " .. goods.name .. "를 더 구매할 수 있다네.", false, true)
			return
		end

		local code = me:exchange({ meso = total }, { item = { [goods.item] = amount } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "자네.. 분명 메소는 제대로 갖고 있는건가? 아니면 인벤토리 공간이 부족한건 아닌가?", false, true)
			return
		end
		if limited ~= nil then
			limited:record(tostring(left - amount))
		end
		me:dialog(npc, "자, 여기있네. 또 필요하면 언제든지 다시 찾아오게나.", false, true)
	end
}
