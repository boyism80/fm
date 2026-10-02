local mini_dungeon = require("script/lib/mini_dungeon")

return {
	on_enter = function(me)
		mini_dungeon.enter(me, 261020300, 261020301)
	end
}
