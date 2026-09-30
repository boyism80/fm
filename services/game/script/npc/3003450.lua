-- NPC name (String.wz/Npc.img.xml): 타나

local COIN = 3980001
local REFUSE = "길가다가 넘어져서 한강으로 굴러떨어져라."
local CATEGORIES = {
	{
		items = { 2049100, 2044601, 2044001, 2043001, 2044501, 2043801, 2043701, 2043301, 2044701, 2044301, 2044901, 2044801, 2040625, 2040418, 2040425, 2040718, 2040321, 2040317, 2040029, 2040516 },
		costs = { 9500, 4500, 4500, 4500, 4500, 4500, 4500, 4500, 4500, 4500, 4500, 4500, 6000, 6000, 6000, 300, 6000, 6000, 6000, 6000 },
	},
	{
		items = { 4170006 },
		costs = { 10 },
	},
	{
		items = { 5060002, 4031461, 5520000, 4031455, 4213001 },
		costs = { 1000, 2000, 1000, 2500, 1000 },
	},
	{
		items = { 1402013, 1052676, 1004650, 1142378 },
		costs = { 2000, 1000, 10000, 15000 },
	},
}

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "특별한 아이템을 교환하기위해선 특정 개수의 낚시코인이 필요합니다.", {
			"#b주문서류#k",
			"#b피그미에그#k",
			"#b물약/특수아이템류#k",
			"#b코디류#k",
		})
		if sel == nil then
			me:dialog(npc, REFUSE)
			return
		end

		local category = CATEGORIES[sel]
		local options = {}
		for i, id in ipairs(category.items) do
			options[i] = "#i" .. COIN .. ":# " .. category.costs[i] .. "개#k   =  #i" .. id .. ":# #b#z" .. id .. ":##k"
		end
		local item_sel = me:dialog_list(npc, "특별한 아이템을 교환 하기 위해선 낚시코인이 필요합니다.", options)
		if item_sel == nil then
			me:dialog(npc, REFUSE)
			return
		end

		local code = me:exchange({ item = { [COIN] = category.costs[item_sel] } }, { item = { [category.items[item_sel]] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "낚시코인이 부족하거나 인벤토리 빈 공간이 부족합니다.")
			return
		end
		me:dialog(npc, "아이템을 교환하였습니다. 즐거운 메이플스토리 되세요~")
	end
}
