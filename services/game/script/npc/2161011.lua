-- NPC name (String.wz/Npc.img.xml): 반 레온

local PENDANT = 4032839

return {
	on_click = function(me, npc)
		if not me:dialog(npc, "넌.. 어떻게 이 길을 알고 있는 거지? 이 길은 왕족, 나와 이피아 뿐만이 모르는 길인데...", false, true) then
			return
		end
		if not me:dialog(npc, "설마... 너 진짜 이피아를 만난거냐? 나를 놀리려는 의도인가?", true, true) then
			return
		end
		if not me:dialog(npc, "마치 진실을 말하는 것 같은 표정이군. 아니... 네 말이 진실일지도 모르지. 이피아가 여기에 있고, 또 내게 말을 걸고 있을지도 몰라. 하지만 그게 무슨 소용이지? 내 손은 이미 더러워져 있는데...", true, true) then
			return
		end
		if not me:dialog(npc, "아아.. 그래서일지도 모르겠군. 내가 검은 마법사에게 영혼을 팔았기 때문에... 너무 많은 사람을 죽였기 때문에 그녀의 목소리를 들을 수 없는 것일지도 몰라. 이것이 내 죄의 대가인가...", true, true) then
			return
		end
		if not me:dialog(npc, "이피아를 알고 있는 자여. 이것을 받아다오. 오래 전, 궁정화가가 그려준 이피아의 그림이 담긴 펜던트지... 이걸 바라보며 그녀를 추억해 왔지만 이제 이런 물건은 내게 어울리지 않아.", true, true) then
			return
		end
		if not me:dialog(npc, "영혼을 바쳐 복수를 마쳤다... 남은것 따위는 없군. 이런 내게 그녀를 추억할 자격조차 없어...", true, true) then
			return
		end
		me:dialog(npc, "다시 한번 그 때로 돌아갈 수 있다면 이런 선택을 하지 않았을까? 수만 번을 생각해 봤지만 알 수 없어. 분노와 허무... 무엇을 택해도 결국 돌아오는 것은 아무것도 없으니까.")
		me:exchange(nil, { item = { [PENDANT] = 1 } })
		me:notice("퀘스트가 정상적으로 완료되지 않는다면 기타창을 비워주세요")
	end
}
