-- NPC name (String.wz/Npc.img.xml): 돌고래

local TICKET = 4031242

return {
	on_click = function(me, npc)
		local meso = 10000
		local meso2 = 1000
		if me:class() == Class.Beginner then
			meso = 1000
			meso2 = 100
		end

		local tickets = 0
		for _, it in pairs(me:item(TICKET)) do
			tickets = tickets + it:count()
		end
		local first = " " .. meso2 .. "메소를 내고 #m230030200#로 갑니다."
		if tickets >= 1 then
			first = " #t4031242#으로 #m230030200#로 갑니다."
		end
		local sel = me:dialog_list(npc, "안녕하세요~ 돌고래 택시입니다! 안전하고 편안하게 사각지대까지 모셔다 드립니다. 초보자는 요금을 90% 할인해 드립니다.\r\n\r\n#b", {
			first,
			" " .. meso .. "메소를 내고 #m251000100#까지 갑니다.",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			local cost = { meso = meso2 }
			if tickets >= 1 then
				cost = { item = { [TICKET] = 1 } }
			end
			local code = me:exchange(cost, nil)
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "티켓이나 메소를 분명 제대로 갖고 계신건지 확인해 보세요.")
				return
			end
			me:map(230030200, 3)
			return
		end

		local code = me:exchange({ meso = meso }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "메소가 부족하신건 아닌지 확인해 보세요.")
			return
		end
		me:map(251000100)
	end
}
