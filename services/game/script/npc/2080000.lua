-- NPC name (String.wz/Npc.img.xml): 모스

local CATEGORIES = {
	{
		recipes = {
			{ item = 1302059, label = "레벨. 110 한손검", stimulator = 4130002, mats = { { 1302056, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 8 } } },
			{ item = 1312031, label = "레벨. 110 한손도끼", stimulator = 4130003, mats = { { 1312030, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 8 } } },
			{ item = 1322052, label = "레벨. 110 한손둔기", stimulator = 4130004, mats = { { 1322045, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 8 } } },
			{ item = 1402036, label = "레벨. 110 두손검", stimulator = 4130005, mats = { { 1402035, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 8 } } },
			{ item = 1412026, label = "레벨. 110 두손도끼", stimulator = 4130006, mats = { { 1412021, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 8 } } },
			{ item = 1422028, label = "레벨. 110 두손둔기", stimulator = 4130007, mats = { { 1422027, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 8 } } },
			{ item = 1432038, label = "레벨. 110 창", stimulator = 4130008, mats = { { 1432030, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 8 } } },
			{ item = 1442045, label = "레벨. 110 폴암", stimulator = 4130009, mats = { { 1442044, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 8 } } },
		},
	},
	{
		recipes = {
			{ item = 1452044, label = "레벨. 110 활", stimulator = 4130012, mats = { { 1452019, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 3 }, { 4005002, 5 } } },
			{ item = 1462039, label = "레벨. 110 석궁", stimulator = 4130013, mats = { { 1462015, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 5 }, { 4005002, 3 } } },
		},
	},
	{
		recipes = {
			{ item = 1372032, label = "레벨. 108 완드", stimulator = 4130010, mats = { { 1372010, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005001, 6 }, { 4005003, 2 } } },
			{ item = 1382036, label = "레벨. 110 스태프", stimulator = 4130011, mats = { { 1382035, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005001, 6 }, { 4005003, 2 } } },
		},
	},
	{
		recipes = {
			{ item = 1332049, label = "레벨. 110 STR 단검", stimulator = 4130014, mats = { { 1332051, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005000, 5 }, { 4005002, 3 } } },
			{ item = 1332050, label = "레벨. 110 LUK 단검", stimulator = 4130014, mats = { { 1332052, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005002, 3 }, { 4005003, 5 } } },
			{ item = 1472051, label = "레벨. 110 아대", stimulator = 4130015, mats = { { 1472053, 1 }, { 4000244, 20 }, { 4000245, 25 }, { 4005002, 2 }, { 4005003, 6 } } },
		},
	},
}
local COST = 120000
local BUSTED_DAGGER = 4001079

