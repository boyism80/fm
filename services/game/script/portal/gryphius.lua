local SCRIPT = "script/portal/gryphius.lua"
local ROOM = 240020101

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	on_enter = function(me)
		local ok, count = run_on_map(ROOM, SCRIPT, "character_count")
		if ok and count >= 6 then
			me:message("그리프의 숲에 입장할 수 없습니다.")
			return
		end
		me:play_portal_sound()
		me:map(ROOM)
	end
}
