-- NPC name (String.wz/Npc.img.xml): 비셔스

local CATEGORIES = {
	{
		text = "활 제작",
		prompt = "난 저격수였지. 하지만 활과 석궁은 그다지 차이가 많이 나진 않아... 어쨌든, 뭘 만들고 싶지?#b",
		equip = true,
		recipes = {
			{ item = 1452002, label = " - 궁수 Lv. 10", mats = { { 4003001, 5 }, { 4000000, 30 } }, cost = 800 },
			{ item = 1452003, label = " - 궁수 Lv. 15", mats = { { 4011001, 1 }, { 4003000, 3 } }, cost = 2000 },
			{ item = 1452001, label = " - 궁수 Lv. 20", mats = { { 4003001, 30 }, { 4000016, 50 } }, cost = 3000 },
			{ item = 1452000, label = " - 궁수 Lv. 25", mats = { { 4011001, 2 }, { 4021006, 2 }, { 4003000, 8 } }, cost = 5000 },
			{ item = 1452005, label = " - 궁수 Lv. 30", mats = { { 4011001, 5 }, { 4011006, 5 }, { 4021003, 3 }, { 4021006, 3 }, { 4003000, 30 } }, cost = 30000 },
			{ item = 1452006, label = " - 궁수 Lv. 35", mats = { { 4011004, 7 }, { 4021000, 6 }, { 4021004, 3 }, { 4003000, 35 } }, cost = 40000 },
			{ item = 1452007, label = " - 궁수 Lv. 40", mats = { { 4021008, 1 }, { 4011001, 10 }, { 4011006, 3 }, { 4003000, 40 }, { 4000014, 50 } }, cost = 80000 },
		},
	},
	{
		text = "석궁 제작",
		prompt = "나는 한때 저격수였어. 석궁은 내 전문이지. 어떤 것을 만들어줄까?#b",
		equip = true,
		recipes = {
			{ item = 1462001, label = " - 궁수 Lv. 12", mats = { { 4003001, 7 }, { 4003000, 2 } }, cost = 1000 },
			{ item = 1462002, label = " - 궁수 Lv. 18", mats = { { 4011001, 1 }, { 4003001, 20 }, { 4003000, 5 } }, cost = 2000 },
			{ item = 1462003, label = " - 궁수 Lv. 22", mats = { { 4011001, 1 }, { 4003001, 50 }, { 4003000, 8 } }, cost = 3000 },
			{ item = 1462000, label = " - 궁수 Lv. 28", mats = { { 4011001, 2 }, { 4021006, 1 }, { 4021002, 1 }, { 4003000, 10 } }, cost = 10000 },
			{ item = 1462004, label = " - 궁수 Lv. 32", mats = { { 4011001, 5 }, { 4011005, 5 }, { 4021006, 3 }, { 4003001, 50 }, { 4003000, 15 } }, cost = 30000 },
			{ item = 1462005, label = " - 궁수 Lv. 38", mats = { { 4021008, 1 }, { 4011001, 8 }, { 4011006, 4 }, { 4021006, 2 }, { 4003000, 30 } }, cost = 50000 },
			{ item = 1462006, label = " - 궁수 Lv. 42", mats = { { 4021008, 2 }, { 4011004, 6 }, { 4003001, 30 }, { 4003000, 30 } }, cost = 80000 },
			{ item = 1462007, label = " - 궁수 Lv. 50", mats = { { 4021008, 2 }, { 4011006, 5 }, { 4021006, 3 }, { 4003001, 40 }, { 4003000, 40 } }, cost = 200000 },
		},
	},
	{
		text = "장갑 제작",
		prompt = "좋아. 어떤 장갑을 만들어줬으면 좋겠지?#b",
		equip = true,
		recipes = {
			{ item = 1082012, label = " - 궁수 Lv. 15", mats = { { 4000021, 15 }, { 4000009, 20 } }, cost = 5000 },
			{ item = 1082013, label = " - 궁수 Lv. 20", mats = { { 4000021, 20 }, { 4000009, 20 }, { 4011001, 2 } }, cost = 10000 },
			{ item = 1082016, label = " - 궁수 Lv. 25", mats = { { 4000021, 40 }, { 4000009, 50 }, { 4011006, 2 } }, cost = 15000 },
			{ item = 1082048, label = " - 궁수 Lv. 30", mats = { { 4000021, 50 }, { 4011006, 2 }, { 4021001, 1 } }, cost = 20000 },
			{ item = 1082068, label = " - 궁수 Lv. 35", mats = { { 4011000, 1 }, { 4011001, 3 }, { 4000021, 60 }, { 4003000, 15 } }, cost = 30000 },
			{ item = 1082071, label = " - 궁수 Lv. 40", mats = { { 4011001, 3 }, { 4021000, 1 }, { 4021002, 3 }, { 4000021, 80 }, { 4003000, 25 } }, cost = 40000 },
			{ item = 1082084, label = " - 궁수 Lv. 50", mats = { { 4011004, 3 }, { 4011006, 1 }, { 4021002, 2 }, { 4000030, 40 }, { 4003000, 35 } }, cost = 50000 },
			{ item = 1082089, label = " - 궁수 Lv. 60", mats = { { 4011006, 2 }, { 4011007, 1 }, { 4021006, 8 }, { 4000030, 50 }, { 4003000, 50 } }, cost = 70000 },
		},
	},
	{
		text = "장갑 합성",
		prompt = "좋아. 어떤 장갑을 만들어줬으면 좋겠지??#b",
		equip = true,
		recipes = {
			{ item = 1082015, label = " - 궁수 Lv. 20", mats = { { 1082013, 1 }, { 4021003, 2 } }, cost = 7000 },
			{ item = 1082014, label = " - 궁수 Lv. 20", mats = { { 1082013, 1 }, { 4021000, 1 } }, cost = 7000 },
			{ item = 1082017, label = " - 궁수 Lv. 25", mats = { { 1082016, 1 }, { 4021000, 3 } }, cost = 10000 },
			{ item = 1082018, label = " - 궁수 Lv. 25", mats = { { 1082016, 1 }, { 4021008, 1 } }, cost = 12000 },
			{ item = 1082049, label = " - 궁수 Lv. 30", mats = { { 1082048, 1 }, { 4021003, 3 } }, cost = 15000 },
			{ item = 1082050, label = " - 궁수 Lv. 30", mats = { { 1082048, 1 }, { 4021008, 1 } }, cost = 20000 },
			{ item = 1082069, label = " - 궁수 Lv. 35", mats = { { 1082068, 1 }, { 4011002, 4 } }, cost = 22000 },
			{ item = 1082070, label = " - 궁수 Lv. 35", mats = { { 1082068, 1 }, { 4011006, 2 } }, cost = 25000 },
			{ item = 1082072, label = " - 궁수 Lv. 40", mats = { { 1082071, 1 }, { 4011006, 4 } }, cost = 30000 },
			{ item = 1082073, label = " - 궁수 Lv. 40", mats = { { 1082071, 1 }, { 4021008, 2 } }, cost = 40000 },
			{ item = 1082085, label = " - 궁수 Lv. 50", mats = { { 1082084, 1 }, { 4011000, 1 }, { 4021000, 5 } }, cost = 55000 },
			{ item = 1082083, label = " - 궁수 Lv. 50", mats = { { 1082084, 1 }, { 4011006, 2 }, { 4021008, 2 } }, cost = 60000 },
			{ item = 1082090, label = " - 궁수 Lv. 60", mats = { { 1082089, 1 }, { 4021000, 5 }, { 4021007, 1 } }, cost = 70000 },
			{ item = 1082091, label = " - 궁수 Lv. 60", mats = { { 1082089, 1 }, { 4021007, 2 }, { 4021008, 2 } }, cost = 80000 },
		},
	},
	{
		text = "재료 제작",
		prompt = "재료? 몇가지 재료 제작 방법을 알고 있으니 만들어 줄 수 있겠어.#b",
		equip = false,
		recipes = {
			{ item = 4003001, text = "나뭇가지로 가공된 나무 제작", mats = { { 4000003, 10 } }, cost = 0 },
			{ item = 4003001, text = "장작으로 가공된 나무 제작", mats = { { 4000018, 5 } }, cost = 0 },
			{ item = 4003000, text = "나사 15개 제작", mats = { { 4011000, 1 }, { 4011001, 1 } }, cost = 0, gain = 15 },
		},
	},
	{
		text = "화살 제작",
		prompt = "화살? 문제 없지.#b",
		equip = true,
		recipes = {
			{ item = 2060000, text = "#t2060000#", mats = { { 4003001, 1 }, { 4003004, 1 } }, cost = 0, gain = 1000 },
			{ item = 2061000, text = "#t2061000#", mats = { { 4003001, 1 }, { 4003004, 1 } }, cost = 0, gain = 1000 },
			{ item = 2060001, text = "#t2060001#", mats = { { 4011000, 1 }, { 4003001, 3 }, { 4003004, 10 } }, cost = 0, gain = 900 },
			{ item = 2061001, text = "#t2061001#", mats = { { 4011000, 1 }, { 4003001, 3 }, { 4003004, 10 } }, cost = 0, gain = 900 },
			{ item = 2060002, text = "#t2060002#", mats = { { 4011001, 1 }, { 4003001, 5 }, { 4003005, 15 } }, cost = 0, gain = 800 },
			{ item = 2061002, text = "#t2061002#", mats = { { 4011001, 1 }, { 4003001, 5 }, { 4003005, 15 } }, cost = 0, gain = 800 },
		},
	},
}

