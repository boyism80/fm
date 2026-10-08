-- NPC name (String.wz/Npc.img.xml): 할로 캣

local boss_entry = require("script/lib/boss_entry")

local BOSS = {
	name = "데미안",
	difficulty = "#k★★★★#k☆#b#n",
	maps = { 123356786, 123356787 },
	time = 60,
	mobs = { 9300890 },
	x = 1222,
	y = 16,
	limit = 3,
	record = "boss_entry.damien",
}

return {
	on_click = function(me, npc)
		boss_entry.on_click(me, npc, BOSS)
	end
}
