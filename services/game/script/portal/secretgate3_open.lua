-- Portal (old/scripts/portal/secretgate3_open.js): 수로의 미로

local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.pass_gate(me, "secretgate3", nil, 990000641, 1, "문이 굳게 닫혀있습니다.")
	end
}
