-- NPC name (String.wz/Npc.img.xml): 돌고래

return {
	on_click = function(me, npc)
		local meso = 10000
		if me:class() == Class.Beginner then
			meso = 1000
		end

		local sel = me:dialog_list(npc, "안녕하세요~ 돌고래 택시입니다! 바닷길은 모두 이어져 있기에 어디로든지 갈 수 있답니다. 초보자는 요금을 90% 할인해 드립니다.\r\n\r\n#b", {
			" " .. meso .. "메소를 내고 #m230000000#까지 갑니다.",
		})
		if sel == nil then
			return
		end

		local code = me:exchange({ meso = meso }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "메소가 부족하신건 아닌지 확인해 보세요.")
			return
		end
		me:map(230000000)
	end
}
