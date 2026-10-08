-- NPC name (String.wz/Npc.img.xml): 할로 캣

local boss_entry = require("script/lib/boss_entry")

local BOSS = {
	name = "매그너스",
	difficulty = "#k★★#k☆☆☆#b#n",
	maps = { 123356785 },
	time = 30,
	mobs = { 8880000 },
	x = 2365,
	y = -1347,
	limit = 5,
	record = "boss_entry.magnus",
}

return {
	on_click = function(me, npc)
		boss_entry.on_click(me, npc, BOSS)
	end
}
