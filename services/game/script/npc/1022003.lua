-- NPC name (String.wz/Npc.img.xml): 선더

local CATEGORIES = {
	{
		text = "광석 제련",
		prompt = "어떤 광석을 제련해 보고 싶나?#b",
		equip = false,
		recipes = {
			{ item = 4011000, text = "#t4011000#", mats = { { 4010000, 10 } }, cost = 300 },
			{ item = 4011001, text = "#t4011001#", mats = { { 4010001, 10 } }, cost = 300 },
			{ item = 4011002, text = "#t4011002#", mats = { { 4010002, 10 } }, cost = 300 },
			{ item = 4011003, text = "#t4011003#", mats = { { 4010003, 10 } }, cost = 500 },
			{ item = 4011004, text = "#t4011004#", mats = { { 4010004, 10 } }, cost = 500 },
			{ item = 4011005, text = "#t4011005#", mats = { { 4010005, 10 } }, cost = 500 },
			{ item = 4011006, text = "#t4011006#", mats = { { 4010006, 10 } }, cost = 800 },
		},
	},
	{
		text = "보석 제련",
		prompt = "어떤 보석을 제련해 보고 싶나?#b",
		equip = false,
		recipes = {
			{ item = 4021000, text = "#t4021000#", mats = { { 4020000, 10 } }, cost = 500 },
			{ item = 4021001, text = "#t4021001#", mats = { { 4020001, 10 } }, cost = 500 },
			{ item = 4021002, text = "#t4021002#", mats = { { 4020002, 10 } }, cost = 500 },
			{ item = 4021003, text = "#t4021003#", mats = { { 4020003, 10 } }, cost = 500 },
			{ item = 4021004, text = "#t4021004#", mats = { { 4020004, 10 } }, cost = 500 },
			{ item = 4021005, text = "#t4021005#", mats = { { 4020005, 10 } }, cost = 500 },
			{ item = 4021006, text = "#t4021006#", mats = { { 4020006, 10 } }, cost = 500 },
			{ item = 4021007, text = "#t4021007#", mats = { { 4020007, 10 } }, cost = 1000 },
			{ item = 4021008, text = "#t4021008#", mats = { { 4020008, 10 } }, cost = 3000 },
		},
	},
	{
		text = "헬멧 합성",
		prompt = "어떤 헬멧을 만들어 보고 싶나?#b",
		equip = true,
		recipes = {
			{ item = 1002042, text = "#z1002042##k - 공용 Lv. 15#b", mats = { { 1002001, 1 }, { 4011002, 1 } }, cost = 500 },
			{ item = 1002041, text = "#z1002041##k - 공용 Lv. 15#b", mats = { { 1002001, 1 }, { 4021006, 1 } }, cost = 300 },
			{ item = 1002002, text = "#z1002002##k - 전사 Lv. 10#b", mats = { { 1002043, 1 }, { 4011001, 1 } }, cost = 500 },
			{ item = 1002044, text = "#z1002044##k - 전사 Lv. 10#b", mats = { { 1002043, 1 }, { 4011002, 1 } }, cost = 800 },
			{ item = 1002003, text = "#z1002003##k - 전사 Lv. 12#b", mats = { { 1002039, 1 }, { 4011001, 1 } }, cost = 500 },
			{ item = 1002040, text = "#z1002040##k - 전사 Lv. 12#b", mats = { { 1002039, 1 }, { 4011002, 1 } }, cost = 800 },
			{ item = 1002007, text = "#z1002007##k - 전사 Lv. 15#b", mats = { { 1002051, 1 }, { 4011001, 2 } }, cost = 1000 },
			{ item = 1002052, text = "#z1002052##k - 전사 Lv. 15#b", mats = { { 1002051, 1 }, { 4011002, 2 } }, cost = 1500 },
			{ item = 1002011, text = "#z1002011##k - 전사 Lv. 20#b", mats = { { 1002059, 1 }, { 4011001, 3 } }, cost = 1500 },
			{ item = 1002058, text = "#z1002058##k - 전사 Lv. 20#b", mats = { { 1002059, 1 }, { 4011002, 3 } }, cost = 2000 },
			{ item = 1002009, text = "#z1002009##k - 전사 Lv. 20#b", mats = { { 1002055, 1 }, { 4011001, 3 } }, cost = 1500 },
			{ item = 1002056, text = "#z1002056##k - 전사 Lv. 20#b", mats = { { 1002055, 1 }, { 4011002, 3 } }, cost = 2000 },
			{ item = 1002087, text = "#z1002087##k - 전사 Lv. 22#b", mats = { { 1002027, 1 }, { 4011002, 4 } }, cost = 2000 },
			{ item = 1002088, text = "#z1002088##k - 전사 Lv. 22#b", mats = { { 1002027, 1 }, { 4011006, 4 } }, cost = 4000 },
			{ item = 1002049, text = "#z1002049##k - 전사 Lv. 25#b", mats = { { 1002005, 1 }, { 4011006, 5 } }, cost = 4000 },
			{ item = 1002050, text = "#z1002050##k - 전사 Lv. 25#b", mats = { { 1002005, 1 }, { 4011005, 5 } }, cost = 5000 },
			{ item = 1002047, text = "#z1002047##k - 전사 Lv. 35#b", mats = { { 1002004, 1 }, { 4021000, 3 } }, cost = 8000 },
			{ item = 1002048, text = "#z1002048##k - 전사 Lv. 35#b", mats = { { 1002004, 1 }, { 4021005, 3 } }, cost = 10000 },
			{ item = 1002099, text = "#z1002099##k - 전사 Lv. 40#b", mats = { { 1002021, 1 }, { 4011002, 5 } }, cost = 12000 },
			{ item = 1002098, text = "#z1002098##k - 전사 Lv. 40#b", mats = { { 1002021, 1 }, { 4011006, 6 } }, cost = 15000 },
			{ item = 1002085, text = "#z1002085##k - 전사 Lv. 50#b", mats = { { 1002086, 1 }, { 4011002, 5 } }, cost = 20000 },
			{ item = 1002028, text = "#z1002028##k - 전사 Lv. 50#b", mats = { { 1002086, 1 }, { 4011004, 4 } }, cost = 25000 },
			{ item = 1002022, text = "#z1002022##k - 전사 Lv. 55#b", mats = { { 1002100, 1 }, { 4011007, 1 }, { 4011001, 7 } }, cost = 30000 },
			{ item = 1002101, text = "#z1002101##k - 전사 Lv. 55#b", mats = { { 1002100, 1 }, { 4011007, 1 }, { 4011002, 7 } }, cost = 30000 },
		},
	},
	{
		text = "방패 합성",
		prompt = "어떤 방패를 합성해 보고 싶나?#b",
		equip = true,
		recipes = {
			{ item = 1092014, text = "#z1092014##k - 전사 Lv. 40#b", mats = { { 1092012, 1 }, { 4011003, 10 } }, cost = 100000 },
			{ item = 1092013, text = "#z1092013##k - 전사 Lv. 40#b", mats = { { 1092012, 1 }, { 4011002, 10 } }, cost = 100000 },
			{ item = 1092010, text = "#z1092010##k - 전사 Lv. 60#b", mats = { { 1092009, 1 }, { 4011007, 1 }, { 4011004, 15 } }, cost = 120000 },
			{ item = 1092011, text = "#z1092011##k - 전사 Lv. 60#b", mats = { { 1092009, 1 }, { 4011007, 1 }, { 4011003, 15 } }, cost = 120000 },
		},
	},
}

