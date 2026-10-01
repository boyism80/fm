-- Portal (old/scripts/portal/guildwaitingenter.js): 유적발굴 현장

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			me:map(101030104)
			return
		end
		if sm:get_property("state") == "waiting" then
			me:message("지금은 포탈이 닫혀있습니다.", Msg.PinkText)
			return
		end
		me:map(990000100)
	end
}
