-- NPC name (String.wz/Npc.img.xml): 할로 캣

local TICKET = 3980005
local TIERS = {
	{ upto = 1, items = { 1032035 }, announce = true },
	{ upto = 39, items = { 3010131, 3010132, 3010682, 3010690, 3010601, 3010064, 3010009, 3010574, 3010590, 3010674, 3010113, 3010319, 3010086, 3010133, 3010694, 3010065 } },
	{ upto = 60, items = { 2049100, 2049100, 2049100, 2049100 } },
	{ upto = 100, items = { 2049100, 2049100 } },
}

return {
	on_click = function(me, npc)
		local message = "#i" .. TICKET .. "# #b#z" .. TICKET .. "##k1개로 진귀한 아이템을 교환해드려요! \r\n #i" .. TICKET .. "# #b#z" .. TICKET .. "##k은 핫타임 으로 획득합니다.\r\n"
		local sel = me:dialog_list(npc, message, {
			"#b#z" .. TICKET .. "#을 사용한다.",
		})
		if sel == nil then
			return
		end

		local roll = math.random(1, 100)
		local tier = nil
		for _, t in ipairs(TIERS) do
			if roll <= t.upto then
				tier = t
				break
			end
		end
		local item_id = tier.items[math.random(1, #tier.items)]
		if me:exchange({ item = { [TICKET] = 1 } }, { item = { [item_id] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 빈 공간이 없거나 #b요정의특제간장#k이 부족합니다.")
			return
		end
		me:dialog(npc, "#i" .. item_id .. ":# #b#z" .. item_id .. "##k이(가) 나왔습니다.")
		if tier.announce then
			me:message("[할로캣 핫타임] : " .. me:name() .. " 님이 1% 확률로 정령의힘이 깃든 귀고리를 획득하셨습니다. ", Msg.Megaphone, MessageScope.World)
		end
	end
}
