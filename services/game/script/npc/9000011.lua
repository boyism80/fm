-- NPC name (String.wz/Npc.img.xml): 마틴

local channel_event = require("script/lib/channel_event")

return {
	on_click = function(me, npc)
		channel_event.guide(me, npc)
	end,
}
