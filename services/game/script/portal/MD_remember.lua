local mini_dungeon = require("script/lib/mini_dungeon")

return {
	on_enter = function(me)
		mini_dungeon.enter(me, 240040511, 240040800)
	end
}
