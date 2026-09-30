-- NPC name (String.wz/Npc.img.xml): 다낭

local COIN = 3980000

local GOODS = {
	{ item = 4001017, cost = 200, extra = "x 1 [200개]" },
	{ item = 4031179, cost = 400, extra = "x 1 [400개]" },
	{ item = 4011006, cost = 100, extra = "x 1 [100개]" },
}

return {
	on_click = function(me, npc)
		local options = {}
		for i, g in ipairs(GOODS) do
			options[i] = "#i" .. g.item .. ":# #b#z" .. g.item .. "##k " .. g.extra
		end
		local sel = me:dialog_list(npc, "#i3980000# #b#z3980000##k을 이용하여 아이템을 교환해드립니다.", options)
		if sel == nil then
			return
		end
		local g = GOODS[sel]
		if not me:dialog_yes_no(npc, "#i" .. g.item .. "# #z" .. g.item .. "# 1개 교환사시겠어요?") then
			return
		end

		local code = me:exchange({ item = { [COIN] = g.cost } }, { item = { [g.item] = 1 } })
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "인벤토리의 빈 공간이 제대로 있는지 확인해주세요.")
			return
		elseif code ~= ExchangeResult.OK then
			me:dialog(npc, "재료 아이템이 제대로 있는지 확인해주세요.")
			return
		end
		me:dialog(npc, "#i" .. g.item .. ":# #b#z" .. g.item .. "##k 1개 완성됐어요.")
	end
}
