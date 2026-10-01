-- Portal (old/scripts/portal/ghostgate_open.js): 샤렌 3세의 무덤

local gq = require("script/lib/guild_quest")

return {
	on_enter = function(me)
		gq.pass_gate(me, "ghostgate", "stage4clear", 990000800, 0, "결계로 가로막혀 있어 지나갈 수 없습니다.")
	end
}
