-- NPC name (String.wz/Npc.img.xml): 에뜨랑

local CATEGORIES = {
	{
		prompt = "어떤 장갑을 만들어 보고 싶나?#b",
		recipes = {
			{ item = 1082019, label = "마법사 Lv. 15", mats = { { 4000021, 15 } }, cost = 7000 },
			{ item = 1082020, label = "마법사 Lv. 20", mats = { { 4000021, 30 }, { 4011001, 1 } }, cost = 15000 },
			{ item = 1082026, label = "마법사 Lv. 25", mats = { { 4000021, 50 }, { 4011006, 2 } }, cost = 20000 },
			{ item = 1082051, label = "마법사 Lv. 30", mats = { { 4000021, 60 }, { 4021006, 1 }, { 4021000, 2 } }, cost = 25000 },
			{ item = 1082054, label = "마법사 Lv. 35", mats = { { 4000021, 70 }, { 4011006, 1 }, { 4011001, 3 }, { 4021000, 2 } }, cost = 30000 },
			{ item = 1082062, label = "마법사 Lv. 40", mats = { { 4000021, 80 }, { 4021000, 3 }, { 4021006, 3 }, { 4003000, 30 } }, cost = 40000 },
			{ item = 1082081, label = "마법사 Lv. 50", mats = { { 4021000, 3 }, { 4011006, 2 }, { 4000030, 35 }, { 4003000, 40 } }, cost = 50000 },
			{ item = 1082086, label = "마법사 Lv. 60", mats = { { 4011007, 1 }, { 4011001, 8 }, { 4021007, 1 }, { 4000030, 50 }, { 4003000, 50 } }, cost = 70000 },
		},
	},
	{
		prompt = "어떤 장갑을 합성해 보고 싶나?#b",
		recipes = {
			{ item = 1082021, label = "마법사 Lv. 20", mats = { { 1082020, 1 }, { 4011001, 1 } }, cost = 20000 },
			{ item = 1082022, label = "마법사 Lv. 20", mats = { { 1082020, 1 }, { 4021001, 2 } }, cost = 25000 },
			{ item = 1082027, label = "마법사 Lv. 25", mats = { { 1082026, 1 }, { 4021000, 3 } }, cost = 30000 },
			{ item = 1082028, label = "마법사 Lv. 25", mats = { { 1082026, 1 }, { 4021008, 1 } }, cost = 40000 },
			{ item = 1082052, label = "마법사 Lv. 30", mats = { { 1082051, 1 }, { 4021005, 3 } }, cost = 35000 },
			{ item = 1082053, label = "마법사 Lv. 30", mats = { { 1082051, 1 }, { 4021008, 1 } }, cost = 40000 },
			{ item = 1082055, label = "마법사 Lv. 35", mats = { { 1082054, 1 }, { 4021005, 3 } }, cost = 40000 },
			{ item = 1082056, label = "마법사 Lv. 35", mats = { { 1082054, 1 }, { 4021008, 1 } }, cost = 45000 },
			{ item = 1082063, label = "마법사 Lv. 40", mats = { { 1082062, 1 }, { 4021002, 4 } }, cost = 45000 },
			{ item = 1082064, label = "마법사 Lv. 40", mats = { { 1082062, 1 }, { 4021008, 2 } }, cost = 50000 },
			{ item = 1082082, label = "마법사 Lv. 50", mats = { { 1082081, 1 }, { 4021002, 5 } }, cost = 55000 },
			{ item = 1082080, label = "마법사 Lv. 50", mats = { { 1082081, 1 }, { 4021008, 3 } }, cost = 60000 },
			{ item = 1082087, label = "마법사 Lv. 60", mats = { { 1082086, 1 }, { 4011004, 3 }, { 4011006, 5 } }, cost = 70000 },
			{ item = 1082088, label = "마법사 Lv. 60", mats = { { 1082086, 1 }, { 4021008, 2 }, { 4011006, 3 } }, cost = 80000 },
		},
	},
	{
		prompt = "어떤 모자를 합성해 보고 싶나?#b",
		recipes = {
			{ item = 1002065, label = "마법사 Lv. 30", mats = { { 1002064, 1 }, { 4011001, 3 } }, cost = 40000 },
			{ item = 1002013, label = "마법사 Lv. 30", mats = { { 1002064, 1 }, { 4011006, 3 } }, cost = 50000 },
		},
	},
	{
		prompt = "어떤 완드를 만들어 보고 싶나?#b",
		recipes = {
			{ item = 1372005, label = "공용 Lv. 8", mats = { { 4003001, 5 } }, cost = 1000 },
			{ item = 1372006, label = "공용 Lv. 13", mats = { { 4003001, 10 }, { 4000001, 50 } }, cost = 3000 },
			{ item = 1372002, label = "공용 Lv. 18", mats = { { 4011001, 1 }, { 4000009, 30 }, { 4003000, 5 } }, cost = 5000 },
			{ item = 1372004, label = "마법사 Lv. 23", mats = { { 4011002, 2 }, { 4003002, 1 }, { 4003000, 10 } }, cost = 12000 },
			{ item = 1372003, label = "마법사 Lv. 28", mats = { { 4011002, 3 }, { 4021002, 1 }, { 4003000, 10 } }, cost = 30000 },
			{ item = 1372001, label = "마법사 Lv. 33", mats = { { 4021006, 5 }, { 4011002, 3 }, { 4011001, 1 }, { 4003000, 15 } }, cost = 60000 },
			{ item = 1372000, label = "마법사 Lv. 38", mats = { { 4021006, 5 }, { 4021005, 5 }, { 4021007, 1 }, { 4003003, 1 }, { 4003000, 20 } }, cost = 120000 },
			{ item = 1372007, label = "마법사 Lv. 48", mats = { { 4011006, 4 }, { 4021003, 3 }, { 4021007, 2 }, { 4021002, 1 }, { 4003002, 1 }, { 4003000, 30 } }, cost = 200000 },
		},
	},
	{
		prompt = "어떤 스태프를 만들어 보고 싶나?#b",
		recipes = {
			{ item = 1382000, label = "마법사 Lv. 10", mats = { { 4003001, 5 } }, cost = 2000 },
			{ item = 1382003, label = "마법사 Lv. 15", mats = { { 4021005, 1 }, { 4011001, 1 }, { 4003000, 5 } }, cost = 2000 },
			{ item = 1382005, label = "마법사 Lv. 15", mats = { { 4021003, 1 }, { 4011001, 1 }, { 4003000, 5 } }, cost = 2000 },
			{ item = 1382004, label = "마법사 Lv. 20", mats = { { 4003001, 50 }, { 4011001, 1 }, { 4003000, 10 } }, cost = 5000 },
			{ item = 1382002, label = "마법사 Lv. 25", mats = { { 4021006, 2 }, { 4021001, 1 }, { 4011001, 1 }, { 4003000, 15 } }, cost = 12000 },
			{ item = 1382001, label = "마법사 Lv. 45", mats = { { 4011001, 8 }, { 4021006, 5 }, { 4021001, 5 }, { 4021005, 5 }, { 4003000, 30 }, { 4000010, 50 }, { 4003003, 1 } }, cost = 180000 },
		},
	},
}

