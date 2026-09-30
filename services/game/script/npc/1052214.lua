-- NPC name (String.wz/Npc.img.xml): 사피

local COIN = 3980002
local TIERS = {
	{ upto = 1, items = { 2040807, 1132152 }, announce = true },
	{ upto = 39, items = { 3980015, 3980015, 3980015, 3980015, 3980015, 3980015, 3980015, 3980015, 2040026, 2040031, 1112401, 1032062, 1132005, 2049100, 2049004 } },
	{ upto = 60, items = { 3980027, 3980027, 3980027, 3980027 } },
	{ upto = 100, items = { 3980018, 4001017 } },
}

return {
	on_click = function(me, npc)
		local message = "#i" .. COIN .. "# #b#z" .. COIN .. "##k1개로 진귀한 아이템을 교환해드려요! \r\n #i" .. COIN .. "# #b#z" .. COIN .. "##k은 후원 으로 획득합니다.\r\n"
		local sel = me:dialog_list(npc, message, {
			"#b#z" .. COIN .. "#을 사용한다.",
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
		if me:exchange({ item = { [COIN] = 1 } }, { item = { [item_id] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 빈 공간이 없거나 #b아르카나코인#k이 부족합니다.")
			return
		end
		me:dialog(npc, "#i" .. item_id .. ":# #b#z" .. item_id .. "##k이(가) 나왔습니다.")
		if tier.announce then
			me:world_message(2, "[후원가챠] : " .. me:name() .. " 님이 1% 소울 벨트 또는 장공100%를 획득하였습니다.")
		end
	end
}
