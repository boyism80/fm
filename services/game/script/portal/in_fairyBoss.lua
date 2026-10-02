local SCRIPT = "script/portal/in_fairyBoss.lua"
local ROOM = 300030310

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
		local ok, count = run_on_map(ROOM, SCRIPT, "character_count")
		if me:quest(31229):started() == false or ok == false or count ~= 0 then
			me:message("이미 누군가 안에 입장해있거나 퀘스트를 진행중이시지 않습니다.", Msg.PinkText)
			return
		end
		run_on_map(ROOM, SCRIPT, "reset_map")
		me:map(ROOM)
	end
}
