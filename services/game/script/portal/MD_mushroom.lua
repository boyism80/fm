local mini_dungeon = require("script/lib/mini_dungeon")

return {
	on_enter = function(me)
		mini_dungeon.enter(me, 105050100, 105050101)
	end
}
