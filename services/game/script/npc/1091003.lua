-- NPC name (String.wz/Npc.img.xml): 세릴

local CATEGORIES = {
	{
		prompt = "만들 너클을 골라봐.#b",
		recipes = {
			{ item = 1482001, label = "#t1482001# (레벨 제한: 15, 해적)", level = 15, cost = 1000, mats = { { 4000021, 20 } } },
			{ item = 1482002, label = "#t1482002# (레벨 제한: 20, 해적)", level = 20, cost = 2000, mats = { { 4011001, 1 }, { 4011000, 1 }, { 4000021, 10 }, { 4003000, 5 } } },
			{ item = 1482003, label = "#t1482003# (레벨 제한: 25, 해적)", level = 25, cost = 5000, mats = { { 4011000, 2 }, { 4011001, 1 }, { 4003000, 10 } } },
			{ item = 1482004, label = "#t1482004# (레벨 제한: 30, 해적)", level = 30, cost = 15000, mats = { { 4011000, 1 }, { 4011001, 1 }, { 4000021, 30 }, { 4003000, 10 } } },
			{ item = 1482005, label = "#t1482005# (레벨 제한: 35, 해적)", level = 35, cost = 30000, mats = { { 4011000, 2 }, { 4011001, 2 }, { 4000021, 30 }, { 4003000, 20 } } },
			{ item = 1482006, label = "#t1482006# (레벨 제한: 40, 해적)", level = 40, cost = 50000, mats = { { 4011000, 1 }, { 4011001, 1 }, { 4021000, 2 }, { 4000021, 50 }, { 4003000, 20 } } },
			{ item = 1482007, label = "#t1482007# (레벨 제한: 50, 해적)", level = 50, cost = 100000, mats = { { 4000039, 150 }, { 4011000, 1 }, { 4011001, 2 }, { 4000030, 20 }, { 4000021, 20 }, { 4003000, 20 } } },
		},
	},
	{
		prompt = "만들 건을 골라봐.#b",
		recipes = {
			{ item = 1492001, label = "#t1492001# (레벨 제한: 15, 해적)", level = 15, cost = 1000, mats = { { 4011000, 1 }, { 4003000, 5 }, { 4003001, 1 } } },
			{ item = 1492002, label = "#t1492002# (레벨 제한: 20, 해적)", level = 20, cost = 2000, mats = { { 4011000, 1 }, { 4003000, 10 }, { 4003001, 5 }, { 4000021, 10 } } },
			{ item = 1492003, label = "#t1492003# (레벨 제한: 25, 해적)", level = 25, cost = 5000, mats = { { 4011000, 2 }, { 4003000, 10 } } },
			{ item = 1492004, label = "#t1492004# (레벨 제한: 30, 해적)", level = 30, cost = 15000, mats = { { 4011001, 2 }, { 4000021, 10 }, { 4003000, 10 } } },
			{ item = 1492005, label = "#t1492005# (레벨 제한: 35, 해적)", level = 35, cost = 30000, mats = { { 4011006, 10 }, { 4011001, 2 }, { 4000021, 5 }, { 4003000, 10 } } },
			{ item = 1492006, label = "#t1492006# (레벨 제한: 40, 해적)", level = 40, cost = 50000, mats = { { 4011004, 1 }, { 4011001, 2 }, { 4000021, 10 }, { 4003000, 20 } } },
			{ item = 1492007, label = "#t1492007# (레벨 제한: 50, 해적)", level = 50, cost = 100000, mats = { { 4011006, 1 }, { 4011004, 2 }, { 4011001, 4 }, { 4000030, 30 }, { 4003000, 30 } } },
		},
	},
	{
		prompt = "만들 장갑을 골라봐.#b",
		recipes = {
			{ item = 1082180, label = "#t1082180#", level = 15, cost = 1000, mats = { { 4000021, 15 }, { 4021003, 1 } } },
			{ item = 1082183, label = "#t1082183#", level = 20, cost = 8000, mats = { { 4000021, 35 } } },
			{ item = 1082186, label = "#t1082186#", level = 25, cost = 15000, mats = { { 4011000, 2 }, { 4000021, 20 } } },
			{ item = 1082189, label = "#t1082189#", level = 30, cost = 25000, mats = { { 4021006, 2 }, { 4000021, 50 }, { 4003000, 10 } } },
			{ item = 1082192, label = "#t1082192#", level = 35, cost = 30000, mats = { { 4011000, 3 }, { 4000021, 60 }, { 4003000, 15 } } },
			{ item = 1082195, label = "#t1082195#", level = 40, cost = 40000, mats = { { 4000021, 80 }, { 4011000, 3 }, { 4011001, 3 }, { 4003000, 25 } } },
			{ item = 1082198, label = "#t1082198#", level = 50, cost = 50000, mats = { { 4011000, 3 }, { 4000021, 20 }, { 4000030, 40 }, { 4003000, 30 } } },
			{ item = 1082201, label = "#t1082201#", level = 60, cost = 70000, mats = { { 4011007, 1 }, { 4021008, 1 }, { 4021007, 1 }, { 4000030, 50 }, { 4003000, 50 } } },
		},
	},
}

return {
	on_click = function(me, npc)
		local type_sel = me:dialog_list(npc, "무기나 장갑을 만들고 싶다구? 그렇다면 제대로 찾아온 것 같군. 해적 세월 20년 동안 갈고 닦은 나의 장인 실력을 보여주지.#b", {
			" 너클 제작",
			" 건 제작",
			" 장갑 제작",
		})
		if type_sel == nil then
			return
		end

		local category = CATEGORIES[type_sel]
		local options = {}
		for i, recipe in ipairs(category.recipes) do
			options[i] = " " .. recipe.label
		end
		local item_sel = me:dialog_list(npc, category.prompt, options)
		if item_sel == nil then
			return
		end

		local recipe = category.recipes[item_sel]
		local prompt = "#t" .. recipe.item .. "# 제작에 필요한 아이템들이야. 레벨 제한은 " .. recipe.level .. " 이니까 착용 가능한지 확인해보라구.\r\n"
		for _, mat in ipairs(recipe.mats) do
			prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "#" .. mat[2] .. " 개 "
		end
		prompt = prompt .. "\r\n#i4031138# " .. recipe.cost .. " 메소"
		if not me:dialog_yes_no(npc, prompt) then
			return
		end

		if me:meso() < recipe.cost then
			me:dialog(npc, "메소는 제대로 갖고 있는거야?")
			return
		end
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			cost_items[mat[1]] = mat[2]
		end
		local code = me:exchange({ item = cost_items, meso = recipe.cost }, { item = { [recipe.item] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료가 부족한 것 같은데? 다시 확인해봐.", false, true)
			return
		end
		me:dialog(npc, "자, 다 됐어. 근사하지 않아?")
	end
}