return {
	on_click = function(me, npc)
		local busted = 0
		for _, it in pairs(me:item(BUSTED_DAGGER)) do
			busted = busted + it:count()
		end
		if busted >= 1 then
			if not me:dialog_yes_no(npc, "#b코니언의 단도#k가 필요하다는 거지? #b코니언의 단도#k를 만드려면 #b망가진 단도#k를 수리하는 방법밖에 없어. #b망가진 단도#k를 수리하려면 다음과 같은 재료가 필요하지. 만들어보겠어?\r\n\r\n#b#i4001079# #z4001079# 1개\r\n#i4011001# #z4011001# 1개\r\n#i4011002# #z4011002# 1개") then
				return
			end
			local code = me:exchange({ item = { [4011001] = 1, [4011002] = 1, [BUSTED_DAGGER] = 1 } }, { item = { [4001078] = 1 } })
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "아이템은 제대로 모아 온거야? 그렇지 않은 것 같은데? 혹은 인벤토리 공간이 부족한건 아닌지 확인해 봐.")
				return
			end
			me:dialog(npc, "어때? 새 것 같지? 잃어버리거나 새로 만들어야 한다면 다시 나를 찾아오라고.")
			return
		end

		local type_sel = me:dialog_list(npc, "용의 힘은 아주 강력하지. 원한다면 그 힘을 네 무기에 담아줄게. 하지만 그 힘을 무기가 과연 버틸 수 있을지..#b", {
			" 촉진제가 뭐죠?",
			" 전사 무기 제작",
			" 궁수 무기 제작",
			" 마법사 무기 제작",
			" 도적 무기 제작",
			" 촉진제로 전사 무기 제작",
			" 촉진제로 궁수 무기 제작",
			" 촉진제로 마법사 무기 제작",
			" 촉진제로 도적 무기 제작",
		})
		if type_sel == nil then
			return
		end
		if type_sel == 1 then
			me:dialog(npc, "촉진제는 여러 아이템을 만들때 첨가할 수 있는 특수한 약이지. 몬스터에게서 구할 수 있을거야. 어쨌든, 촉진제는 아이템을 만들때 일정 확률로 옵션이 더 좋아질 수 있게 되지. 하지만 10%의 확률로 제작에 실패할 수도 있으니 신중하게 결정해야하지.", false, true)
			return
		end

		local stimulator = type_sel > 5
		local category = CATEGORIES[type_sel - 1]
		if stimulator then
			category = CATEGORIES[type_sel - 5]
		end
		local options = {}
		for i, recipe in ipairs(category.recipes) do
			options[i] = " #z" .. recipe.item .. "##k - " .. recipe.label .. "#b"
		end
		local item_sel = me:dialog_list(npc, "좋아. 드래곤의 힘을 담을 무기를 골라봐.#b", options)
		if item_sel == nil then
			return
		end

		local recipe = category.recipes[item_sel]
		local prompt = "흐음. 필요한 아이템이 #t" .. recipe.item .. "# 인거야? 재료는 다음과 같아. 인벤토리 공간이 충분한지도 확인해 보도록 해.#b"
		if stimulator then
			prompt = prompt .. "\r\n#i" .. recipe.stimulator .. "# #t" .. recipe.stimulator .. "# 1개"
		end
		for _, mat in ipairs(recipe.mats) do
			prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] .. "개"
		end
		prompt = prompt .. "\r\n#i4031138# " .. COST .. " 메소"
		if not me:dialog_yes_no(npc, prompt) then
			return
		end

		if me:meso() < COST then
			me:dialog(npc, "흐음.. 이정도의 수수료면 싼 편인데.. 그정도 돈도 없는거야?")
			return
		end
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			cost_items[mat[1]] = mat[2]
		end

		if stimulator == false then
			local code = me:exchange({ item = cost_items, meso = COST }, { item = { [recipe.item] = 1 } })
			if code == ExchangeResult.LackCapacity then
				me:dialog(npc, "인벤토리 공간이 부족한걸?")
				return
			end
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "재료가 부족한 것 같은데? 아니면 인벤토리 공간이 부족한건 아닌지 확인해 봐!")
				return
			end
			me:dialog(npc, "자, 여기 드래곤의 힘이 담긴 무기야. 유용하게 쓰도록 해.")
			return
		end

		cost_items[recipe.stimulator] = 1
		local fail = "재료가 부족한 것 같은데? 아니면 인벤토리 공간이 부족한건 아닌지 확인해 봐!"
		if math.random(0, 9) == 0 then
			local code = me:exchange({ item = cost_items, meso = COST }, nil)
			if code ~= ExchangeResult.OK then
				me:dialog(npc, fail)
				return
			end
			me:dialog(npc, "음.. 이 무기는 드래곤의 힘을 담기에는 너무 약했었나봐.. 무기가 부서진 것 같아.")
			return
		end
		local code = me:exchange({ item = cost_items, meso = COST }, { item = { [recipe.item] = 1 }, random = true })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, fail)
			return
		end
		me:dialog(npc, "자, 여기 드래곤의 힘이 담긴 무기야. 유용하게 쓰도록 해.")
	end
}
