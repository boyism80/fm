return {
	on_enter = function(me)
		me:quest(126640):record("0")
		me:quest(126641):record("0")
		me:play_portal_sound()
		me:map(120000103, 1)
	end
}
