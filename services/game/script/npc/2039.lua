-- NPC name (String.wz/Npc.img.xml): 할로 캣

local boss_entry = require("script/lib/boss_entry")

local BOSS = {
	name = "루시드",
	difficulty = "#k★★★★#k☆#b#n",
	maps = { 450004150, 450004550 },
	time = 60,
	mobs = { 8880166, 8880140 },
	x = 1008,
	y = 48,
	limit = 4,
	record = "boss_entry.lucid",
}

return {
	on_click = function(me, npc)
		boss_entry.on_click(me, npc, BOSS)
	end
}
