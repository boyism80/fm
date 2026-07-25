local quest_id = 2962

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "나에게 무슨 할말이라도 응? 골드비치에 갑자기 나타난 블랙 슬라임에 대해 묻는 거냐? 난, 무언가를 보았다. 하지만, 그게 뭔지는 모른다. 내가 본 것을 당시도 알고 싶은가?") then
			me:dialog(npc, "음.. 진짠데..", false, false)
			return
		end

		me:dialog(npc, "그 날의 기억을 되살려야한다. 잠시만 기다려라. 그 날은 달이 밝은 좋은 밤이 었다. 낮선 섬에서 마음을 달래기 위해 악기를 연주하던중이었다. 그런데 먼 바다에서 무언가 움직이는 것을 보았다.", false, true)
		q:start(npc, true)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "응? #b전직관의 추천서#k!!! 뭐야 네가.. 아니 당신이 우리 머쉬킹 왕국을 구해주러 온 용사란 말이오?", false, true)
		me:dialog(npc, "음... 알겠소. 전직관들이 인정했을 정도면 용사님이 맞겠지요. 인사가 늦었군요. 저는 머쉬킹 왕실의 경호를 맡고있는 #b경호대장#k이라고 합니다. 보시다시피 지금은 이 곳 임시거처의 경비와 요인들의 경호를 책임지고 있습니다. 상황은 좋지않지만 어쨌든 머쉬킹 왕국에 오신 것을 환영합니다.", true, true)

		local code = me:exchange({ item = { [4032375] = 1 } }, {})
		if code ~= ExchangeResult.OK then
			return
		end
		q:force_complete(npc)
		local q2312 = me:quest(2312)
		if q2312 ~= nil then
			q2312:start(npc, true)
		end
	end
}
