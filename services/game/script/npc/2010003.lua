-- NPC name (String.wz/Npc.img.xml): 네브

local CATEGORIES = {
	{
		prompt = "전사의 장갑이 필요하다구? 문제없지. 원하는게 뭐야?#b",
		recipes = {
			{ item = 1082103, label = "전사 레벨. 70", cost = 90000, mats = { { 4005000, 2 }, { 4011000, 8 }, { 4011006, 3 }, { 4000030, 70 }, { 4003000, 55 } } },
			{ item = 1082104, label = "전사 레벨. 70", cost = 90000, mats = { { 1082103, 1 }, { 4011002, 6 }, { 4021006, 4 } } },
			{ item = 1082105, label = "전사 레벨. 70", cost = 100000, mats = { { 1082103, 1 }, { 4021006, 8 }, { 4021008, 3 } } },
			{ item = 1082114, label = "전사 레벨. 80", cost = 100000, mats = { { 4005000, 2 }, { 4005002, 1 }, { 4021005, 8 }, { 4000030, 90 }, { 4003000, 60 } } },
			{ item = 1082115, label = "전사 레벨. 80", cost = 110000, mats = { { 1082114, 1 }, { 4005000, 1 }, { 4005002, 1 }, { 4021003, 7 } } },
			{ item = 1082116, label = "전사 레벨. 80", cost = 110000, mats = { { 1082114, 1 }, { 4005002, 3 }, { 4021000, 8 } } },
			{ item = 1082117, label = "전사 레벨. 80", cost = 120000, mats = { { 1082114, 1 }, { 4005000, 2 }, { 4005002, 1 }, { 4021008, 4 } } },
		},
	},
	{
		prompt = "궁수의 장갑이 필요하다구? 문제없지. 원하는게 뭐야?#b",
		recipes = {
			{ item = 1082106, label = "궁수 레벨. 70", cost = 90000, mats = { { 4005002, 2 }, { 4021005, 8 }, { 4011004, 3 }, { 4000030, 70 }, { 4003000, 55 } } },
			{ item = 1082107, label = "궁수 레벨. 70", cost = 90000, mats = { { 1082106, 1 }, { 4021006, 5 }, { 4011006, 3 } } },
			{ item = 1082108, label = "궁수 레벨. 70", cost = 100000, mats = { { 1082106, 1 }, { 4021007, 2 }, { 4021008, 3 } } },
			{ item = 1082109, label = "궁수 레벨. 80", cost = 100000, mats = { { 4005002, 2 }, { 4005000, 1 }, { 4021000, 8 }, { 4000030, 90 }, { 4003000, 60 } } },
			{ item = 1082110, label = "궁수 레벨. 80", cost = 110000, mats = { { 1082109, 1 }, { 4005002, 1 }, { 4005000, 1 }, { 4021005, 7 } } },
			{ item = 1082111, label = "궁수 레벨. 80", cost = 110000, mats = { { 1082109, 1 }, { 4005002, 1 }, { 4005000, 1 }, { 4021003, 7 } } },
			{ item = 1082112, label = "궁수 레벨. 80", cost = 120000, mats = { { 1082109, 1 }, { 4005002, 2 }, { 4005000, 1 }, { 4021008, 4 } } },
		},
	},
	{
		prompt = "마법사의 장갑이 필요하다구? 문제 없지. 원하는게 뭐야?#b",
		recipes = {
			{ item = 1082098, label = "마법사 레벨. 70", cost = 90000, mats = { { 4005001, 2 }, { 4011000, 6 }, { 4011004, 6 }, { 4000030, 70 }, { 4003000, 55 } } },
			{ item = 1082099, label = "마법사 레벨. 70", cost = 90000, mats = { { 1082098, 1 }, { 4021002, 6 }, { 4021007, 2 } } },
			{ item = 1082100, label = "마법사 레벨. 70", cost = 100000, mats = { { 1082098, 1 }, { 4021008, 3 }, { 4011006, 3 } } },
			{ item = 1082121, label = "마법사 레벨. 80", cost = 100000, mats = { { 4005001, 2 }, { 4005003, 1 }, { 4021003, 8 }, { 4000030, 90 }, { 4003000, 60 } } },
			{ item = 1082122, label = "마법사 레벨. 80", cost = 110000, mats = { { 1082121, 1 }, { 4005001, 1 }, { 4005003, 1 }, { 4021005, 7 } } },
			{ item = 1082123, label = "마법사 레벨. 80", cost = 120000, mats = { { 1082121, 1 }, { 4005001, 2 }, { 4005003, 1 }, { 4021008, 4 } } },
		},
	},
	{
		prompt = "돚거의 장갑이 필요하다구? 문제 없지. 원하는게 뭐야?#b",
		recipes = {
			{ item = 1082095, label = "도적 레벨. 70", cost = 90000, mats = { { 4005003, 2 }, { 4011000, 6 }, { 4011003, 6 }, { 4000030, 70 }, { 4003000, 55 } } },
			{ item = 1082096, label = "도적 레벨. 70", cost = 90000, mats = { { 1082095, 1 }, { 4011004, 6 }, { 4021007, 2 } } },
			{ item = 1082097, label = "도적 레벨. 70", cost = 100000, mats = { { 1082095, 1 }, { 4021007, 3 }, { 4011006, 3 } } },
			{ item = 1082118, label = "도적 레벨. 80", cost = 100000, mats = { { 4005003, 2 }, { 4005002, 1 }, { 4011002, 8 }, { 4000030, 90 }, { 4003000, 60 } } },
			{ item = 1082119, label = "도적 레벨. 80", cost = 110000, mats = { { 1082118, 1 }, { 4005003, 1 }, { 4005002, 1 }, { 4021001, 7 } } },
			{ item = 1082120, label = "도적 레벨. 80", cost = 120000, mats = { { 1082118, 1 }, { 4005003, 2 }, { 4005002, 1 }, { 4021000, 8 } } },
		},
	},
}

