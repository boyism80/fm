local cpq = require("script/lib/carnival")

return {
	on_enter = function(me)
		cpq.revive(me)
	end
}
