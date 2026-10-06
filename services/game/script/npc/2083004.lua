-- NPC name (String.wz/Npc.img.xml): 원정대의 표식

local ex = require("script/lib/expedition")

return {
	on_click = function(me, npc)
		ex.talk(me, npc, "horntail")
	end
}
