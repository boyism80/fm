-- NPC name (String.wz/Npc.img.xml): 할로 캣

local catalogs = {
	{ "전사", { { 1302020, 40000 }, { 1302030, 75000 }, { 1302064, 190000 }, { 1402039, 230000 }, { 1422014, 100000 }, { 1422029, 240000 }, { 1312032, 240000 }, { 1412011, 100000 }, { 1412027, 240000 }, { 1432012, 100000 }, { 1432040, 240000 } } },
	{ "법사", { { 1382009, 35000 }, { 1382012, 80000 }, { 1382039, 230000 }, { 1372034, 230000 } } },
	{ "궁수", { { 1452016, 50000 }, { 1452022, 50000 }, { 1452045, 270000 }, { 1462014, 45000 }, { 1462019, 90000 }, { 1462040, 270000 } } },
	{ "도적", { { 1472030, 40000 }, { 1472032, 70000 }, { 1472055, 270000 }, { 1332057, 10000 }, { 1332025, 70000 }, { 1332055, 270000 } } },
	{ "해적", { { 1492020, 40000 }, { 1492021, 70000 }, { 1492022, 270000 }, { 1482020, 40000 }, { 1482021, 75000 }, { 1482022, 280000 } } },
}

return {
	on_click = function(me, npc)
		local categories = {}
		for i, catalog in ipairs(catalogs) do
			categories[i] = catalog[1]
		end
		local cat = me:dialog_list(npc, "#e[ 소지중인 메소 : " .. me:meso() .. " ]#n", categories)
		if cat == nil then
			return
		end
		local goods = catalogs[cat][2]
		local choices = {}
		for i, row in ipairs(goods) do
			choices[i] = "#i" .. row[1] .. ":# #b#z" .. row[1] .. "##k (가격 : " .. row[2] .. ")"
		end
		local sel = me:dialog_list(npc, "#e[ 소지중인 메소 : " .. me:meso() .. " ]#n", choices)
		if sel == nil then
			return
		end
		local row = goods[sel]
		if row == nil then
			return
		end
		if not me:dialog_yes_no(npc, "#i" .. row[1] .. ":# #b#z" .. row[1] .. "##k을(를) " .. row[2] .. " 메소로 구입하시겠습니까?") then
			return
		end
		local code = me:exchange({ meso = row[2] }, { item = { [row[1]] = 1 } })
		if code == ExchangeResult.LackCost then
			me:dialog(npc, "메소가 " .. row[2] .. "만큼 필요합니다.")
			return
		end
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "아이템창이 부족합니다.")
			return
		end
		me:dialog(npc, "#b#z" .. row[1] .. "##k을(를) " .. row[2] .. " 메소로 구입하셨습니다.")
	end
}
