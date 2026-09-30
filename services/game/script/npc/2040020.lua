-- NPC name (String.wz/Npc.img.xml): 지로큰

local STIMULATOR = 4130000

local CATEGORIES = {
	{
		prompt = "전사 장갑이 필요하신가요? 어떤 아이템이 필요하세요?#b",
		recipes = {
			{ item = 1082007, label = "전사 레벨. 30", cost = 18000, mats = { { 4011000, 3 }, { 4011001, 2 }, { 4003000, 15 } } },
			{ item = 1082008, label = "전사 레벨. 35", cost = 27000, mats = { { 4000021, 30 }, { 4011001, 4 }, { 4003000, 15 } } },
			{ item = 1082023, label = "전사 레벨. 40", cost = 36000, mats = { { 4000021, 50 }, { 4011001, 5 }, { 4003000, 40 } } },
			{ item = 1082009, label = "전사 레벨. 50", cost = 45000, mats = { { 4011001, 3 }, { 4021007, 2 }, { 4000030, 30 }, { 4003000, 45 } } },
		},
	},
	{
		prompt = "궁수 장갑이 필요하신가요? 어떤 아이템이 필요하세요?#b",
		recipes = {
			{ item = 1082048, label = "궁수 레벨. 30", cost = 18000, mats = { { 4000021, 50 }, { 4011006, 2 }, { 4021001, 1 } } },
			{ item = 1082068, label = "궁수 레벨. 35", cost = 27000, mats = { { 4011000, 1 }, { 4011001, 3 }, { 4000021, 60 }, { 4003000, 15 } } },
			{ item = 1082071, label = "궁수 레벨. 40", cost = 36000, mats = { { 4011001, 3 }, { 4021000, 1 }, { 4021002, 3 }, { 4000021, 80 }, { 4003000, 25 } } },
			{ item = 1082084, label = "궁수 레벨. 50", cost = 45000, mats = { { 4011004, 3 }, { 4011006, 1 }, { 4021002, 2 }, { 4000030, 40 }, { 4003000, 35 } } },
		},
	},
	{
		prompt = "마법사 장갑이 필요하신가요? 어떤 아이템이 필요하세요?#b",
		recipes = {
			{ item = 1082051, label = "마법사 레벨. 30", cost = 22500, mats = { { 4000021, 60 }, { 4021006, 1 }, { 4021000, 2 } } },
			{ item = 1082054, label = "마법사 레벨. 35", cost = 27000, mats = { { 4000021, 70 }, { 4011006, 1 }, { 4011001, 3 }, { 4021000, 2 } } },
			{ item = 1082062, label = "마법사 레벨. 40", cost = 36000, mats = { { 4000021, 80 }, { 4021000, 3 }, { 4021006, 3 }, { 4003000, 30 } } },
			{ item = 1082081, label = "마법사 레벨. 50", cost = 45000, mats = { { 4021000, 3 }, { 4011006, 2 }, { 4000030, 35 }, { 4003000, 40 } } },
		},
	},
	{
		prompt = "돚거 장갑이 필요하신가요? 어떤 아이템이 필요하세요?#b",
		recipes = {
			{ item = 1082042, label = "도적 레벨. 30", cost = 22500, mats = { { 4011001, 2 }, { 4000021, 50 }, { 4003000, 10 } } },
			{ item = 1082046, label = "도적 레벨. 35", cost = 27000, mats = { { 4011001, 3 }, { 4011000, 1 }, { 4000021, 60 }, { 4003000, 15 } } },
			{ item = 1082075, label = "도적 레벨. 40", cost = 36000, mats = { { 4021000, 3 }, { 4000101, 100 }, { 4000021, 80 }, { 4003000, 30 } } },
			{ item = 1082065, label = "도적 레벨. 50", cost = 45000, mats = { { 4021005, 3 }, { 4021008, 1 }, { 4000030, 40 }, { 4003000, 30 } } },
		},
	},
	{
		prompt = "촉진제를 통한 전사 장갑이 필요하신가요? 어떤 아이템이 필요하세요?#b",
		stimulator = true,
		recipes = {
			{ item = 1082005, label = "전사 레벨. 30", cost = 18000, mats = { { 1082007, 1 }, { 4011001, 1 } } },
			{ item = 1082006, label = "전사 레벨. 30", cost = 22500, mats = { { 1082007, 1 }, { 4011005, 2 } } },
			{ item = 1082035, label = "전사 레벨. 35", cost = 27000, mats = { { 1082008, 1 }, { 4021006, 3 } } },
			{ item = 1082036, label = "전사 레벨. 35", cost = 36000, mats = { { 1082008, 1 }, { 4021008, 1 } } },
			{ item = 1082024, label = "전사 레벨. 40", cost = 40500, mats = { { 1082023, 1 }, { 4011003, 4 } } },
			{ item = 1082025, label = "전사 레벨. 40", cost = 45000, mats = { { 1082023, 1 }, { 4021008, 2 } } },
			{ item = 1082010, label = "전사 레벨. 50", cost = 49500, mats = { { 1082009, 1 }, { 4011002, 5 } } },
			{ item = 1082011, label = "전사 레벨. 50", cost = 54000, mats = { { 1082009, 1 }, { 4011006, 4 } } },
		},
	},
	{
		prompt = "촉진제를 통한 궁수 장갑이 필요하신가요? 어떤 아이템이 필요하세요?#b",
		stimulator = true,
		recipes = {
			{ item = 1082049, label = "궁수 레벨. 30", cost = 13500, mats = { { 1082048, 1 }, { 4021003, 3 } } },
			{ item = 1082050, label = "궁수 레벨. 30", cost = 18000, mats = { { 1082048, 1 }, { 4021008, 1 } } },
			{ item = 1082069, label = "궁수 레벨. 35", cost = 19800, mats = { { 1082068, 1 }, { 4011002, 4 } } },
			{ item = 1082070, label = "궁수 레벨. 35", cost = 22500, mats = { { 1082068, 1 }, { 4011006, 2 } } },
			{ item = 1082072, label = "궁수 레벨. 40", cost = 27000, mats = { { 1082071, 1 }, { 4011006, 4 } } },
			{ item = 1082073, label = "궁수 레벨. 40", cost = 36000, mats = { { 1082071, 1 }, { 4021008, 2 } } },
			{ item = 1082085, label = "궁수 레벨. 50", cost = 49500, mats = { { 1082084, 1 }, { 4011000, 1 }, { 4021000, 5 } } },
			{ item = 1082083, label = "궁수 레벨. 50", cost = 54000, mats = { { 1082084, 1 }, { 4011006, 2 }, { 4021008, 2 } } },
		},
	},
	{
		prompt = "촉진제를 통한 마법사 장갑이 필요하신가요? 어떤 아이템이 필요하세요?#b",
		stimulator = true,
		recipes = {
			{ item = 1082052, label = "마법사 레벨. 30", cost = 31500, mats = { { 1082051, 1 }, { 4021005, 3 } } },
			{ item = 1082053, label = "마법사 레벨. 30", cost = 36000, mats = { { 1082051, 1 }, { 4021008, 1 } } },
			{ item = 1082055, label = "마법사 레벨. 35", cost = 36000, mats = { { 1082054, 1 }, { 4021005, 3 } } },
			{ item = 1082056, label = "마법사 레벨. 35", cost = 40500, mats = { { 1082054, 1 }, { 4021008, 1 } } },
			{ item = 1082063, label = "마법사 레벨. 40", cost = 40500, mats = { { 1082062, 1 }, { 4021002, 4 } } },
			{ item = 1082064, label = "마법사 레벨. 40", cost = 45000, mats = { { 1082062, 1 }, { 4021008, 2 } } },
			{ item = 1082082, label = "마법사 레벨. 50", cost = 49500, mats = { { 1082081, 1 }, { 4021002, 5 } } },
			{ item = 1082080, label = "마법사 레벨. 50", cost = 54000, mats = { { 1082081, 1 }, { 4021008, 3 } } },
		},
	},
	{
		prompt = "촉진제를 통한 돚거 장갑이 필요하신가요? 어떤 아이템이 필요하세요?#b",
		stimulator = true,
		recipes = {
			{ item = 1082043, label = "도적 레벨. 30", cost = 13500, mats = { { 1082042, 1 }, { 4011004, 2 } } },
			{ item = 1082044, label = "도적 레벨. 30", cost = 18000, mats = { { 1082042, 1 }, { 4011006, 1 } } },
			{ item = 1082047, label = "도적 레벨. 35", cost = 19800, mats = { { 1082046, 1 }, { 4011005, 3 } } },
			{ item = 1082045, label = "도적 레벨. 35", cost = 22500, mats = { { 1082046, 1 }, { 4011006, 2 } } },
			{ item = 1082076, label = "도적 레벨. 40", cost = 36000, mats = { { 1082075, 1 }, { 4011006, 4 } } },
			{ item = 1082074, label = "도적 레벨. 40", cost = 45000, mats = { { 1082075, 1 }, { 4021008, 2 } } },
			{ item = 1082067, label = "도적 레벨. 50", cost = 49500, mats = { { 1082065, 1 }, { 4021000, 5 } } },
			{ item = 1082066, label = "도적 레벨. 50", cost = 54000, mats = { { 1082065, 1 }, { 4011006, 2 }, { 4021008, 1 } } },
		},
	},
}

