-- NPC name (String.wz/Npc.img.xml): 도적 전직교관

local second_job = require("script/lib/second_job")

return {
	on_click = function(me, npc)
		second_job.turn_in_marbles(me, npc, "다크로드", 102040000)
	end
}
