-- NPC name (String.wz/Npc.img.xml): 마법사 전직교관

local second_class = require("script/lib/second_class")

return {
	on_click = function(me, npc)
		second_class.turn_in_marbles(me, npc, "하인즈", 101010000)
	end
}
