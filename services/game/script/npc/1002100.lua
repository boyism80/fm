-- NPC name (String.wz/Npc.img.xml): 제인

local GOODS = {
	{ item = 2000002, cost = 310, recover = 300, stat = "HP" },
	{ item = 2022003, cost = 1060, recover = 1000, stat = "HP" },
	{ item = 2022000, cost = 1600, recover = 800, stat = "MP" },
	{ item = 2001000, cost = 3120, recover = 1000, stat = "HP 와 MP" },
}

return {
	on_click = function(me, npc)
		if not me:quest(2013):completed() then
			if me:quest(2010):completed() then
				me:dialog(npc, "아직 제 물약을 구매하시기엔 충분히 강하지 않으신 것 같아요 ...", false, true)
			else
				me:dialog(npc, "제 꿈은 이곳 저곳 여행을 다니는 거랍니다. 여행자님처럼 말이죠... 하지만 저희 아버진 위험하다며 계속 허락해 주시지를 않아요...")
			end
			return
		end
		if not me:dialog(npc, "당신이군요. 당신 덕분에 많은 것들을 해낼 수 있었어요. 여행자님을 위해서 물약을 만들어 드리겠어요.", false, true) then
			return
		end

		local options = {}
		for i, goods in ipairs(GOODS) do
			options[i] = "#z" .. goods.item .. "# (가격 : " .. goods.cost .. " 메소)"
		end
		local sel = me:dialog_list(npc, "어떤 아이템을 구매하시고 싶으세요? #b", options)
		if sel == nil then
			return
		end
		local goods = GOODS[sel]

		local input = me:dialog_input(npc, "#b#t" .. goods.item .. "##k 아이템을 구매하고 싶으세요? #t" .. goods.item .. "# 아이템은 " .. goods.recover .. " " .. goods.stat .. "를 회복시켜 줍니다. 몇개를 구매하시고 싶으세요?")
		if input == nil then
			return
		end
		local amount = tonumber(input)
		if amount == nil or amount < 1 or amount > 100 then
			me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
			return
		end
		if not me:dialog_yes_no(npc, "#b#t" .. goods.item .. "##k 아이템을 #r" .. amount .. "#k 개 구매하시고 싶으세요? 개당 가격은 " .. goods.cost .. " 메소 이며, 총 가격은 " .. goods.cost * amount .. " 메소 입니다.") then
			me:dialog(npc, "아직 재료는 많이 남아 있답니다. 천천히 생각해보시고 다시 말을 걸어주세요.", false, true)
			return
		end

		local code = me:exchange({ meso = goods.cost * amount }, { item = { [goods.item] = amount } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "메소가 부족하신건 아닌가요? 혹은 인벤토리 공간이 충분한지 확인해주세요.", false, true)
			return
		end
		me:dialog(npc, "와주셔서 고마워요. 더 도와드릴 일이 있다면 언제든지 다시 찾아와주세요.", false, true)
	end
}
