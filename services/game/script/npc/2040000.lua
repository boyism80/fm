-- NPC name (String.wz/Npc.img.xml): 멜

local trains = require("script/lib/trains")

return {
	on_click = function(me, npc)
		trains.sell_train_ticket(me, npc)
	end,
}