return {
	on_click = function(me, npc)
		local category_options = {}
		for i, category in ipairs(CATEGORIES) do
			category_options[i] = category.text
		end
		local category_sel = me:dialog_list(npc, "흐음.. 궁수의 아이템이 필요한건가? 그렇다면 내가 도와줄 수 있는데.. 필요한 물건이라도 있어? 어떤 물건이 필요하지?#b", category_options)
		if category_sel == nil then
			return
		end
		local category = CATEGORIES[category_sel]

		local recipe_options = {}
		for i, recipe in ipairs(category.recipes) do
			if recipe.text ~= nil then
				recipe_options[i] = recipe.text
			else
				recipe_options[i] = "#z" .. recipe.item .. "##k" .. recipe.label
			end
		end
		local recipe_sel = me:dialog_list(npc, category.prompt, recipe_options)
		if recipe_sel == nil then
			return
		end
		local recipe = category.recipes[recipe_sel]

		local qty = 1
		if not category.equip then
			local input = me:dialog_input(npc, "만들고 싶은 아이템이 #b#t" .. recipe.item .. "##k 인가? 몇개를 만들고 싶어?")
			if input == nil then
				return
			end
			qty = tonumber(input)
			if qty == nil or qty < 1 or qty > 100 then
				me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
				return
			end
		end

		local prompt = "만들고 싶은 아이템이 "
		if qty == 1 then
			prompt = prompt .. "#t" .. recipe.item .. "#"
		else
			prompt = prompt .. " #t" .. recipe.item .. "#" .. qty .. "개"
		end
		prompt = prompt .. " 인가? 재료는 아래를 참조해.\r\n#b"
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] * qty .. "개 "
			cost_items[mat[1]] = mat[2] * qty
		end
		local meso = recipe.cost * qty
		if meso > 0 then
			prompt = prompt .. "\r\n#i4031138# " .. meso .. " 메소"
		end
		if not me:dialog_yes_no(npc, prompt) then
			return
		end

		if me:meso() < meso then
			me:dialog(npc, "흐음.. 내 서비스를 이용하려면 약간의 수수료가 필요하다네. 수수료를 충분히 갖고 있는지 확인해 보게나.")
			return
		end
		local cost = { item = cost_items }
		if meso > 0 then
			cost.meso = meso
		end
		local code = me:exchange(cost, { item = { [recipe.item] = (recipe.gain or 1) * qty } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료는 충분히 갖고 있는지, 인벤토리 공간이 부족한 건 아닌지 다시한번 확인해 보게나.")
			return
		end
		me:dialog(npc, "다 됐다네. 다른 원하는것이 있다면 언제든지 내게 말을 걸어주게.")
	end
}
