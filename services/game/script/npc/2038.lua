-- NPC name (String.wz/Npc.img.xml): 할로 캣

local boss_entry = require("script/lib/boss_entry")

local BOSS = {
	name = "힐라",
	difficulty = "#k★#k☆☆☆☆#b#n",
	maps = { 123356780 },
	time = 40,
	mobs = { 8870000 },
	x = 164,
	y = 196,
	limit = 5,
	quest = 19022100,
	day_quest = 19022101,
}

return {
	on_click = function(me, npc)
		boss_entry.on_click(me, npc, BOSS)
	end
}
