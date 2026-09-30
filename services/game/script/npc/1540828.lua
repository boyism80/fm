-- NPC name (String.wz/Npc.img.xml): 라니아

local REWARD = 1022501
local JOURNAL = 3980000

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local message = "안녕하세요 ? 토벌 보상 세탁 담당, 라니아 입니다.\r\n"
		message = message .. "토벌퀘스트 아이템을 세탁, 교환해주고있습니다.\r\n"
		message = message .. "토벌퀘스트 아이템을 세탁 하시겠어요?\r\n"
		local sel = me:dialog_list(npc, message, {
			"#b세탁을 시도한다.#k",
		})
		if sel == nil then
			return
		end
		if me:dialog_yes_no(npc, "세탁 하시려면, #b토벌 보상 아이템#k 1개와 #b퀘스트 일지#k 200개가 필요합니다.") == false then
			return
		end

		if me:empty_slots(InventoryType.Equipment) < 1 then
			me:dialog(npc, "장비창 한칸을 비우시오.")
			return
		end
		if item_count(me, REWARD) < 1 or item_count(me, JOURNAL) < 200 then
			me:dialog(npc, "재료가 부족합니다")
			return
		end

		local cost = { item = { [REWARD] = 1, [JOURNAL] = 200 } }
		local reward = {
			item = { [REWARD] = 1 },
			bonus = { str = 8, dex = 8, int = 15, luk = 15, watk = 5, matk = 3 },
		}
		if me:exchange(cost, reward) ~= ExchangeResult.OK then
			me:dialog(npc, "재료가 부족합니다")
			return
		end
		me:dialog(npc, "전리품을 교환하여 아이템이 지급되었습니다.무슨 아이템이 나왔는지는 장비템을 확인해주세요!")
	end
}
