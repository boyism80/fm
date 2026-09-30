-- NPC name (String.wz/Npc.img.xml): 센츠

local SHOPS = {
	{
		label = "여제 옷 을 구매한다.",
		items = {
			{ id = 1052314, mat = 4021019, qty = 10 },
			{ id = 1052315, mat = 4021019, qty = 10 },
			{ id = 1052318, mat = 4021019, qty = 10 },
			{ id = 1052316, mat = 4021019, qty = 10 },
			{ id = 1052317, mat = 4021019, qty = 10 },
		},
	},
	{
		label = "여제 망토를 구매한다.",
		items = {
			{ id = 1102278, mat = 4021019, qty = 45 },
			{ id = 1102277, mat = 4021019, qty = 45 },
			{ id = 1102276, mat = 4021019, qty = 45 },
			{ id = 1102275, mat = 4021019, qty = 45 },
			{ id = 1102279, mat = 4021019, qty = 45 },
		},
	},
}

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local labels = {}
		for _, shop in ipairs(SHOPS) do
			table.insert(labels, shop.label)
		end
		local sel = me:dialog_list(npc, "영주 직할령, 특별 상점에 어서오세요~#r[꿈의돌]#k은 #k광장#k의 #b'리린'#k 에게서 제작 가능합니다!\r\n#b", labels)
		if sel == nil then
			return
		end

		local shop = SHOPS[sel]
		local options = {}
		for _, entry in ipairs(shop.items) do
			table.insert(options, "#i" .. entry.id .. ":# #b#z" .. entry.id .. "##k")
		end
		local pick = me:dialog_list(npc, "어떤 아이템을 제작하고 싶으신가요?\r\n필요한 아이템이 있으시다면 선택해주세요.", options)
		if pick == nil then
			return
		end

		local entry = shop.items[pick]
		local prompt = "#i" .. entry.id .. ":# #b#z" .. entry.id .. "##k 제작 하시겠습니까?\r\n재료 아이템의 개수가 제대로 있는지 확인해보세요.\r\n"
		prompt = prompt .. "\r\n#i" .. entry.mat .. ":# #z" .. entry.mat .. "# ( #b#c" .. entry.mat .. "#개#k / #r" .. entry.qty .. "개#k ) "
		if me:dialog_yes_no(npc, prompt) == false then
			return
		end

		if item_count(me, entry.mat) < entry.qty then
			me:dialog(npc, "재료 아이템이 제대로 있는지 확인해주세요.")
			return
		end
		if me:exchange({ item = { [entry.mat] = entry.qty } }, { item = { [entry.id] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "재료가 제대로 있거나 인벤토리의 빈 공간이 제대로 있는지 확인해주세요.")
			return
		end
		me:dialog(npc, "#i" .. entry.id .. ":# #b#z" .. entry.id .. "##k 완성됐어요. 근사하지 않나요?")
	end
}
