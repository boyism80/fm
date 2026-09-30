-- NPC name (String.wz/Npc.img.xml): 보겐

local CATEGORIES = {
	{
		text = "광석 제련",
		prompt = "어떤 종류의 광석을 제련하고 싶나?#b",
		equip = false,
		recipes = {
			{ item = 4011000, text = "청동", mats = { { 4010000, 10 } }, cost = 300 },
			{ item = 4011001, text = "강철", mats = { { 4010001, 10 } }, cost = 300 },
			{ item = 4011002, text = "미스릴", mats = { { 4010002, 10 } }, cost = 300 },
			{ item = 4011003, text = "아다만티움", mats = { { 4010003, 10 } }, cost = 500 },
			{ item = 4011004, text = "은", mats = { { 4010004, 10 } }, cost = 500 },
			{ item = 4011005, text = "오리할콘", mats = { { 4010005, 10 } }, cost = 500 },
			{ item = 4011006, text = "금", mats = { { 4010006, 10 } }, cost = 800 },
		},
	},
	{
		text = "보석 제련",
		prompt = "어떤 종류의 보석을 제련하고 싶나?#b",
		equip = false,
		recipes = {
			{ item = 4021000, text = "가넷", mats = { { 4020000, 10 } }, cost = 500 },
			{ item = 4021001, text = "자수정", mats = { { 4020001, 10 } }, cost = 500 },
			{ item = 4021002, text = "아쿠아마린", mats = { { 4020002, 10 } }, cost = 500 },
			{ item = 4021003, text = "에메랄드", mats = { { 4020003, 10 } }, cost = 500 },
			{ item = 4021004, text = "오팔", mats = { { 4020004, 10 } }, cost = 500 },
			{ item = 4021005, text = "사파이어", mats = { { 4020005, 10 } }, cost = 500 },
			{ item = 4021006, text = "토파즈", mats = { { 4020006, 10 } }, cost = 500 },
			{ item = 4021007, text = "다이아몬드", mats = { { 4020007, 10 } }, cost = 1000 },
			{ item = 4021008, text = "흑수정", mats = { { 4020008, 10 } }, cost = 3000 },
		},
	},
	{
		text = "희귀 보석 제련",
		prompt = "희귀 보석이라.. 어떤걸 생각하고 있나?#b",
		equip = false,
		recipes = {
			{ item = 4011007, text = "달의 돌", mats = { { 4011000, 1 }, { 4011001, 1 }, { 4011002, 1 }, { 4011003, 1 }, { 4011004, 1 }, { 4011005, 1 }, { 4011006, 1 } }, cost = 10000 },
			{ item = 4021009, text = "별의 돌", mats = { { 4021000, 1 }, { 4021001, 1 }, { 4021002, 1 }, { 4021003, 1 }, { 4021004, 1 }, { 4021005, 1 }, { 4021006, 1 }, { 4021007, 1 }, { 4021008, 1 } }, cost = 15000 },
		},
	},
	{
		text = "크리스탈 제련",
		prompt = "크리스탈 제련..? 흐음.. 구하기 어려웠을텐데.. 무슨 크리스탈 원석을 가져왔는가?#b",
		equip = false,
		recipes = {
			{ item = 4005000, text = "힘의 크리스탈", mats = { { 4004000, 10 } }, cost = 5000 },
			{ item = 4005001, text = "지혜의 크리스탈", mats = { { 4004001, 10 } }, cost = 5000 },
			{ item = 4005002, text = "민첩의 크리스탈", mats = { { 4004002, 10 } }, cost = 5000 },
			{ item = 4005003, text = "행운의 크리스탈", mats = { { 4004003, 10 } }, cost = 5000 },
			{ item = 4005004, text = "어둠의 크리스탈", mats = { { 4004004, 10 } }, cost = 100000 },
		},
	},
	{
		text = "재료 제작",
		prompt = "자네에게 몇가지 재료를 만들어 줄 수 있다네.#b",
		equip = false,
		recipes = {
			{ item = 4003001, text = "나뭇가지로 가공된 나무 제작", mats = { { 4000003, 10 } }, cost = 0 },
			{ item = 4003001, text = "장작으로 가공된 나무 제작", mats = { { 4000018, 5 } }, cost = 0 },
			{ item = 4003000, text = "나사 15개 제작", mats = { { 4011000, 1 }, { 4011001, 1 } }, cost = 0, gain = 15 },
		},
	},
	{
		text = "화살 제작",
		prompt = "화살? 흐음. 이 몸을 뭘로 보고. 당연히 문제 없다네.#b",
		equip = true,
		recipes = {
			{ item = 2060000, text = "활 전용 화살", mats = { { 4003001, 1 }, { 4003004, 1 } }, cost = 0, gain = 1000 },
			{ item = 2061000, text = "석궁 전용 화살", mats = { { 4003001, 1 }, { 4003004, 1 } }, cost = 0, gain = 1000 },
			{ item = 2060001, text = "활 전용 청동화살", mats = { { 4011000, 1 }, { 4003001, 3 }, { 4003004, 10 } }, cost = 0, gain = 900 },
			{ item = 2061001, text = "석궁 전용 청동화살", mats = { { 4011000, 1 }, { 4003001, 3 }, { 4003004, 10 } }, cost = 0, gain = 900 },
			{ item = 2060002, text = "활 전용 강철화살", mats = { { 4011001, 1 }, { 4003001, 5 }, { 4003005, 15 } }, cost = 0, gain = 800 },
			{ item = 2061002, text = "석궁 전용 강철화살", mats = { { 4011001, 1 }, { 4003001, 5 }, { 4003005, 15 } }, cost = 0, gain = 800 },
		},
	},
}

return {
	on_click = function(me, npc)
		local category_options = {}
		for i, category in ipairs(CATEGORIES) do
			category_options[i] = category.text
		end
		local category_sel = me:dialog_list(npc, "흠? 무얼 하러 온겐가? 아, 내 뛰어난 제련 솜씨를 듣고 찾아온거로구만? 그렇다면 실력발휘를 한번 해보도록 하지. 자. 원하는게 뭔가?#b", category_options)
		if category_sel == nil then
			return
		end
		local category = CATEGORIES[category_sel]

		local recipe_options = {}
		for i, recipe in ipairs(category.recipes) do
			recipe_options[i] = recipe.text
		end
		local recipe_sel = me:dialog_list(npc, category.prompt, recipe_options)
		if recipe_sel == nil then
			return
		end
		local recipe = category.recipes[recipe_sel]

		local qty = 1
		if not category.equip then
			local input = me:dialog_input(npc, "#t" .. recipe.item .. "# 아이템을 만들어 보고 싶은건가? 몇개를 만들어 보고 싶은가?")
			if input == nil then
				return
			end
			qty = tonumber(input)
			if qty == nil or qty < 1 or qty > 100 then
				me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
				return
			end
		end

		local prompt = "그렇다면 "
		if qty == 1 then
			prompt = prompt .. "#t" .. recipe.item .. "# 아이템을 만들고 싶다는 것이로군?"
		else
			prompt = prompt .. "#t" .. recipe.item .. "# " .. qty .. "개를 만들고 싶다는 것이로군?"
		end
		prompt = prompt .. " 그렇다면 다음과 같은 재료를 구해와야 하네. 그리고 인벤토리 공간도 충분한지 확인해 보게나.\r\n#b"
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] * qty .. "개"
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
