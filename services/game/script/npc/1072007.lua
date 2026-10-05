-- NPC name (String.wz/Npc.img.xml): 도적 전직교관

local second_class = require("script/lib/second_class")

return {
	on_click = function(me, npc)
		second_class.turn_in_marbles(me, npc, "다크로드", 102040000)
	end
}
