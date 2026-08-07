local function relocate(me, portal_name)
	local map = me:map()
	if map == nil then
		return
	end
	local portal = map:portal(portal_name)
	if portal == nil then
		return
	end
	me:play_portal_sound()
	me:map(map, portal:id(), true)
end

return {
	on_enter = function(me, portal)
		local sm = me:state_machine()
		if sm == nil or portal == nil then
			return
		end
		if sm:get_property("stage4_rand") == "" then
			sm:set_property("stage4_rand", "1")
			sm:set_property("stage4_r4way1", tostring(math.random(1, 3)))
			sm:set_property("stage4_r4way2", tostring(math.random(1, 3)))
		end
		local pid = portal:id()
		local way1 = sm:get_property("stage4_r4way1")
		local way2 = sm:get_property("stage4_r4way2")
		if pid == 11 then
			if way1 == "1" then
				relocate(me, "np00")
			else
				relocate(me, "np02")
			end
		elseif pid == 12 then
			if way1 == "2" then
				relocate(me, "np00")
			else
				relocate(me, "np02")
			end
		elseif pid == 13 then
			if way1 == "3" then
				relocate(me, "np00")
			else
				relocate(me, "np02")
			end
		elseif pid == 14 then
			if way2 == "1" then
				relocate(me, "np01")
			else
				relocate(me, "np02")
			end
		elseif pid == 15 then
			if way2 == "2" then
				relocate(me, "np01")
			else
				relocate(me, "np02")
			end
		elseif pid == 16 then
			if way2 == "3" then
				relocate(me, "np01")
			else
				relocate(me, "np02")
			end
		end
	end
}
