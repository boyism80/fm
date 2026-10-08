local dojo = require("script/lib/dojo")

return {
	on_map_enter = function(me, map)
		dojo.taunt(map)
	end
}
