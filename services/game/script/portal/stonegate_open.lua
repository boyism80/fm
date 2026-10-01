-- Portal (old/scripts/portal/stonegate_open.js): 기사의 홀

local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.pass_gate(me, "stonegate", nil, 990000430)
	end
}
