local quest_id = 2950

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "이제 전직도 했고, 충분한 자격을 갖춘 것 같군. 자네에게 긴히 부탁 하고 싶은 것이 있는데 도와줄 수 있겠나?") then
			me:dialog(npc, "정말로 머쉬킹 왕국에 도움을 줘볼 생각이 없나? 언제든지 생각이 바뀌면 나를 찾아오게나.", false, false)
			return
		end

		me:dialog(npc, "안녕하신가 난 대 무역상 골드리치라고 하네. 자네에게 특별한 기회를 주려고 하네. 곧 오픈 예정인 메이플 최고의 휴양지 골드비치 리조트를 오픈 전에 미리 체험해 볼 수 있는 기회를 주겠네. 지금 바로 이동하겠나?", false, true)
		if not me:dialog_yes_no(npc, "그런데, 자네 골드비치 리조트의 위치를 알고 있나? 이번에는 나를 통해 가겠지만 다음부터는 플로리나 비치를 통해 가면 된다네.") then
			me:dialog(npc, "Okay. In that case, I'll just give you the routes to the Kingdom of Mushroom. #bNear the west entrance of Henesys,#k you'll find an #bempty house#k. Enter the house, and turn left to enter#b<Themed Dungeon : Mushroom Castle>#k. That's the entrance to the Kingdom of Mushroom. There's not much time!", false, true)
			q:start(npc, true)
			return
		end

		me:map(120040300)
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
