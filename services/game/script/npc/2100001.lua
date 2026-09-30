-- NPC name (String.wz/Npc.img.xml): 무하마드

local CATEGORIES = {
	{
		text = "광석 제련",
		prompt = "어떤 광석을 제련하고 싶은가?#b",
		recipes = {
			{ item = 4011000, text = "청동", mats = { { 4010000, 10 } }, cost = 270 },
			{ item = 4011001, text = "강철", mats = { { 4010001, 10 } }, cost = 270 },
			{ item = 4011002, text = "미스릴", mats = { { 4010002, 10 } }, cost = 270 },
			{ item = 4011003, text = "아다만티움", mats = { { 4010003, 10 } }, cost = 450 },
			{ item = 4011004, text = "은", mats = { { 4010004, 10 } }, cost = 450 },
			{ item = 4011005, text = "오리할콘", mats = { { 4010005, 10 } }, cost = 450 },
			{ item = 4011006, text = "금", mats = { { 4010006, 10 } }, cost = 720 },
			{ item = 4011008, text = "리튬", mats = { { 4010007, 10 } }, cost = 270 },
		},
	},
	{
		text = "보석 제련",
		prompt = "어떤 보석을 제련하고 싶은가?#b",
		recipes = {
			{ item = 4021000, text = "가넷", mats = { { 4020000, 10 } }, cost = 450 },
			{ item = 4021001, text = "자수정", mats = { { 4020001, 10 } }, cost = 450 },
			{ item = 4021002, text = "아쿠아마린", mats = { { 4020002, 10 } }, cost = 450 },
			{ item = 4021003, text = "에메랄드", mats = { { 4020003, 10 } }, cost = 450 },
			{ item = 4021004, text = "오팔", mats = { { 4020004, 10 } }, cost = 450 },
			{ item = 4021005, text = "사파이어", mats = { { 4020005, 10 } }, cost = 450 },
			{ item = 4021006, text = "토파즈", mats = { { 4020006, 10 } }, cost = 450 },
			{ item = 4021007, text = "다이아몬드", mats = { { 4020007, 10 } }, cost = 900 },
			{ item = 4021008, text = "흑수정", mats = { { 4020008, 10 } }, cost = 2700 },
		},
	},
	{
		text = "크리스탈 제련",
		prompt = "크리스탈? 그 희귀한 물건 말인가? 그것의 원석들을 모아온건가? 자네도 대단하구만. 자 어떤 크리스탈을 제련해줬으면 좋겠나? #b",
		recipes = {
			{ item = 4005000, text = "힘의 크리스탈", mats = { { 4004000, 10 } }, cost = 4500 },
			{ item = 4005001, text = "지혜의 크리스탈", mats = { { 4004001, 10 } }, cost = 4500 },
			{ item = 4005002, text = "민첩성의 크리스탈", mats = { { 4004002, 10 } }, cost = 4500 },
			{ item = 4005003, text = "행운의 크리스탈", mats = { { 4004003, 10 } }, cost = 4500 },
		},
	},
}

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "광석이나 보석 제련에 관심이 있는가? 특히 이곳의 특산물인 리튬은 나밖에 제련할 수 없는 광석이라네.") then
			me:dialog(npc, "흐음? 그런가? 마음이 바뀌면 다시 찾아오게나.", false, true)
			return
		end

		local category_options = {}
		for i, category in ipairs(CATEGORIES) do
			category_options[i] = category.text
		end
		local category_sel = me:dialog_list(npc, "좋아, 자네. 마음에 드는군. 어떤 광석이나 보석을 제련해줄까? #b", category_options)
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

		local prompt = "#t" .. recipe.item .. "#.. 이것을 제련하기 위해선 아래와 같은 재료가 필요하다네. 얼마나 만들고 싶은가?"
		for _, mat in ipairs(recipe.mats) do
			prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] .. "개 "
		end
		if recipe.cost > 0 then
			prompt = prompt .. "\r\n#i4031138# " .. recipe.cost .. " 메소"
		end
		local input = me:dialog_input(npc, prompt)
		if input == nil then
			return
		end
		local qty = tonumber(input)
		if qty == nil or qty < 1 or qty > 100 then
			me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
			return
		end

		local meso = recipe.cost * qty
		if me:meso() < meso then
			me:dialog(npc, "흐음.. 수수료가 부족한 것 같구만.")
			return
		end
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			cost_items[mat[1]] = mat[2] * qty
		end
		local code = me:exchange({ item = cost_items, meso = meso }, { item = { [recipe.item] = qty } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족한건 아닌지, 혹은 부족한 재료가 있는건 아닌지 확인해보게나.", false, true)
			return
		end
		me:dialog(npc, "자, 다 됐다네. 훌륭한 예술 작품이지. 다른 물건이 필요하다면 언제든지 다시 나를 찾게나.")
	end
}
