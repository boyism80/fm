-- NPC name (String.wz/Npc.img.xml): 주황버섯

local CATEGORIES = {
	{
		text = "전사 무기",
		goods = { { 1302020, 40000 }, { 1302030, 75000 }, { 1302064, 190000 }, { 1402039, 230000 }, { 1422014, 100000 }, { 1422029, 240000 }, { 1312032, 240000 }, { 1412011, 100000 }, { 1412027, 240000 }, { 1432012, 100000 }, { 1432040, 240000 } },
	},
	{
		text = "법사 무기",
		goods = { { 1382009, 35000 }, { 1382012, 80000 }, { 1382039, 230000 }, { 1372034, 230000 } },
	},
	{
		text = "궁수 무기",
		goods = { { 1452016, 50000 }, { 1452022, 50000 }, { 1452045, 270000 }, { 1462014, 45000 }, { 1462019, 90000 }, { 1462040, 270000 } },
	},
	{
		text = "도적 무기",
		goods = { { 1472030, 40000 }, { 1472032, 70000 }, { 1472055, 270000 }, { 1332057, 10000 }, { 1332025, 70000 }, { 1332055, 270000 } },
	},
	{
		text = "해적 무기",
		goods = { { 1492020, 40000 }, { 1492021, 70000 }, { 1492022, 270000 }, { 1482020, 40000 }, { 1482021, 75000 }, { 1482022, 280000 } },
	},
}

return {
	on_click = function(me, npc)
		local names = {}
		for i, c in ipairs(CATEGORIES) do
			names[i] = c.text
		end
		local cat = me:dialog_list(npc, "#e[ 소지중인 메소 : " .. me:meso() .. " ]#n", names)
		if cat == nil then
			return
		end

		local goods = {}
		local options = {}
		for _, g in ipairs(CATEGORIES[cat].goods) do
			if item_wz(g[1]) ~= nil then
				goods[#goods + 1] = g
				options[#options + 1] = "#i" .. g[1] .. ":##b#z" .. g[1] .. "##k (가격 : " .. g[2] .. ")"
			end
		end
		local sel = me:dialog_list(npc, "	#e[ 소지중인 메소 : " .. me:meso() .. " ]#n", options)
		if sel == nil then
			return
		end
		local item = goods[sel][1]
		local price = goods[sel][2]
		if me:meso() < price then
			me:dialog(npc, "메소가 " .. price .. "만큼 필요합니다.")
			return
		end
		if not me:dialog_yes_no(npc, "#i" .. item .. ":# #b#z" .. item .. "##k을(를) " .. price .. " 메소로 구입하시겠습니까?") then
			return
		end

		local code = me:exchange({ meso = price }, { item = { [item] = 1 } })
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "아이템창이 부족합니다.")
			return
		elseif code ~= ExchangeResult.OK then
			me:dialog(npc, "메소가 " .. price .. "만큼 필요합니다.")
			return
		end
		me:dialog(npc, "#b#z" .. item .. "##k을(를) " .. price .. " 메소로 구입하셨습니다.")
	end
}
