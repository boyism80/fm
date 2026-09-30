-- NPC name (String.wz/Npc.img.xml): 라니아

local GLOVE = 1082149
local BELT = 1132151
local CHICKEN = 2022096

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local message = "안녕하세요 ? 파티퀘스트 컨텐츠 담당, 라니아 입니다.\r\n"
		message = message .. "갈색노가다목장갑을 세탁, 교환해주고있습니다.\r\n"
		message = message .. "장갑의 세탁 시도 하시겠어요? #r주의: 장갑이 아닌 다른 아이템이 나올 수 있습니다.\r\n"
		local sel = me:dialog_list(npc, message, {
			"#b세탁을 시도한다.#k",
		})
		if sel == nil then
			return
		end
		if me:dialog_yes_no(npc, "파퀘 아이템을 세탁 하시려면, #b갈색노가다목장갑#k 1개와 #b닭튀김#k 5개가 필요합니다. 세탁시, 벨트나 장갑이 랜덤으로 나타납니다.") == false then
			return
		end

		if me:empty_slots(InventoryType.Equipment) < 1 then
			me:dialog(npc, "장비창 한칸을 비우시오.")
			return
		end
		if item_count(me, GLOVE) < 1 or item_count(me, CHICKEN) < 5 then
			me:dialog(npc, "재료가 부족합니다")
			return
		end

		local reward = nil
		if math.random(0, 1) == 0 then
			reward = {
				item = { [GLOVE] = 1 },
				bonus = { watk = 5, matk = 1 },
			}
		else
			reward = {
				item = { [BELT] = 1 },
				bonus = { str = 30, dex = 30, int = 30, luk = 30, watk = 5, matk = 1 },
			}
		end
		if me:exchange({ item = { [GLOVE] = 1, [CHICKEN] = 5 } }, reward) ~= ExchangeResult.OK then
			me:dialog(npc, "재료가 부족합니다")
			return
		end
		me:dialog(npc, "전리품을 교환하여 아이템이 지급되었습니다.무슨 아이템이 나왔는지는 장비템을 확인해주세요!")
	end
}
