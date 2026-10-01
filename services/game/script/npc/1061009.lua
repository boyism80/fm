-- NPC name (String.wz/Npc.img.xml): 차원의 문

local TRIAL_QUEST = 195000
local BLACK_CHARM = 4031059
local DEFAULT_TEXT = "다른 세계로 통할 것 같은 이상한 모양의 차원의 균열이다."

local CRACKS = {
	[Class.Fighter] = 105070001,
	[Class.Page] = 105070001,
	[Class.Spearman] = 105070001,
	[Class.FpWizard] = 100040106,
	[Class.IlWizard] = 100040106,
	[Class.Cleric] = 100040106,
	[Class.Hunter] = 105040305,
	[Class.Crossbowman] = 105040305,
	[Class.Assassin] = 107000402,
	[Class.Bandit] = 107000402,
	[Class.Brawler] = 105070200,
	[Class.Gunslinger] = 105070200,
}

return {
	on_click = function(me, npc)
		if me:quest(TRIAL_QUEST):record() ~= "job3_trial1_2" or next(me:item(BLACK_CHARM)) ~= nil then
			me:dialog(npc, DEFAULT_TEXT)
			return
		end
		if CRACKS[me:class()] ~= me:map():template_id() then
			me:dialog(npc, DEFAULT_TEXT)
			return
		end
		local sm, err = state_machine("third_job_trial"):start_solo(me)
		if sm == nil then
			log("third_job_trial start_solo:", err)
			me:dialog(npc, "이미 이 균열 안에는 다른 누군가가 들어가 있는 것 같다. 지금은 들어갈 수 없을 것 같다..")
		end
	end
}