return {
	on_click = function(me, npc)
		local type_sel = me:dialog_list(npc, "안녕하세요, 루디브리엄 장갑 상점이에요. 무엇을 도와드릴까요?#b", {
			" 촉진제가 뭐죠?",
			" 전사 장갑 제작",
			" 궁수 장갑 제작",
			" 마법사 장갑 제작",
			" 도적 장갑 제작",
			" 촉진제로 전사 장갑 제작",
			" 촉진제로 궁수 장갑 제작",
			" 촉진제로 마법사 장갑 제작",
			" 촉진제로 돚거 장갑 제작",
		})
		if type_sel == nil then
			return
		end
		if type_sel == 1 then
			me:dialog(npc, "촉진제는 여러 아이템을 만들때 첨가할 수 있는 특수한 약이죠. 몬스터에게서 구하실 수 있을거에요. 어쨌든, 촉진제는 아이템을 만들때 일정 확률로 옵션이 더 좋아질 수 있게 됩니다. 하지만 10%의 확률로 제작에 실패할 수도 있으니 신중하게 결정해주세요.", false, true)
			return
		end

		local category = CATEGORIES[type_sel - 1]
		local options = {}
		for i, recipe in ipairs(category.recipes) do
			options[i] = " #z" .. recipe.item .. "##k - " .. recipe.label .. "#b"
		end
		local item_sel = me:dialog_list(npc, category.prompt, options)
		if item_sel == nil then
			return
		end

		local recipe = category.recipes[item_sel]
		local prompt = "#t" .. recipe.item .. "# 아이템이 필요한가요? 그렇다면 재료는 아래와 같이 모아오셔야 합니다. 그리고 인벤토리 공간이 충분한지도 확인해 주세요.\r\n#b"
		for _, mat in ipairs(recipe.mats) do
			prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] .. "개"
		end
		prompt = prompt .. "\r\n#i4031138# " .. recipe.cost .. " 메소"
		if not me:dialog_yes_no(npc, prompt) then
			return
		end

		if me:meso() < recipe.cost then
			me:dialog(npc, "음.. 메소가 부족하신 것 같은데요?")
			return
		end
		local fail = "죄송하지만 재료를 올바르게 가져오신 것 같지 않네요. 혹은 인벤토리 공간이 부족한건 아니신가요?"
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			cost_items[mat[1]] = mat[2]
		end

		if category.stimulator ~= true then
			local code = me:exchange({ item = cost_items, meso = recipe.cost }, { item = { [recipe.item] = 1 } })
			if code ~= ExchangeResult.OK then
				me:dialog(npc, fail)
				return
			end
			me:dialog(npc, "자, 장갑이 다 만들어 졌어요. 조심하세요. 아직 뜨겁거든요.")
			return
		end

		cost_items[STIMULATOR] = 1
		if math.random(0, 9) == 0 then
			local code = me:exchange({ item = cost_items, meso = recipe.cost }, nil)
			if code ~= ExchangeResult.OK then
				me:dialog(npc, fail)
				return
			end
			me:dialog(npc, "헉! 아무래도.. 촉진제를 사용하다가 사고가 발생한 모양이에요. 아쉽하지만 촉진제로 아이템 제작에 실패한 것 같네요.")
			return
		end
		local code = me:exchange({ item = cost_items, meso = recipe.cost }, { item = { [recipe.item] = 1 }, random = true })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, fail)
			return
		end
		me:dialog(npc, "자, 장갑이 다 만들어 졌어요. 조심하세요. 아직 뜨겁거든요.")
	end
}
