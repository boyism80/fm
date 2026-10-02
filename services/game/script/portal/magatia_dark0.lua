return {
	on_enter = function(me)
		me:play_portal_sound()
		if me:quest(3309):started() then
			me:map(926120000, 1)
		else
			me:map(261020700, 4)
		end
	end
}
