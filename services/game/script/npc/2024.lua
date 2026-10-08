-- NPC name (String.wz/Npc.img.xml): 할로 캣

local boss_entry = require("script/lib/boss_entry")

local BOSS = {
	name = "카오스 파풀라투스",
	difficulty = "#k★★★#k☆☆#b#n",
	maps = { 123356781 },
	time = 50,
	mobs = { 8820119 },
	x = -393,
	y = -386,
	limit = 10,
	record = "boss_entry.chaos_papulatus",
}

return {
	on_click = function(me, npc)
		boss_entry.on_click(me, npc, BOSS)
	end
}
