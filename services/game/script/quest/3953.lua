-- 무하마드 설득하기 (Quest.wz/QuestData/3953.img): 무하마드 설득하기

local quest_id = 3953

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

		local sel = me:dialog_list(npc, "데우가 몬스터라는 헛소리를 할 거라면 들을 생각 없으니 어서 저리 비키게! ...응? 흠... 이건 리튬이 아닌가. 색을 보니 최상급 리튬인데... 상태도 훌륭하고... 응? 이걸 주겠다고? 허험... 리튬이라면 사양할 수 없지. 그래... 무슨 일인가?", {
			"데우는 몬스터라는 사실을 꼭 알려드리고 싶습니다.",
			"사막을 이동하던 상단이 몬스터에게 공격받았다는 이야기를 들으셨나요?",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			me:dialog(npc, "데우는 몬스터가 아니라고 하지 않았나! 여기서 썩 꺼지게!", false, false)
			return
		end

		sel = me:dialog_list(npc, "상단이? ...호위가 부족했던 모양이군. 버닝로드에는 크게 위험한 몬스터는 없지만 그래도 방심하면 안 되는 법인데... 사막은 항상 먼저 주의하는 수밖에 없다네.", {
			"데우만 퇴치하면 이런 일은 없을 겁니다.",
			"이게 다 왕비가 마을 주변 치안에 너무 소홀한 탓이에요.",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			me:dialog(npc, "데우에게 모든 죄를 뒤집어 씌우려는 겐가? 더 말할 필요가 없네. 여기서 썩 사라지게.", false, false)
			return
		end

		sel = me:dialog_list(npc, "맞아! 왕비 때문이지! 그 여자가 온 이후로... 총명하던 압둘라 8세께서는 변해 버리고 아리안트는 점점 말라가고만 있어! 마치 오아시스가 마르는 것처럼! 다 그 여자 때문이야!", {
			"빨리 군대를 모아 왕비의 압제에서 벗어나야 해요!",
			"왕비가 이렇게 폭정을 펼치는데, 사막의 수호신은 뭘 하는지 모르겠어요.",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			me:dialog(npc, "쉿, 목소리를 낮추게. 누가 들을라. 그 생각은 아직은 위험한 것 같군. 사막의 수호신이 지금은 뭘 하는지..", false, false)
			return
		end

		sel = me:dialog_list(npc, "...그러게 말일세. 데우가 좀 더 힘을 내주었더라면 좋았을 것을. 수호신께선 야박하시기도 하지...", {
			"데우는 몬스터일 뿐이니, 어쩔 수 없는 일이잖아요?",
			"그래서 말인데... 데우는 이미 몬스터가 되어 버린 거 아닐까요?",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			me:dialog(npc, "사막의 수호신 데우를 몬스터로 의심하지 말게. 흐음.. 아쉽지만 더 할 얘기는 없어 보이는군.", false, false)
			return
		end

		me:dialog(npc, "그래... 자네 말이 맞을지도 몰라. 아리안트가 이리 변하다니... 그건 필시 데우가 변해버린 탓일지도 모르지. 데우는 이미 몬스터가 되어 버린 걸지도... 젊은 사람들 말대로, 이젠 데우를 퇴치해야 할 때인가...", false, false)
		q:force_complete(npc)
		me:show_quest_completion(3953)
		me:show_effect(EffectType.QuestCompletion)
	end
}
