-- NPC name (String.wz/Npc.img.xml): 조수 레드 - 몬스터 카니발

local cpq = require("script/lib/carnival")

return {
	on_click = function(me, npc)
		cpq.leave(me)
	end
}
