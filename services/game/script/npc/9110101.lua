-- NPC name (String.wz/Npc.img.xml): 티옌데

local COIN = 4009097
local CANCEL = "길가다가 넘어져서 한강으로 굴러떨어져라."

local ITEMS = { 1142472, 1302024, 1302019, 1322012, 1332024, 1402017, 1442025, 1452018 }
local COSTS = { 10000, 3000, 2000, 1000, 1000, 1000, 1000, 3500 }

return {
	on_click = function(me, npc)
		local menu = me:dialog_list(npc, "특별한 아이템을 교환하기위해선 특정 개수의 벚꽃 코인이 필요합니다.", { "벚꽃 코인을 사용한다#k" })
		if menu == nil then
			me:dialog(npc, CANCEL)
			return
		end

		local options = {}
		for i, item in ipairs(ITEMS) do
			options[i] = "#i" .. COIN .. ":# " .. COSTS[i] .. "개#k   =  #i" .. item .. ":# #b#z" .. item .. ":##k"
		end
		local sel = me:dialog_list(npc, "어서오세요~#r미리보는 더미 아이템의 옵션은 실 옵션과 다릅니다.#k구매 해보시면 후회 없어요~", options)
		if sel == nil then
			me:dialog(npc, CANCEL)
			return
		end

		local code = me:exchange({ item = { [COIN] = COSTS[sel] } }, { item = { [ITEMS[sel]] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "코인이 부족하거나 인벤토리 빈 공간이 부족합니다.")
			return
		end
		me:dialog(npc, "아이템을 교환하였습니다. 즐거운 메이플스토리 되세요~")
	end
}
