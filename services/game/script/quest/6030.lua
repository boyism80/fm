-- 카슨의 연금술 수업 (Quest.wz/QuestData/6030.img.xml): 카슨의 연금술 수업

local quest_id = 6030

return {
	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not q:started() and not q:completed() then
			q:start(npc, true)
			return
		end

		if not me:dialog(npc, "어서 오게. 자네를 보낸다는 스탠의 연락은 이미 받았네. 기초 연금술에 대한 수업을 받으면 된다고 했던가?", false, true) then
			return
		end
		if not me:dialog(npc, "연금술의 기본을 정의하자면 #b순환#k과 #b교환#k이라고 할 수 있네.우선 순환이란 물질을 구성하는 성질이 서로 일정한 규칙을 가지고 변화한다는 것을 말하지. 물은 나무로, 나무는 불로, 불은 흙으로, 흙은 금으로, 금은 다시 물로 변하는 것처럼 말이지. 이것은 연금술의 가장 기초이면서 세상 모든 물질의 기본적인 성질이라네.", false, true) then
			return
		end
		if not me:dialog(npc, "교환은 순환과는 조금 다르지. 순환이 정해진 범위 안에서 물질의 성질을 바꾸는 힘이라면 교환은 물질의 절대적인 양에 대한 개념이라네. 무에서 유를 만들 수 없듯이 연금술을 통해서 전에 없던 새로운 것을 만들 수는 없다네. 항상 그곳에 존재했던 것 존재했던 물질을 순환법칙으로 가공하여 새로운 것으로 교환하는 것 그것이 연금술의 핵심이지.", false, true) then
			return
		end
		if not me:dialog_yes_no(npc, "자네의 표정을 보아하니 과연 내 말을 잘 알아들었는지 알 수가 없구만. 하긴 연금술을 그렇게 단시간 안에 모두 이해할 수 있다면 자네는 천년에 한번 나올까 말까한 인재일테니까. 어쨌든 수업을 마치도록 하겠네. 스탠에게는 수업을 들었다고 연락해 두지.") then
			me:dialog(npc, "수업을 다시 듣고 싶은건가? 그렇다면 내게 다시 말을 걸게나.", false, false)
			return
		end

		local qr = me:quest(6029)
		if qr ~= nil then
			local info = qr:record()
			if info == nil or info == "" then
				info = "000"
			end
			qr:record("1" .. info:sub(2, 3))
		end
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
