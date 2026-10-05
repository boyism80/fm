-- NPC name (String.wz/Npc.img.xml): 궁수 전직교관

local second_class = require("script/lib/second_class")

return {
	on_click = function(me, npc)
		second_class.turn_in_marbles(me, npc, "헬레나", 106010000)
	end
}
