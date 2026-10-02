local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/s4berserk_move.lua"
local ALTAR = 910500200

return {
	mob_count = function(map)
		return pq.mob_count(map)
	end,

	on_enter = function(me)
		local ok, count = run_on_map(ALTAR, SCRIPT, "mob_count")
		if ok == false or count > 0 then
			me:message("포탈이 봉인되어 있습니다.")
			return
		end
		me:play_portal_sound()
		me:map(ALTAR, "pt00")
	end
}
