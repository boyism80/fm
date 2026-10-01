-- Portal (old/scripts/portal/watergate_open.js): 현자의 분수

local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.pass_gate(me, "watergate", "stage3clear", 990000600)
	end
}
