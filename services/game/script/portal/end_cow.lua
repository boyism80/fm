return {
	on_enter = function(me)
		me:records():remove("cow.milk")
		me:records():remove("cow.last")
		me:play_portal_sound()
		me:map(120000103, 1)
	end
}
