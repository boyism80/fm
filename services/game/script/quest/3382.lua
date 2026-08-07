-- Quest name (Quest.wz/Quest.img.xml): 유레테의 보답

local rj = require("script/lib/romeo_juliet")

local quest_id = 3382

return {
	on_end = function(me, npc)
		local q = me:quest(quest_id)
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
