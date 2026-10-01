-- Portal (old/scripts/portal/kinggate_open.js): 왕의 회랑

local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.pass_gate(me, "kinggate", nil, 990000900, 1, "문이 굳게 닫혀있습니다.")
	end
}
