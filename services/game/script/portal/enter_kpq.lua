return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			me:notice("지금은 이 포탈을 사용할 수 없습니다.", Msg.PinkText)
			return
		end
		local map = me:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		if wz == nil then
			return
		end
		local map_id = wz.id
		if map_id < 103000800 or map_id > 103000804 then
			return
		end
		if sm:get_property("stage") == "" then
			sm:set_property("stage", "1")
		end
		local stage = tonumber(sm:get_property("stage")) or 1
		local curstage = (map_id % 10) + 1
		if stage <= curstage then
			me:notice("지금은 이 포탈을 사용할 수 없습니다.", Msg.PinkText)
			return
		end
		local portal = map:portal("next00")
		if portal == nil then
			return
		end
		local target_map_id = portal:target_map()
		local target_name = portal:target()
		local group = sm:group()
		if group == nil then
			return
		end
		local dest = group:map(target_map_id)
		if dest == nil then
			return
		end
		local dest_portal = dest:portal(target_name)
		local spawn = 0
		if dest_portal ~= nil then
			spawn = dest_portal:id()
		end
		me:play_portal_sound()
		me:map(dest, spawn)
	end
}
