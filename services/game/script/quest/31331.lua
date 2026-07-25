-- Quest name (custom Resonance): [암벽 거인] 대정령 구와르

local quest_id = 31331

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog(npc, "미나르숲 남부는 예로부터 이상한 일이 일어나기로 유명했지. 하지만 이번처럼 이상한 일은 처음이야. 암벽 산이 살아서 벌떡 일어나다니 말이야. 이야기를 더 들어보고 싶나?", false, true) then
			me:dialog(npc, "정말로 머쉬킹 왕국에 도움을 줘볼 생각이 없나? 언제든지 생각이 바뀌면 나를 찾아오게나.", false, false)
			return
		end
		if not me:dialog(npc, "자네도 언뜻 들으면 무슨 말인지 이해가 잘 가지 않지? 하지만 그런 일이 실제로 일어났다네.", false, true) then
			return
		end
		if not me:dialog(npc, "이러한 초 자연적인 현상을 가장 잘 설명해줄 수 있는 자라면 딱 하나 있지. 그는 바로 대정령 구와르... 한때 검은 마법사에게 현혹되어 군단장이 된 적도 있었지만. 지금은 더 이상 약한 존재가 아니야. 이곳이 아닌 다른 어딘가에서 휴식을 취하고 있다네.", false, true) then
			return
		end
		if not me:dialog_yes_no(npc, "우리 하프링족은 대대로 하늘과 바람과 숲의 친구였지. 부족대대로 내려오는 비술을 사용하면 일시적으로 대정령 구와르와 접촉할 수 있네만... 지금 그를 만나보겠는가?") then
			me:dialog(npc, "Okay. In that case, I'll just give you the routes to the Kingdom of Mushroom. #bNear the west entrance of Henesys,#k you'll find an #bempty house#k. Enter the house, and turn left to enter#b<Themed Dungeon : Mushroom Castle>#k. That's the entrance to the Kingdom of Mushroom. There's not much time!", false, true)
			q:start(npc, true)
			return
		end

		q:start(npc, true)
		me:open_npc(2210011)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog(npc, "대정령 구와르에게 들었겠지만, 이건 중대한 일일세. 하지만 자네 같은 용사들에게는 흥미가 동하는 일이기도 하겠지. 어때, 준비가 되었나? \r\n\r\n#b암벽 거인에게는 어떻게 가면 되죠? \r\n\r\n허허허, 바로 떠나려는가? 성미가 급하군. 이미 우리 하프링들 중 탐사대원 몇몇이 그곳에 머물고 있으니 도움을 받으면 될걸세. \r\n\r\n#b하프링들이?", false, true) then
			return
		end
		if not me:dialog_yes_no(npc, "그래, 우리 종족은 대부분 조용하고 평화롭고 순박한 삶은 좋아하는 편이지만... 가끔은 탐험가의 피를 가지고 태어나는 녀석들이 있단 말일세 그런 녀석들은 도무지 말릴 수가 없지. 원한다면 지금 바로 자네를 그곳으로 이동시켜주지. 어떠한가?") then
			return
		end
		q:force_complete(npc)
		me:map(240090000)
	end
}
