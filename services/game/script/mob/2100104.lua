-- Mob name (String.wz/Mob.img.xml): 로얄 카투스

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 35, max = 55, chance = 500000 },
	{ item = 2000001, min = 1, max = 1, chance = 10000 },
	{ item = 2000003, min = 1, max = 1, chance = 10000 },
	{ item = 2060000, min = 20, max = 30, chance = 8000 },
	{ item = 2061000, min = 20, max = 30, chance = 8000 },
	{ item = 4000331, min = 1, max = 1, chance = 400000 },
	{ item = 2002003, min = 1, max = 1, chance = 10000 },
	{ item = 4020006, min = 1, max = 1, chance = 150 },
	{ item = 4020005, min = 1, max = 1, chance = 150 },
	{ item = 4004001, min = 1, max = 1, chance = 80 },
	{ item = 2041002, min = 1, max = 1, chance = 70 },
	{ item = 1442001, min = 1, max = 1, chance = 40 },
	{ item = 1322003, min = 1, max = 1, chance = 40 },
	{ item = 1092019, min = 1, max = 1, chance = 40 },
	{ item = 1072087, min = 1, max = 1, chance = 40 },
	{ item = 1040068, min = 1, max = 1, chance = 40 },
	{ item = 1060057, min = 1, max = 1, chance = 40 },
	{ item = 1082052, min = 1, max = 1, chance = 40 },
	{ item = 1002004, min = 1, max = 1, chance = 40 },
	{ item = 2022155, min = 1, max = 1, chance = 6000 },
	{ item = 4010007, min = 1, max = 1, chance = 1000 },
	{ item = 1082186, min = 1, max = 1, chance = 40 },
	{ item = 2043214, min = 1, max = 1, chance = 70 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
