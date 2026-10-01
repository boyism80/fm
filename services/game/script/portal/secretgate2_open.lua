-- Portal (old/scripts/portal/secretgate2_open.js): 수로의 미로

local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.pass_gate(me, "secretgate2", nil, 990000631, 1, "문이 굳게 닫혀있습니다.")
	end
}
