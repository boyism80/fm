-- NPC name (String.wz/Npc.img.xml): 카오

local COIN = 3980015
local TRADES = {
	{ item = 3980030, qty = 15, price = 1, extra = "x 15 [1개]" },
	{ item = 3980030, qty = 30, price = 2, extra = "x 30 [2개]" },
	{ item = 1902382, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902381, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902401, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902201, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902374, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902145, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902146, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902179, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902186, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902188, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902368, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902437, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902142, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902264, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902372, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902364, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902394, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902387, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902389, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902403, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902435, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902417, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902340, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902348, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902346, qty = 1, price = 20, extra = "x 1 [20개]" },
	{ item = 1902379, qty = 1, price = 20, extra = "x 1 [20개]" },
}

return {
	on_click = function(me, npc)
		local options = {}
		for i, trade in ipairs(TRADES) do
			options[i] = "#i" .. trade.item .. ":# #b#z" .. trade.item .. "##k " .. trade.extra
		end
		local sel = me:dialog_list(npc, "#i3980015# #b#z3980015##k를 이용하여 아이템을 교환해드립니다.", options)
		if sel == nil then
			return
		end

		local trade = TRADES[sel]
		if not me:dialog_yes_no(npc, "#i" .. trade.item .. "# #z" .. trade.item .. "# " .. trade.qty .. "개 교환하시겠어요?") then
			return
		end

		local code = me:exchange({ item = { [COIN] = trade.price } }, { item = { [trade.item] = trade.qty } })
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "인벤토리의 빈 공간이 제대로 있는지 확인해주세요.")
			return
		end
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료 아이템이 제대로 있는지 확인해주세요.")
			return
		end
		me:dialog(npc, "#i" .. trade.item .. ":# #b#z" .. trade.item .. "##k " .. trade.qty .. "개 완성됐어요.")
	end
}
