-- NPC name (String.wz/Npc.img.xml): 할로 캣

local boss_entry = require("script/lib/boss_entry")

local BOSS = {
	name = "검은 마법사",
	difficulty = "#r★★★★★(헬)#b#n",
	maps = { 123356788, 123356789 },
	time = 120,
	mobs = { 8880500 },
	x = -3,
	y = 85,
	limit = 3,
	quest = 19022150,
	day_quest = 19022151,
}

return {
	on_click = function(me, npc)
		boss_entry.on_click(me, npc, BOSS)
	end
}
