local quest_id = 2148

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	local sel = me:dialog_list(npc, "무슨 일이지?", {
		"귀신나무에 대해 들어보신 적이 있나요?",
	})
	if sel == nil then
		return
	end

	sel = me:dialog_list(npc, "귀신나무? 아, 오래 전에 사라졌던 그 거대한 스텀프를 말하는 건가? 아버지의 아버지가 어릴 적에 그런 나무가 있었다는 이야기를 들은 적이 있었다네. 전해져 오는 소문에는 가지마다 붉은 천이 달려있는데 혼령의 피로 물들은거라고 하더군. 하지만 나도 실제로 본 적은 한번도 없다네. 그러니 진실인지는 알 수 없지.", {
		"다른 소문은 듣지 못했나요?",
	})
	if sel == nil then
		return
	end

	me:dialog(npc, "애석하게도 나는 소문에 밝은 사람이 아니라네.", false, true)

	q:start(npc, true)
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
	me:exp(me:exp() + 100)
end
