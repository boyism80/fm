-- NPC name (String.wz/Npc.img.xml): 루덴

local SHOPS = {
	{
		label = "코인을 정화한다.",
		items = {
			{ id = 4310010, mat = 4310009, qty = 1000 },
		},
	},
	{
		label = "코인을 사용한다.",
		items = {
			{ id = 1312186, mat = 4310010, qty = 4 },
			{ id = 1322237, mat = 4310010, qty = 4 },
			{ id = 1332261, mat = 4310010, qty = 4 },
			{ id = 1372208, mat = 4310010, qty = 4 },
			{ id = 1382246, mat = 4310010, qty = 4 },
			{ id = 1402237, mat = 4310010, qty = 4 },
			{ id = 1412179, mat = 4310010, qty = 4 },
			{ id = 1422186, mat = 4310010, qty = 4 },
			{ id = 1432201, mat = 4310010, qty = 4 },
			{ id = 1442255, mat = 4310010, qty = 4 },
			{ id = 1452239, mat = 4310010, qty = 4 },
			{ id = 1462226, mat = 4310010, qty = 4 },
			{ id = 1472248, mat = 4310010, qty = 4 },
			{ id = 1482203, mat = 4310010, qty = 4 },
			{ id = 1492213, mat = 4310010, qty = 4 },
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
		local sel = me:dialog_list(npc, "사자 상점에 어서오세요~#r노블메달은 #b사자성#k#r 의 '잡몹' 에게서 드롭됩니다!\r\n#b", labels)
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
			me:dialog(npc, "인벤토리의 빈 공간이 제대로 있는지 확인해주세요.")
			return
		end
		me:dialog(npc, "#i" .. entry.id .. ":# #b#z" .. entry.id .. "##k 완성됐어요. 근사하지 않나요?")
	end
}
