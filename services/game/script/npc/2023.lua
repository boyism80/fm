-- NPC name (String.wz/Npc.img.xml): 할로 캣

local boss_entry = require("script/lib/boss_entry")

local BOSS = {
	name = "스우",
	difficulty = "#k★#k☆☆☆☆#b#n",
	maps = { 123356782, 123356783, 123356784 },
	time = 40,
	mobs = { 9801028 },
	x = -110,
	y = -16,
	limit = 5,
	quest = 19022130,
	day_quest = 19022131,
}

return {
	on_click = function(me, npc)
		boss_entry.on_click(me, npc, BOSS)
	end
}
