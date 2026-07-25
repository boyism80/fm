local quest_id = 3321

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "...어째서 이런 곳에 오신 건지 모르지만... 연금술사의 실험실은 그리 즐거운 곳이 아닙니다. 연금술사가 아닌 사람의 눈에는 무척 지루하다더군요. 하긴... 그녀야 요정이니 더 재미없어 보일지도 모르겠군요...  ", false, true)
		me:dialog(npc, "그녀가 누구냐고요? 그녀는... 제 아내입니다. 그러고 보니 그녀의 얼굴을 못 본지도 꽤 오래 되었군요... 딸 아이의 얼굴이 가물가물할 정도이니... 그녀가 무척 화를 내겠군요. 물론 상냥한 그녀는 곧 용서해 줄 테지만요... ", false, true)
		me:dialog(npc, "...하지만 어쩔 수 없지요. 이 연구를 마치기 전까지 그녀의 얼굴을 보지 않겠다고 결심했으니까요. 무척 보고 싶지만... 연구를 마치기 전까지는... 연구만 마치면 영원히 #b#p2111004##k의 얼굴을 볼 수 있을 겁니다.", false, true)
		me:dialog(npc, "그러고 보니 #b펜던트#k를 아직도 그녀에게 선물하지 못했군요. 그녀에게 들킬까봐 #b액자 뒤#k에 숨겨 놓기까지 했었는데... 그녀의 얼굴을 볼 수 없으니 선물조차 할 수 없네요. 언제쯤이면 그녀를 볼 수 있을까요... ", false, true)
		if not me:dialog_accept(npc, "...쓸데없는 이야기가 너무 길어졌군요. 죄송합니다만, 연구를 계속 해야 해서... 그만 이 연구실에서 나가주십시오.") then
			me:dialog(npc, "무례하신 분이군요...", false, false)
			return
		end

		q:start(npc, true)
		me:map(261020401, 0)
	end
}
