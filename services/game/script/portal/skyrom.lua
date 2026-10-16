local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/skyrom.lua"
local ROOM = 926000010

return {

	prepare_room = function(map)
		map:reset()
		pq.shuffle_reactors(map)
	end,

	on_enter = function(me)
		if me:quest(3935):started() == false or pq.has_item(me, 4031574) then
			return
		end
		local ok, count = run_on_map(ROOM, "script/lib/party_quest.lua", "character_count")
		if ok == false or count ~= 0 then
			me:message("이미 이 안에 다른 누군가가 들어가서 스카이롬을 훔치는 중입니다.")
			return
		end
		run_on_map(ROOM, SCRIPT, "prepare_room")
		me:play_portal_sound()
		me:map(ROOM)
	end
}
