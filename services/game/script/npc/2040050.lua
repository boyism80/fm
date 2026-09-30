-- NPC name (String.wz/Npc.img.xml): 떠돌이 연금술사

local COST = 4000
local PRODUCTS = { 4006000, 4006001 }
local RECIPES = {
	{
		{ { 4000046, 20 }, { 4000027, 20 }, { 4021001, 1 } },
		{ { 4000025, 20 }, { 4000049, 20 }, { 4021006, 1 } },
		{ { 4000129, 15 }, { 4000130, 15 }, { 4021002, 1 } },
		{ { 4000074, 15 }, { 4000057, 15 }, { 4021005, 1 } },
		{ { 4000054, 7 }, { 4000053, 7 }, { 4021003, 1 } },
	},
	{
		{ { 4000046, 20 }, { 4000027, 20 }, { 4011001, 1 } },
		{ { 4000014, 20 }, { 4000049, 20 }, { 4011003, 1 } },
		{ { 4000132, 15 }, { 4000128, 15 }, { 4011005, 1 } },
		{ { 4000074, 15 }, { 4000069, 15 }, { 4011002, 1 } },
		{ { 4000080, 7 }, { 4000079, 7 }, { 4011004, 1 } },
	},
}
local REFUSE = "재료가 부족하신 모양이군요? 저는 당분간 이곳에 머물 예정이니 재료를 모으셨다면 언제든지 제게 찾아오세요."

return {
	on_click = function(me, npc)
		if not me:dialog(npc, "Alright, mix up the frog's tongue with the squirrel's tooth and ... oh yeah! Forgot to put in the sparkling white powder!! Man, that could have been really bad ... Whoa!! How long have you been standing there? I maaaay have been a little carried away with my work ... hehe.", false, true) then
			me:dialog(npc, REFUSE, false, true)
			return
		end

		local set = me:dialog_list(npc, "음.. 이걸 이렇게 섞고.. 저렇게 섞으면.. 안녕하세요? 저는 떠돌이 연금술사 입니다. 세상 최고의 연금술을 위해서 이곳 저곳을 여행하는 중이지요. 필요한 물건이라도 있으신가요?\r\n\r\n", {
			"#b마법의 돌 제작#k",
			"#b소환의 돌 제작#k",
		})
		if set == nil then
			return
		end

		local product = PRODUCTS[set]
		local options = {}
		for i, recipe in ipairs(RECIPES[set]) do
			options[i] = "#b#t" .. recipe[1][1] .. "#, #t" .. recipe[2][1] .. "##k"
		end
		local sel = me:dialog_list(npc, "하하.. #b#t" .. product .. "##k은 HP와 MP만 소모하는 스킬보다 훨씬 강력한 스킬을 사용할 수 있게 해주는.. 저만 만들 수 있는 특별한 아이템이죠. #t" .. product .. "#를 만드려면 5가지 방법이 있습니다. 어떻게 만들어 보시겠어요?", options)
		if sel == nil then
			return
		end

		local recipe = RECIPES[set][sel]
		local menu = ""
		for _, mat in ipairs(recipe) do
			menu = menu .. "\r\n#v" .. mat[1] .. "# #b" .. mat[2] .. " #t" .. mat[1] .. "##k"
		end
		menu = menu .. "\r\n#i4031138# #b" .. COST .. " 메소#k"
		if not me:dialog_yes_no(npc, "#b5개의 #t" .. product .. "##k을 만드려면, 다음과 같은 재료가 필요합니다. 사냥을 통해서도 구할 수 있고, 다른 플레이어에게서도 구할 수 있습니다. 정말 만들어 보시겠어요?\r\n" .. menu) then
			me:dialog(npc, REFUSE, false, true)
			return
		end

		local cost_items = {}
		for _, mat in ipairs(recipe) do
			cost_items[mat[1]] = mat[2]
		end
		local code = me:exchange({ item = cost_items, meso = COST }, { item = { [product] = 5 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "아이템은 분명 잘 가지고 계신건지, 또는 인벤토리에 빈 공간이 있는지 확인해 주세요.", false, true)
			return
		end
		me:dialog(npc, "여기, #b#t" .. product .. "##k 5개가 있습니다. 제 도움이 더 필요하시면 언제든지 찾아와 주세요!")
	end
}
