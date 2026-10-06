-- NPC name (String.wz/Npc.img.xml): 아도비스

local ex = require("script/lib/expedition")

return {
	on_click = function(me, npc)
		ex.talk(me, npc, "zakum")
	end
}
