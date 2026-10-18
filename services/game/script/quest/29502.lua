-- 칭호 도전 - 핑크빈 슬레이어 (Quest.wz/QuestData/29502.img): 칭호 도전 - 핑크빈 슬레이어

local medal = require("script/lib/medal")

local quest_id = 29502

return {
	on_start = function(me, npc)
		medal.start_slayer(me, npc, quest_id, medal.PINKBEAN_SLAYER, "핑크빈")
	end,

	on_end = function(me, npc)
		medal.end_slayer(me, npc, quest_id, medal.PINKBEAN_SLAYER)
	end,
}
