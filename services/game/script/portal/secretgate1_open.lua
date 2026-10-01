-- Portal (old/scripts/portal/secretgate1_open.js): 수로의 미로

local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.pass_gate(me, "secretgate1", nil, 990000611, 1, "문이 굳게 닫혀있습니다.")
	end
}
