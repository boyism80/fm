local ALTAR = 2408003

return {
	on_enter = function(me)
		local map = me:map()
		local altar = map:reactor(ALTAR)
		if altar == nil or altar:state() > 0 then
			return
		end
		me:message("깊은 동굴 속에서 거대한 생물체가 다가오고 있습니다.", Msg.LightBlueText, MessageScope.Map)
		local x, y = altar:position()
		map:spawn_mob(8810025, x, y)
		altar:hit(0)
		altar:hit(1)
	end
}
