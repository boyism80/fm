-- NPC name (String.wz/Npc.img.xml): 유레테

local rj = require("script/lib/romeo_juliet")

local QUEST_ID = 3382

return {
	on_click = function(me, npc)
		local q = me:quest(QUEST_ID)
		if q == nil then
			return
		end
		if not q:started() then
			q:start(npc, true)
			return
		end
		if rj.exchange_marbles(me, npc) then
			q:force_complete(npc)
		end
	end
}
