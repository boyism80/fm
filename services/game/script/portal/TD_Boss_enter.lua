local SCRIPT = "script/portal/TD_Boss_enter.lua"
local ROOMS = 6

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	reset_map = function(map)
		map:reset()
	end,

	on_enter = function(me)
		local map_id = me:map():wz():id()
		for i = 1, ROOMS do
			local ok, count = run_on_map(map_id + i, SCRIPT, "character_count")
			if ok and count == 0 then
				run_on_map(map_id + i, SCRIPT, "reset_map")
				me:play_portal_sound()
				me:map(map_id + i, 1)
				return
			end
		end
		me:message("이미 다른 누군가가 들어가 있는 것 같다. 지금은 들어갈 수 없을 것 같다.", Msg.PinkText)
	end
}