return {
	on_click = function(me, npc)
		local type_sel = me:dialog_list(npc, "후후후후... 여기 좋은물건이 많으니 구경이라도 해보고 가지 그래..? \r\n#b", {
			" 장갑 제작",
			" 장갑 합성",
			" 모자 합성",
			" 완드 제작",
			" 스태프 제작",
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
		local prompt = "만들고 싶은 아이템이  #t" .. recipe.item .. "# 1 개 인가? 재료는 아래를 참조하게.\r\n#b"
		for _, mat in ipairs(recipe.mats) do
			prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] .. " 개"
		end
		prompt = prompt .. "\r\n#i4031138# " .. recipe.cost .. " 메소"
		if not me:dialog_yes_no(npc, prompt) then
			return
		end

		if me:meso() < recipe.cost then
			me:dialog(npc, "메소#k 는 제대로 갖고 있는건나? 다시 한번 확인해보게.")
			return
		end
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			cost_items[mat[1]] = mat[2]
		end
		local code = me:exchange({ item = cost_items, meso = recipe.cost }, { item = { [recipe.item] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료는 제대로 갖고 있는건가? 다시 한번 확인해보게. 만일 재료를 장비하고 있다면 장비를 해제하도록 하게. 혹은 인벤토리 공간이 부족한거 아닌가..? 제대로 다시 한번 확인해 보게.")
			return
		end
		me:dialog(npc, "자아.. 다 됐다구. 역시 완벽한 아이템이 탄생했잖아? 다른 작업도 필요하다면 다시 나에게 찾아오라구.")
	end
}
