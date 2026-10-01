return {
	on_enter = function(me, portal)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		for _, map in ipairs(sm:maps()) do
			if map:template_id() == portal:target_map() then
				local target = map:portal(portal:target())
				local spawn = 0
				if target ~= nil then
					spawn = target:id()
				end
				me:play_portal_sound()
				me:map(map, spawn)
				return
			end
		end
	end,
}
