-- 휴즈의 과학 수업 (Quest.wz/QuestData/6031.img.xml): 휴즈의 과학 수업

local quest_id = 6031

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not q:started() and not q:completed() then
		q:start(npc, true)
		return
	end

	if not me:dialog(npc, "어서 와. 스탠에게 이야기는 들었어. 귀찮게 말이지. 나한테 왜 이런 걸 부탁하는건지 모르겠네. 과학이라는 게 뭐라고 생각해? 난 한 번도 생각해 본 적이 없어. 그런 나에게 과학의 기본을 설명하라니 그거 좀 웃기지 않아?", false, true) then
		return
	end
	if not me:dialog(npc, "과학이란 어려운 게 아니라구. 생활 속에 녹아 있단 말이지. 여기 이 기계를 좀 봐. 복잡해 보여? 어려워 보이냐고! 아냐 하나도 다르지 않아! 마치 내가 정상인 것처럼 말이지. 하지만 사람들은 지레 겁을 먹지. 선입견을 가지면서 말이야. 그건 좋지 않아 좋지 않다고!!", false, true) then
		return
	end
	if not me:dialog(npc, "잘~ 생각해봐야 해. 중요한 건 이해한다는 거야. 이해하려고 노력하는 거지. 자신이 이해하는 것! 그것을 넘어서는 것! 그리고 그걸 다시 이해하는 것! 그게 바로 핵심이라고... 뭔지 알아? 알아 듣겠어?", false, true) then
		return
	end
	if not me:dialog_yes_no(npc, "물리적인 법칙? 수식? 실험? 그런 건 다 어중이떠중이야. 중요한건 그게 아니라고! 아...오랫만에 말을 많이 했더니 피곤하군. 이만 가봐. 수업은 끝났어.") then
		me:dialog(npc, "젠장.. 피곤하다니까. 뭐 그래도 수업을 다시 듣고 싶다면 특별히 한번쯤은 더 들려주도록 하지.", false, false)
		return
	end

	local qr = me:quest(6029)
	if qr ~= nil then
		local info = qr:record()
		if info == nil or info == "" then
			info = "000"
		end
		qr:record(info:sub(1, 1) .. "1" .. info:sub(3, 3))
	end
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
end