return {
	on_click = function(me, npc)
		local type_sel = me:dialog_list(npc, "오르비스에서 최고의 장갑 제작자인 나에게 원하는거라도 있어?#b", {
			"  전사 장갑 제작/합성",
			"  궁수 장갑 제작/합성",
			"  마법사 장갑 제작/합성",
			"  돚거 장갑 제작/합성",
		})
		if type_sel == nil then
			return
		end

		local category = CATEGORIES[type_sel]
		local options = {}
		for i, recipe in ipairs(category.recipes) do
			options[i] = " #z" .. recipe.item .. "##k - " .. recipe.label .. "#b"
		end
		local item_sel = me:dialog_list(npc, category.prompt, options)
		if item_sel == nil then
			return
		end

		local recipe = category.recipes[item_sel]
		local prompt = "#t" .. recipe.item .. "# 아이템을 만들어 보고 싶은거야? 필요한 재료는 아래와 같으니 재료를 제대로 갖고 있는지, 인벤토리 공간은 충분한지 확인해 봐.#b"
		for _, mat in ipairs(recipe.mats) do
			prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] .. "개"
		end
		prompt = prompt .. "\r\n#i4031138# " .. recipe.cost .. " 메소"
		if not me:dialog_yes_no(npc, prompt) then
			return
		end

		if me:meso() < recipe.cost then
			me:dialog(npc, "메소가 부족한것 같은데? 미안하지만 요금을 제대로 지불해야 서비스를 이용할 수 있다구.")
			return
		end
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			cost_items[mat[1]] = mat[2]
		end
		local code = me:exchange({ item = cost_items, meso = recipe.cost }, { item = { [recipe.item] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흐음. 재료는 제대로 갖고 있는지, 인벤토리 공간은 충분한지 다시 한번 확인해 봐.")
			return
		end
		me:dialog(npc, "다 됐다구. 더 이상 필요한게 있으면 다시 내게 말을 걸어줘.")
	end
}
