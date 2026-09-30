-- NPC name (String.wz/Npc.img.xml): 타나

local COIN = 3980003
local GLOVE = 1082149
local BELT = 1132005

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local message = "안녕하세요 ? 파티퀘스트 컨텐츠 담당, 타나 입니다.\r\n"
		message = message .. "파티퀘스트로 얻는 전리품을 진귀한 아이템으로 교환해주고있습니다.\r\n"
		message = message .. "전리품을 교환하시겠어요 ?\r\n"
		local sel = me:dialog_list(npc, message, {
			"#b전리품을 교환한다.#k",
		})
		if sel == nil then
			return
		end
		if me:dialog_yes_no(npc, "전리품을 교환하기 위해서는 #i3980003 ##b 파티퀘스트코인#k  35개를 모아오셔야합니다!") == false then
			return
		end

		if me:empty_slots(InventoryType.Equipment) < 1 then
			me:dialog(npc, "장비창 한칸을 비우시오.")
			return
		end
		if item_count(me, COIN) < 35 then
			me:dialog(npc, "재료가 부족합니다")
			return
		end

		local reward = nil
		if math.random(0, 1) == 0 then
			reward = {
				item = { [GLOVE] = 1 },
				bonus = { watk = 2, matk = 2 },
			}
		else
			reward = {
				item = { [BELT] = 1 },
				bonus = { str = 30, dex = 30, int = 30, luk = 30, watk = 5, matk = 1 },
			}
		end
		if me:exchange({ item = { [COIN] = 35 } }, reward) ~= ExchangeResult.OK then
			me:dialog(npc, "재료가 부족합니다")
			return
		end
		me:dialog(npc, "전리품을 교환하여 아이템이 지급되었습니다.무슨 아이템이 나왔는지는 장비템을 확인해주세요!")
	end
}
