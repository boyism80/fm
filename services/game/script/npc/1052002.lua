-- NPC name (String.wz/Npc.img.xml): 뒷골목의 제이엠

local CATEGORIES = {
	{
		text = "장갑 제작",
		prompt = "어떤 장갑을 만들어 보고 싶어?#b",
		equip = true,
		recipes = {
			{ item = 1082002, text = "#z1082002##k - 공용 Lv. 10#b", mats = { { 4000021, 15 } }, cost = 1000 },
			{ item = 1082029, text = "#z1082029##k - 도적 Lv. 15#b", mats = { { 4000021, 30 }, { 4000018, 20 } }, cost = 7000 },
			{ item = 1082030, text = "#z1082030##k - 도적 Lv. 15#b", mats = { { 4000021, 30 }, { 4000015, 20 } }, cost = 7000 },
			{ item = 1082031, text = "#z1082031##k - 도적 Lv. 15#b", mats = { { 4000021, 30 }, { 4000020, 20 } }, cost = 7000 },
			{ item = 1082032, text = "#z1082032##k - 도적 Lv. 20#b", mats = { { 4011000, 2 }, { 4000021, 40 } }, cost = 10000 },
			{ item = 1082037, text = "#z1082037##k - 도적 Lv. 25#b", mats = { { 4011000, 2 }, { 4011001, 1 }, { 4000021, 10 } }, cost = 15000 },
			{ item = 1082042, text = "#z1082042##k - 도적 Lv. 30#b", mats = { { 4011001, 2 }, { 4000021, 50 }, { 4003000, 10 } }, cost = 25000 },
			{ item = 1082046, text = "#z1082046##k - 도적 Lv. 35#b", mats = { { 4011001, 3 }, { 4011000, 1 }, { 4000021, 60 }, { 4003000, 15 } }, cost = 30000 },
			{ item = 1082075, text = "#z1082075##k - 도적 Lv. 40#b", mats = { { 4021000, 3 }, { 4000014, 200 }, { 4000021, 80 }, { 4003000, 30 } }, cost = 40000 },
			{ item = 1082065, text = "#z1082065##k - 도적 Lv. 50#b", mats = { { 4021005, 3 }, { 4021008, 1 }, { 4000030, 40 }, { 4003000, 30 } }, cost = 50000 },
			{ item = 1082092, text = "#z1082092##k - 도적 Lv. 60#b", mats = { { 4011007, 1 }, { 4011000, 8 }, { 4021007, 1 }, { 4000030, 50 }, { 4003000, 50 } }, cost = 70000 },
		},
	},
	{
		text = "장갑 합성",
		prompt = "어떤 장갑을 합성해 보고 싶어?#b",
		equip = true,
		recipes = {
			{ item = 1082033, text = "#z1082033##k - 도적 Lv. 20#b", mats = { { 1082032, 1 }, { 4011002, 1 } }, cost = 5000 },
			{ item = 1082034, text = "#z1082034##k - 도적 Lv. 20#b", mats = { { 1082032, 1 }, { 4021004, 1 } }, cost = 7000 },
			{ item = 1082038, text = "#z1082038##k - 도적 Lv. 25#b", mats = { { 1082037, 1 }, { 4011002, 2 } }, cost = 10000 },
			{ item = 1082039, text = "#z1082039##k - 도적 Lv. 25#b", mats = { { 1082037, 1 }, { 4021004, 2 } }, cost = 12000 },
			{ item = 1082043, text = "#z1082043##k - 도적 Lv. 30#b", mats = { { 1082042, 1 }, { 4011004, 2 } }, cost = 15000 },
			{ item = 1082044, text = "#z1082044##k - 도적 Lv. 30#b", mats = { { 1082042, 1 }, { 4011006, 1 } }, cost = 20000 },
			{ item = 1082047, text = "#z1082047##k - 도적 Lv. 35#b", mats = { { 1082046, 1 }, { 4011005, 3 } }, cost = 22000 },
			{ item = 1082045, text = "#z1082045##k - 도적 Lv. 35#b", mats = { { 1082046, 1 }, { 4011006, 2 } }, cost = 25000 },
			{ item = 1082076, text = "#z1082076##k - 도적 Lv. 40#b", mats = { { 1082075, 1 }, { 4011006, 4 } }, cost = 40000 },
			{ item = 1082074, text = "#z1082074##k - 도적 Lv. 40#b", mats = { { 1082075, 1 }, { 4021008, 2 } }, cost = 50000 },
			{ item = 1082067, text = "#z1082067##k - 도적 Lv. 50#b", mats = { { 1082065, 1 }, { 4021000, 5 } }, cost = 55000 },
			{ item = 1082066, text = "#z1082066##k - 도적 Lv. 50#b", mats = { { 1082065, 1 }, { 4011006, 2 }, { 4021008, 1 } }, cost = 60000 },
			{ item = 1082093, text = "#z1082093##k - 도적 Lv. 60#b", mats = { { 1082092, 1 }, { 4011001, 7 }, { 4000014, 200 } }, cost = 70000 },
			{ item = 1082094, text = "#z1082094##k - 도적 Lv. 60#b", mats = { { 1082092, 1 }, { 4011006, 7 }, { 4000027, 150 } }, cost = 80000 },
		},
	},
	{
		text = "아대 제작",
		prompt = "어떤 아대를 제작해 보고 싶어?#b",
		equip = true,
		recipes = {
			{ item = 1472001, text = "#z1472001##k - 도적 Lv. 15#b", mats = { { 4011001, 1 }, { 4000021, 20 }, { 4003000, 5 } }, cost = 2000 },
			{ item = 1472004, text = "#z1472004##k - 도적 Lv. 20#b", mats = { { 4011000, 2 }, { 4011001, 1 }, { 4000021, 30 }, { 4003000, 10 } }, cost = 3000 },
			{ item = 1472007, text = "#z1472007##k - 도적 Lv. 25#b", mats = { { 1472000, 1 }, { 4011001, 3 }, { 4000021, 20 }, { 4003001, 30 } }, cost = 5000 },
			{ item = 1472008, text = "#z1472008##k - 도적 Lv. 30#b", mats = { { 4011000, 3 }, { 4011001, 2 }, { 4000021, 50 }, { 4003000, 20 } }, cost = 15000 },
			{ item = 1472011, text = "#z1472011##k - 도적 Lv. 35#b", mats = { { 4011000, 4 }, { 4011001, 2 }, { 4000021, 80 }, { 4003000, 25 } }, cost = 30000 },
			{ item = 1472014, text = "#z1472014##k - 도적 Lv. 40#b", mats = { { 4011000, 3 }, { 4011001, 2 }, { 4000021, 100 }, { 4003000, 30 } }, cost = 40000 },
			{ item = 1472018, text = "#z1472018##k - 도적 Lv. 50#b", mats = { { 4011000, 4 }, { 4011001, 2 }, { 4000030, 40 }, { 4003000, 35 } }, cost = 50000 },
		},
	},
	{
		text = "아대 합성",
		prompt = "어떤 아대를 합성해 보고 싶어?#b",
		equip = true,
		recipes = {
			{ item = 1472002, text = "#z1472002##k - 도적 Lv. 15#b", mats = { { 1472001, 1 }, { 4011002, 1 } }, cost = 1000 },
			{ item = 1472003, text = "#z1472003##k - 도적 Lv. 15#b", mats = { { 1472001, 1 }, { 4011006, 1 } }, cost = 2000 },
			{ item = 1472005, text = "#z1472005##k - 도적 Lv. 20#b", mats = { { 1472004, 1 }, { 4011001, 2 } }, cost = 3000 },
			{ item = 1472006, text = "#z1472006##k - 도적 Lv. 20#b", mats = { { 1472004, 1 }, { 4011003, 2 } }, cost = 5000 },
			{ item = 1472009, text = "#z1472009##k - 도적 Lv. 30#b", mats = { { 1472008, 1 }, { 4011002, 3 } }, cost = 10000 },
			{ item = 1472010, text = "#z1472010##k - 도적 Lv. 30#b", mats = { { 1472008, 1 }, { 4011003, 3 } }, cost = 15000 },
			{ item = 1472012, text = "#z1472012##k - 도적 Lv. 35#b", mats = { { 1472011, 1 }, { 4011004, 4 } }, cost = 20000 },
			{ item = 1472013, text = "#z1472013##k - 도적 Lv. 35#b", mats = { { 1472011, 1 }, { 4021008, 1 } }, cost = 25000 },
			{ item = 1472015, text = "#z1472015##k - 도적 Lv. 40#b", mats = { { 1472014, 1 }, { 4021000, 5 } }, cost = 30000 },
			{ item = 1472016, text = "#z1472016##k - 도적 Lv. 40#b", mats = { { 1472014, 1 }, { 4011003, 5 } }, cost = 30000 },
			{ item = 1472017, text = "#z1472017##k - 도적 Lv. 40#b", mats = { { 1472014, 1 }, { 4021008, 2 } }, cost = 35000 },
			{ item = 1472019, text = "#z1472019##k - 도적 Lv. 50#b", mats = { { 1472018, 1 }, { 4021000, 6 } }, cost = 40000 },
			{ item = 1472020, text = "#z1472020##k - 도적 Lv. 50#b", mats = { { 1472018, 1 }, { 4021005, 6 } }, cost = 40000 },
		},
	},
	{
		text = "재료 제작",
		prompt = "재료? 그정도는 문제없어.#b",
		equip = false,
		recipes = {
			{ item = 4003001, text = "나뭇가지로 가공된 나무 제작", mats = { { 4000003, 10 } }, cost = 0 },
			{ item = 4003001, text = "장작으로 가공된 나무 제작", mats = { { 4000018, 5 } }, cost = 0 },
			{ item = 4003000, text = "나사 15개 제작", mats = { { 4011000, 1 }, { 4011001, 1 } }, cost = 0, gain = 15 },
		},
	},
}

return {
	on_click = function(me, npc)
		local category_options = {}
		for i, category in ipairs(CATEGORIES) do
			category_options[i] = category.text
		end
		local category_sel = me:dialog_list(npc, "쉿... 조용히 하라구.. 자.. 좋은물건들이 많이 있으니 천천히 골라보도록 하라구.#b ", category_options)
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
		local code = me:exchange(cost, { item = { [recipe.item] = (recipe.gain or 1) * qty } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료는 제대로 갖고 있는건가? 다시 한번 확인해보게. 만일 재료를 장비하고 있다면 장비를 해제하도록 하게. 혹은 인벤토리 공간이 부족한거 아닌가..? 제대로 다시 한번 확인해 보게.")
			return
		end
		me:dialog(npc, "자아.. 다 됐다구. 역시 완벽한 아이템이 탄생했잖아? 다른 작업도 필요하다면 다시 나에게 찾아오라구.")
	end
}
