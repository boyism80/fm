-- NPC name (String.wz/Npc.img.xml): 에릭 무사

local COIN = 3980019
local CANCEL = "길가다가 넘어져서 한강으로 굴러떨어져라."

local ITEMS = { 1003535, 1003534 }
local COSTS = { 3, 1 }

return {
	on_click = function(me, npc)
		local menu = me:dialog_list(npc, "환상의 나라, 말레이시아에 오신것을 환영합니다!", { "스카 코인 으로 기념품을 살래요!#k" })
		if menu == nil then
			me:dialog(npc, CANCEL)
			return
		end

		local options = {}
		for i, item in ipairs(ITEMS) do
			options[i] = "#i" .. COIN .. ":# " .. COSTS[i] .. "개#k   =  #i" .. item .. ":# #b#z" .. item .. ":##k"
		end
		local sel = me:dialog_list(npc, "어서오세요~ 지역 특산 기념품을 판매중입니다!#k구매 해보시면 후회 없어요~", options)
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
