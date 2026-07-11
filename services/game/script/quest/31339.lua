-- [암벽 거인] 카푸의 라이딩 (Quest.wz/QuestInfo.img.xml): [암벽 거인] 카푸의 라이딩

local quest_id = 31339

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_accept(npc, "자, 이쯤 되었으면 연료는 충분해. 준비가 되었으면 출발할까?") then
		return
	end

	me:map(240091000)
	me:dialog(npc, "자 그럼 시동을 켜볼까? #b(쿠쿠쿠쿠쿵... 푸시시...)#k 어? 이게 왜이러지..? 큰일났잖아 내 소중한 라이딩이 고장난것 같아... 오늘은 정비소도 쉬는날인데..이 고물단지 같은 라이딩... 진작에 팔아버렸어야 됬는데 어쩔수 없지.... 라이딩은 다음에 타는걸로 하자고 친구...", false, false)
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
