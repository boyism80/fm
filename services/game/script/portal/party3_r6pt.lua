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
		local way = sm:get_property("stage6_way")
		if way == "" then
			local digits = {}
			for i = 1, 16 do
				digits[i] = tostring(math.random(1, 4))
			end
			way = table.concat(digits)
			sm:set_property("stage6_way", way)
		end
		local pid = portal:id()
		local t1 = pid - 24
		local wayfloor = math.floor(t1 / 4)
		local wayans = way:sub(wayfloor + 1, wayfloor + 1)
		local portalindex = tostring(((wayfloor + 1) * 4 + 24) - pid)
		if wayans ~= portalindex then
			if wayfloor <= 4 then
				relocate(me, "np16")
			elseif wayfloor <= 8 then
				relocate(me, "np03")
			elseif wayfloor <= 12 then
				relocate(me, "np07")
			else
				relocate(me, "np11")
			end
		else
			if wayfloor >= 10 then
				relocate(me, "np" .. tostring(wayfloor))
			else
				relocate(me, string.format("np%02d", wayfloor))
			end
		end
	end
}
