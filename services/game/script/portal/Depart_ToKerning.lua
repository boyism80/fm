return {
	on_enter = function(me)
		me:play_portal_sound()
		local function on_arrive(me)
			me:clock(20, function(me)
				me:map(103000100)
			end)
		end
		me:map(103000301, 0, { callback = on_arrive })
	end
}
