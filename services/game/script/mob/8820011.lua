-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")

return {
	on_revive = function(mob, map, x, y, revives)
		pinkbean.advance_sponge(map, x, y, 8820012, {
			{ marker = 8820024, part = 8820003 },
			{ marker = 8820025, part = 8820004 },
			{ marker = 8820022, part = 8820005 },
		})
	end
}
