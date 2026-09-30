-- NPC name (String.wz/Npc.img.xml): 마이크

return {
	on_click = function(me, npc)
		local text = "흐음.. 이곳은 빅토리아 아일랜드 중앙 던전으로 들어가는 입구지.. 조심하라구.."
		if not me:quest(2048):started() or me:class() == Class.Beginner or me:level() <= 44 then
			me:dialog(npc, text)
			return
		end

		if not me:dialog(npc, text, false, true) then
			return
		end
		if not me:dialog(npc, "흠.. #b#t4021009##k, #b#t4003002##k, #b#t4001005##k, #b#t4001006##k 재료를 얻고싶다고? 나는 경비로 일하기 전에 이 섬에 대해 약간의 공부를 했었지.", true, true) then
			return
		end
		if not me:dialog(npc, "#b#t4021009##k과 #b#t4003002##k은 아마 #m101000000#에 사는 요정이 뭔가 알고있을거야. 만약 요정이 있다면 #t4003002# 재료는 구할 수 있을거야.", true, true) then
			return
		end
		if not me:dialog(npc, "#b#t4001005##k와 #b#t4001006##k이 문제인데.. 우선 #b#t4001005##k는 골렘에게 구할 수 있지 않을까?", true, true) then
			return
		end
		me:dialog(npc, "#b#t4001006##k은.. 깃털 처럼 생긴 불꽃.. 그렇다면 불 속성의 드레이크 에게서 구할 수 있지 않을까? 그렇다면 #t4001006# 구하는 건 힘들지도 모르겠는걸...  어쨌든! 행운을 빌어!", true, false)
	end
}
