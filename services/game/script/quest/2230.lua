local quest_id = 2230
local SNAIL = 5000054
local EGG = 4032086

return {
	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end
		if not q:started() then
			q:start(npc, true)
			return
		end

		if me:dialog_list(npc, "드디어 저를 찾아오셨군요...여행자여...자신에게 주어진 책임을 다 했나요?", { "저에게 주어진 책임은 무엇이죠? 당신은 누구시죠?" }) == nil then
			return
		end
		me:dialog(npc, "당신의 주머니에 작은 알을 발견하셨나요? 그 알이 바로 당신에게 주어진 책임이랍니다. 세상을 혼자서 살아가기랑 어려운 일이죠. 이럴 때 마음의 친구가 있다면 그 어려움이 줄어든답니다. 펫에 대해 들어 본 적이 있나요?\r\n사람들은 펫을 기르면서 마음의 위안을 얻는 답니다. 하지만 모든 일에는 대가와 책임이 따르듯이...", false, true)
		me:dialog(npc, "펫을 키우는 데는 책임감이 필요하답니다. 펫도 하나의 생명이기 때문이죠. 먹이도 주고, 이름도 지어주고, 정신적 교감을 나누면서 친밀감을 형성하면서 모험의 친구가 되어주고 마음의 위로를 얻는 사람들이 많답니다.", true, true)
		me:dialog(npc, "당신에게 그런 마음을 알려드리고 싶어서 제가 아끼는 아이를 보내드렸답니다. 당신이 가지고 오신 그 알은 마법의 힘을 가지고 태어나는 #b룬 달팽이#k의 알이랍니다. 당신이 여기까지 소중하게 지켜서 가져와 주셨기 때문에 곧 부화하게 될거에요.", true, true)
		me:dialog(npc, "룬 달팽이는 아주 재주가 많답니다. 아이템과 메소를 줍기도 하고 HP회복 물약과 MP회복 물약을 먹어주기도 하죠. 아주 똑똑하거든요. 하지만 룬 달팽이는 마력의 힘으로 태어나기 때문에 수명이 아주 짧답니다. 수명이 지나면 인형으로 변해버리고 다시 살려낼 수도 없죠.", true, true)
		if not me:dialog_yes_no(npc, "이해하시겠어요? 모든 일에는 책임이 따르듯이 펫도 마찬가지랍니다. 이제 곧 달팽이의 알이 부화하겠군요.") then
			me:dialog(npc, "아직 마음의 준비가 되지 않으셨나요? 곧 알이 부화하니 빨리 다시 찾아와 주세요.")
			return
		end

		if me:exchange({ item = { [EGG] = 1 } }, { item = { [SNAIL] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족한건 아닌지 확인해보세요.")
			return
		end
		q:force_complete(npc)
		me:dialog(npc, "이 달팽이의 수명은 고작 #b5시간#k이랍니다. 많이 사랑해 주세요. 그러면 분명 당신도 얻는 것이 있을거에요.")
	end,
}
