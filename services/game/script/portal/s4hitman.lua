local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/s4hitman.lua"
local ROOM = 910200000
local EXIT = 101030104

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
		if me:quest(6201):started() == false then
			me:message("알 수 없는 힘으로 봉인되어 있습니다.")
			return
		end
		if pq.has_item(me, 4031452) then
			me:message("샨이 요청한 것을 끝냈습니다. 더 이상 입장할 필요가 없습니다.")
			return
		end
		local ok, count = run_on_map(ROOM, SCRIPT, "character_count")
		if ok == false or count ~= 0 then
			me:message("이미 다른 누군가가 퀘스트에 도전중입니다.")
			return
		end
		me:play_portal_sound()
		run_on_map(ROOM, SCRIPT, "reset_map")
		me:map(ROOM, 0)
		me:warp_later(EXIT, 1200)
	end
}
