local pq = require("script/lib/party_quest")

local EGG = 4001094

return {
	on_enter = function(me)
		if pq.has_item(me, EGG) == false then
			return
		end
		local nest = me:map():find_reactor_name("dragonBaby")
		if nest ~= nil then
			nest:hit(1)
		end
		me:rmitem(EGG, 1)
		me:message("품안에 있던 나인스피릿의 알이 신비한 빛을 내며 둥지로 돌아갔다.", Msg.PinkText)
	end
}
