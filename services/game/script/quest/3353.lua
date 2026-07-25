-- Quest name (Quest.wz/Quest.img.xml): 드랭이 바라는 것

local quest_id = 3353

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "호오~ 친절한 모험가 친구! 왔는가? 오랜만이지? 자네가 정말정말 보고 싶었다네! 왜냐고? 후후후후... 전에 자네가 물어봤던 것에 대해 알아냈거든! 그, 왜 있잖은가. 그 성격 어두운 연금술사의 사념 말이야.", false, true)
		if not me:dialog_accept(npc, "그 사람의 또 다른 사념의 흔적을 알아냈거든. 자네가 관심이 있는 것 같아서 열심히 찾아봤지. 후후후... 자, 그럼 어서 그에게 그에게 가보게") then
			me:dialog(npc, "엥? 싫은가? 자네가 싫다면 하는 수 없지만... 자네가 관심 있어 하는 것 같아서 일부러 고생고생해서 알아놨는데 파웬의 호의를 이리도 무시하다니... 훌쩍훌쩍.", false, false)
			return
		end
		q:start(npc, true)
		me:map(926120200)
	end
}
