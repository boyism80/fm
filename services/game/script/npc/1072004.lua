-- NPC name (String.wz/Npc.img.xml): 전사 전직교관

local second_job = require("script/lib/second_job")

return {
	on_click = function(me, npc)
		second_job.turn_in_marbles(me, npc, "주먹펴고일어서", 102020300)
	end
}
