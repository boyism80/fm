local free_market = require("script/lib/free_market")

return {
	on_enter = function(me)
		free_market.enter(me)
	end
}
