-- Portal (old/scripts/portal/statuegate_open.js): 샤레니안 성문

local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.pass_gate(me, "statuegate", "stage1clear", 990000301)
	end
}
