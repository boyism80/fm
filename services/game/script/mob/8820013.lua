-- Mob name (String.wz/Mob.img.xml): 핑크빈

local pinkbean = require("script/lib/pinkbean")

return {
	on_revive = function(mob, map, x, y, revives)
		pinkbean.advance_sponge(map, x, y, 8820014, {
			{ marker = 8820024, part = 8820015 },
			{ marker = 8820025, part = 8820016 },
			{ marker = 8820026, part = 8820017 },
			{ marker = 8820027, part = 8820018 },
			{ marker = 8820019, part = 8820002 },
		})
	end
}
