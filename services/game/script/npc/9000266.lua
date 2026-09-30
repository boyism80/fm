-- NPC name (String.wz/Npc.img.xml): 밤둥이

local CATEGORIES = {
	{
		text = "10% 무기 주문서",
		goods = { { 2044502, 4000000 }, { 2043302, 4000000 }, { 2043002, 4000000 }, { 2044002, 4000000 }, { 2044702, 4000000 }, { 2043202, 4000000 }, { 2044202, 4000000 }, { 2044302, 4000000 }, { 2044402, 4000000 }, { 2044802, 500000 }, { 2044902, 500000 }, { 2044602, 4000000 }, { 2043702, 4000000 }, { 2043802, 4000000 } },
	},
	{
		text = "100% 무기 주문서",
		goods = { { 2043303, 100000000 }, { 2043003, 100000000 }, { 2044003, 100000000 }, { 2043703, 100000000 }, { 2043803, 100000000 }, { 2044503, 100000000 }, { 2044603, 100000000 }, { 2044303, 100000000 }, { 2044403, 100000000 }, { 2043203, 100000000 }, { 2044203, 100000000 }, { 2044703, 100000000 } },
	},
	{
		text = "10% 방어구 주문서",
		goods = { { 2041017, 3000000 }, { 2040302, 3000000 }, { 2040026, 3000000 }, { 2040514, 3000000 }, { 2040805, 5000000 } },
	},
	{
		text = "100%방어구 무기",
		goods = { { 2040303, 100000000 }, { 2040807, 200000000 } },
	},
	{
		text = "방어구/유틸 상점",
		goods = { { 1082149, 2000000000 }, { 1132152, 2000000000 }, { 4031348, 5000000 }, { 2049004, 25000000 } },
	},
}

return {
	on_click = function(me, npc)
		local names = {}
		for i, c in ipairs(CATEGORIES) do
			names[i] = c.text
		end
		local cat = me:dialog_list(npc, "주문서 상점입니다! 해적은 #e100% 무기 주문서#n 를 구할 수 없습니다.#n", names)
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