return {
	on_click = function(me, npc)
		local category_options = {}
		for i, category in ipairs(CATEGORIES) do
			category_options[i] = category.text
		end
		local category_sel = me:dialog_list(npc, "난 세상 최고의 대장장이, 선더라고 한다네.#b", category_options)
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
			local input = me:dialog_input(npc, "만들고 싶은 아이템이 #t" .. recipe.item .. "# 인가? 몇개를 만들고 싶나?\r\n")
			if input == nil then
				return
			end
			qty = tonumber(input)
			if qty == nil or qty < 1 or qty > 100 then
				me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
				return
			end
		end

		local prompt = "만들고 싶은 아이템이  #t" .. recipe.item .. "# " .. qty .. " 개 인가? 재료는 아래를 참조하게.\r\n#b"
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] * qty .. " 개"
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
			me:dialog(npc, "메소#k 는 제대로 갖고 있는건나? 다시 한번 확인해보게.")
			return
		end
		local cost = { item = cost_items }
		if meso > 0 then
			cost.meso = meso
		end
		local code = me:exchange(cost, { item = { [recipe.item] = qty } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료는 제대로 갖고 있는건가? 다시 한번 확인해보게. 만일 재료를 장비하고 있다면 장비를 해제하도록 하게. 혹은 인벤토리 공간이 부족한거 아닌가..? 제대로 다시 한번 확인해 보게.")
			return
		end
		me:dialog(npc, "자아.. 다 됐다구. 역시 완벽한 아이템이 탄생했잖아? 다른 작업도 필요하다면 다시 나에게 찾아오라구.")
	end
}
