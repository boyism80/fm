local ARENA_PORTAL = 1

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		local arena_id = me:map():template_id() + 1
		for _, map in ipairs(sm:maps()) do
			if map:template_id() == arena_id then
				me:play_portal_sound()
				me:map(map, ARENA_PORTAL)
				return
			end
		end
	end,
}
