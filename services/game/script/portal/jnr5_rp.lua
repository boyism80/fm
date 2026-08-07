local function relocate(me, name)
	local map = me:map()
	if map == nil then
		return
	end
	local portal = map:portal(name)
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
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local room = (map:wz():id() % 10) - 1
		local pname = portal:name()
		if #pname < 4 then
			return
		end
		local b = pname:sub(3, 3)
		local c = pname:sub(4, 4)
		local key = string.format("stage6_%d_%s_%s", room, b, c)
		if sm:get_property(key) == "1" then
			relocate(me, "np0" .. b)
		else
			relocate(me, "npFail")
		end
	end
}
