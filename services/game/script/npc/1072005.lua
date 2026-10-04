-- NPC name (String.wz/Npc.img.xml): 마법사 전직교관

local second_job = require("script/lib/second_job")

return {
	on_click = function(me, npc)
		second_job.turn_in_marbles(me, npc, "하인즈", 101010000)
	end
}
