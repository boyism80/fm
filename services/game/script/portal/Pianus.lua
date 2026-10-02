local SCRIPT = "script/portal/Pianus.lua"
local ROOM = 230040420

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
		if ok and count > 10 then
			me:message("이 방은 이미 피아누스와의 전투를 위한 최대 인원수 만큼 가득 찼습니다.")
			return
		end
		me:play_portal_sound()
		me:map(ROOM)
	end
}
