-- Portal (old/scripts/portal/metalgate_open.js): 신념의 방

local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.pass_gate(me, "metalgate", nil, 990000431)
	end
}
