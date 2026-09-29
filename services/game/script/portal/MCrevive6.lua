return {
	on_enter = function(me)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local match = carnival.map_match(map:wz():id())
		if match == nil then
			return
		end
		local portal = "sp"
		local team = me:carnival_team()
		if team ~= nil then
			local team_id = team:id()
			if team_id == CARNIVAL_TEAM.RED then
				portal = "red_revive"
			elseif team_id == CARNIVAL_TEAM.BLUE then
				portal = "blue_revive"
			end
		end
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		local field = sm:group():map(match:field_map_id())
		if field == nil then
			return
		end
		local target = field:portal(portal)
		if target == nil then
			return
		end
		me:map(field, target:id())
	end
}
