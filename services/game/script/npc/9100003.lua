-- NPC name (String.wz/Npc.img.xml): 페리온 마네키네코

local items = {
	1012289, 1072407, 1072395, 1072394, 1072483, 1072457, 1072426, 1001084, 1072514, 1001083,
	1002845, 1051371, 1050303, 1051285, 1012137, 1032063, 1050235, 1002907, 1002920, 1002921,
	1002970, 1002877, 1002978, 1062136, 1062138, 1062139, 1062147, 1062145, 1062153, 1062152,
	1062157, 1062107, 1062119, 1062122, 1062123, 1062124, 1062126, 1062131, 1062113, 1062110,
	1062108, 1042222, 1042200, 1042202, 1042198, 1042199, 1042193, 1042194, 1042188, 1042207,
	1042168, 1042170, 1042163, 1042183, 1042157, 1042158, 1042156, 1092067, 1342069, 1022110,
	1022079, 1022075, 1022109, 1022083, 1022102, 1052354, 1102532, 1102267, 1052253, 1050177,
	1052291, 1051219, 1051218, 1052203, 1051189, 1052236, 1052224, 1052200, 1052179, 1102270,
	1102157, 1102223, 1102224, 1102215,
}

local PAGE = 50

local function show_page(me, npc, page)
	local first = (page - 1) * PAGE + 1
	if first > #items or first < 1 then
		me:dialog(npc, "그 페이지는 없습니다.")
		return
	end
	while first <= #items do
		local last = math.min(first + PAGE - 1, #items)
		local text = "현재 페이지 : " .. page .. "/" .. math.floor((#items - 1) / PAGE + 1) .. "\r\n\r\n"
		for i = first, last do
			text = text .. "#i" .. items[i] .. "# #z" .. items[i] .. "#\r\n"
		end
		local more = last < #items
		local forward = me:dialog(npc, text, false, more)
		if more == false or forward == false then
			return
		end
		page = page + 1
		first = last + 1
	end
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "부화기 아이템 검색", { "페이지 찾아가기" })
		if sel == nil then
			return
		end
		local text = me:dialog_input(npc, "페이지 적어")
		local page = tonumber(text)
		if page == nil then
			return
		end
		show_page(me, npc, page)
	end
}
