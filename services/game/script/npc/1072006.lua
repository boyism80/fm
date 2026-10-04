-- NPC name (String.wz/Npc.img.xml): 궁수 전직교관

local second_job = require("script/lib/second_job")

return {
	on_click = function(me, npc)
		second_job.turn_in_marbles(me, npc, "헬레나", 106010000)
	end
}
