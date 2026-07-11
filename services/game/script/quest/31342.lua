-- [암벽 거인] 치노와 함께 (Quest.wz/QuestInfo.img.xml): [암벽 거인] 치노와 함께

local quest_id = 31342

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_yes_no(npc, "그럼 출발할까? 승강기를 타고 암벽거인의 몸을 타고 올라 가는 거야. 워낙 거대한 몸 위에 올라가는 것이니만큼 시간이 좀 걸려. 준비를 단단히 하도록 해. 대부분 벌떼가 습격을 하는데 가끔식 운이 좋은 확률으로 벌떼들이 습격을 하지 않을수도 있어~") then
		return
	end

	me:map(240091600)
	me:notice("운이 좋게 벌떼가 나타나지 않았습니다.", 5)
	q:start(npc, true)
end

function on_end(me, npc)
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
