local quest_id = 3320

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "오오! 자네 왔군! 반갑네, 반가워! 자네 덕분에 요즘은 심심하지 않다네~ ...응? 뭐라고? 여기서 연구하던 연금술사가 누구냐고? 음... 그의 이름을 알기는 했는데...  ", false, true)
	me:dialog(npc, "뭐더라? 뭐더라... 뭐더라아... 아아! 도무지 떠오르지 않는군. 혹시 그 사람 자네에게 중요한 사람인가? 웬만하면 그냥 잊어버리면... 안된다고? 그럼 어쩐다아... ", false, true)
	if not me:dialog_accept(npc, "에잇! 모르겠다. 그냥 자네가 직접 보게!") then
		me:dialog(npc, "엥? 싫은가? 자네가 싫다면 하는 수 없지만... 그럼 여기서 연구하던 연금술사는 알려줄 수가 없는데?", false, false)
		return
	end

	q:start(npc, true)
	me:map(926120200)
end
