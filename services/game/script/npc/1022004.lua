-- NPC name (String.wz/Npc.img.xml): 스미스

local CATEGORIES = {
	{
		text = "장갑 제작",
		prompt = "어떤 장갑을 만들어 보고 싶어?#b",
		equip = true,
		recipes = {
			{ item = 1082003, text = "#z1082003##k - 전사 Lv. 10#b", mats = { { 4000021, 15 }, { 4011001, 1 } }, cost = 1000 },
			{ item = 1082000, text = "#z1082000##k - 전사 Lv. 15#b", mats = { { 4011001, 2 } }, cost = 2000 },
			{ item = 1082004, text = "#z1082004##k - 전사 Lv. 20#b", mats = { { 4000021, 40 }, { 4011000, 2 } }, cost = 5000 },
			{ item = 1082001, text = "#z1082001##k - 전사 Lv. 25#b", mats = { { 4011001, 2 } }, cost = 10000 },
			{ item = 1082007, text = "#z1082007##k - 전사 Lv. 30#b", mats = { { 4011000, 3 }, { 4011001, 2 }, { 4003000, 15 } }, cost = 20000 },
			{ item = 1082008, text = "#z1082008##k - 전사 Lv. 35#b", mats = { { 4000021, 30 }, { 4011001, 4 }, { 4003000, 15 } }, cost = 30000 },
			{ item = 1082023, text = "#z1082023##k - 전사 Lv. 40#b", mats = { { 4000021, 50 }, { 4011001, 5 }, { 4003000, 40 } }, cost = 40000 },
			{ item = 1082009, text = "#z1082009##k - 전사 Lv. 50#b", mats = { { 4011001, 3 }, { 4021007, 2 }, { 4000030, 30 }, { 4003000, 45 } }, cost = 50000 },
			{ item = 1082059, text = "#z1082059##k - 전사 Lv. 60#b", mats = { { 4011007, 1 }, { 4011000, 8 }, { 4011006, 2 }, { 4000030, 50 }, { 4003000, 50 } }, cost = 70000 },
		},
	},
	{
		text = "장갑 합성",
		prompt = "어떤 장갑을 합성해 보고 싶어?#b",
		equip = true,
		recipes = {
			{ item = 1082005, text = "#z1082005##k - 전사 Lv. 30#b", mats = { { 1082007, 1 }, { 4011001, 1 } }, cost = 20000 },
			{ item = 1082006, text = "#z1082006##k - 전사 Lv. 30#b", mats = { { 1082007, 1 }, { 4011005, 2 } }, cost = 25000 },
			{ item = 1082035, text = "#z1082035##k - 전사 Lv. 35#b", mats = { { 1082008, 1 }, { 4021006, 3 } }, cost = 30000 },
			{ item = 1082036, text = "#z1082036##k - 전사 Lv. 35#b", mats = { { 1082008, 1 }, { 4021008, 1 } }, cost = 40000 },
			{ item = 1082024, text = "#z1082024##k - 전사 Lv. 40#b", mats = { { 1082023, 1 }, { 4011003, 4 } }, cost = 45000 },
			{ item = 1082025, text = "#z1082025##k - 전사 Lv. 40#b", mats = { { 1082023, 1 }, { 4021008, 2 } }, cost = 50000 },
			{ item = 1082010, text = "#z1082010##k - 전사 Lv. 50#b", mats = { { 1082009, 1 }, { 4011002, 5 } }, cost = 55000 },
			{ item = 1082011, text = "#z1082011##k - 전사 Lv. 50#b", mats = { { 1082009, 1 }, { 4011006, 4 } }, cost = 60000 },
			{ item = 1082060, text = "#z1082060##k - 전사 Lv. 60#b", mats = { { 1082059, 1 }, { 4011002, 3 }, { 4021005, 5 } }, cost = 70000 },
			{ item = 1082061, text = "#z1082061##k - 전사 Lv. 60#b", mats = { { 1082059, 1 }, { 4021007, 2 }, { 4021008, 2 } }, cost = 80000 },
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
		local category_sel = me:dialog_list(npc, "난 선더님의 조수야. 뭘 해보고싶어?#b ", category_options)
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
