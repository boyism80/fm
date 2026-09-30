-- NPC name (String.wz/Npc.img.xml): 낙타 택시

return {
	on_click = function(me, npc)
		local field = me:map()
		local wz = field ~= nil and field:wz() or nil
		if wz == nil then
			return
		end

		local here = wz:id()
		local dest
		local text
		if here == 260020000 then
			dest = 260020700
			text = "#b5000메소#k를 지불하고 사헬지대1으로 가시겠습니까?"
		elseif here == 260020700 then
			dest = 260020000
			text = "#b5000메소#k를 지불하고 아리안트 북문 밖으로 가시겠습니까?"
		else
			return
		end
		if not me:dialog_yes_no(npc, text) then
			return
		end

		local code = me:exchange({ meso = 5000 }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "잠깐... 메소가 부족하신데요!")
			return
		end
		me:map(dest)
	end
}
