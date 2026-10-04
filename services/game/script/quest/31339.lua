-- [암벽 거인] 카푸의 라이딩 (Quest.wz/QuestInfo.img.xml): [암벽 거인] 카푸의 라이딩

local quest_id = 31339

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_accept(npc, "자, 이쯤 되었으면 연료는 충분해. 준비가 되었으면 출발할까?") then
			return
		end

		me:dialog(npc, "자 그럼 시동을 켜볼까? #b(쿠쿠쿠쿠쿵... 푸시시...)#k 어? 이게 왜이러지..? 큰일났잖아 내 소중한 라이딩이 고장난것 같아... 오늘은 정비소도 쉬는날인데..이 고물단지 같은 라이딩... 진작에 팔아버렸어야 됬는데 어쩔수 없지.... 라이딩은 다음에 타는걸로 하자고 친구...", false, false)
		q:start(npc, true)
		me:map(240091000)
	end
}
