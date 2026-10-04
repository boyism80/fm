-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")

return {
	on_revive = function(mob, map, x, y, revives)
		pinkbean.advance_sponge(map, x, y, 8820013, {
			{ marker = 8820024, part = 8820003 },
			{ marker = 8820025, part = 8820004 },
			{ marker = 8820026, part = 8820005 },
			{ marker = 8820023, part = 8820006 },
		})
	end
}
