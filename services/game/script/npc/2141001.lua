-- NPC name (String.wz/Npc.img.xml): 잊혀진 신전관리인

local ex = require("script/lib/expedition")

return {
	on_click = function(me, npc)
		ex.talk(me, npc, "pink_bean")
	end
}
