-- 콜로세움 조사 (Quest.wz/QuestInfo.img.xml): 콜로세움 조사

local quest_id = 31010

return {
	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local data = q:record()
		if data == nil or data == "" then
			me:dialog(npc, "무슨일이시죠? 콜로세움 내부에서는 어떤 적들이 기다리고있는지 확인 해봐야 할텐데...", false, false)
			return
		end

		if not me:dialog(npc, "적의 상황을 알아야 승리할 수 있을텐데...", false, true) then
			return
		end
		if not me:dialog(npc, "네? 콜로세움까지 가셔서 크세르크세스 일당의 병력을 확인하고 오셨따고요? 후아... 역시 대단 하시네요. 그럼 몇가지 질문을 드려도 될까요?", true, true) then
			return
		end

		local sel = me:dialog_list(npc, "콜로세움 내부에는 어떤 몬스터가 있었나요?#b\r\n", {
			"저빌",
			"스코피",
			"페넬",
			"헤라크",
			"맘무트",
		})
		if sel == nil then
			return
		end
		if sel < 4 then
			me:dialog(npc, "정말 다녀오신 게 맞으신가요? 아닌 것 같은데...", false, false)
			return
		end

		sel = me:dialog_list(npc, "맘무트...맘무트라... 굉장한 힘을 가진 녀석들인데...그럼 어떻게 대비를 하는게 좋을까요?#b\r\n", {
			"나도 어떻게 해야할 지 모르겠어.",
			"공격에 큰 피해를 입지 않도록 방어구에 신경을 써야 할 것 같아.",
			"대비 같은건 필요 없어. 나 혼자서도 충분하다고.",
		})
		if sel == nil then
			return
		end
		if sel == 0 then
			me:dialog(npc, "그렇게 막무가내로 전투에 임했다가는 큰 일 날거에요.", false, false)
			return
		end
		if sel == 2 then
			me:dialog(npc, "혼자서는 위험합니다. 저희와 같이 준비하시는 게 좋을겁니다.", false, false)
			return
		end

		local code = me:exchange({}, { exp = 4500 })
		if code ~= ExchangeResult.OK then
			return
		end
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
