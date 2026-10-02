local mini_dungeon = require("script/lib/mini_dungeon")

return {
	on_enter = function(me)
		mini_dungeon.enter(me, 260020600, 260020630)
	end
}
