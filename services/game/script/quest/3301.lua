local quest_id = 3301

local gem_items = {
	4010006, 4010004, 4010007, 4020000, 4020001, 4020002, 4020003, 4020004, 4020005, 4020006, 4020007, 4020008,
}

local function item_count(me, item_id)
	local slots = me:item(item_id)
	local count = 0
	for _, it in pairs(slots) do
		count = count + it:count()
	end
	return count
end

return {
	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local choices = {}
		local labels = {}
		for _, item_id in ipairs(gem_items) do
			if item_count(me, item_id) >= 2 then
				choices[#choices + 1] = item_id
				labels[#labels + 1] = string.format("#i%d# #t%d#", item_id, item_id)
			end
		end

		if #choices == 0 then
			me:dialog(npc, "오호... 표정을 보아하니 거래할 준비가 된 모양이군. 그렇게 까지 해서 제뉴미스트에 가입하고 싶다니... 이해할 수 없지만, 뭐 좋아. 우선 그렇다면 내게 #b아무 보석의 원석 두개#k를 가져오게.", false, false)
			return
		end

		local str = "오호... 표정을 보아하니 거래할 준비가 된 모양이군. 그렇게 까지 해서 제뉴미스트에 가입하고 싶다니... 이해할 수 없지만, 뭐 좋아. 댓가로 무엇을 주겠어?\r\n\r\n#b"
		local sel = me:dialog_list(npc, str, labels)
		if sel == nil then
			return
		end

		local need_item = choices[sel]
		if need_item == nil then
			return
		end
		if item_count(me, need_item) < 2 then
			me:dialog(npc, "아이템이 없는 것 같군.", false, false)
			return
		end

		local code = me:exchange({ item = { [need_item] = 2 } }, {})
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "아이템이 없는 것 같군.", false, false)
			return
		end
		me:dialog(npc, "그럼 잠시만 기다려. 네가 제뉴미스트 협회장의 시험을 통과하도록 만들어줄, 그 물건을 구해 놓을테니.", false, true)
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
	end
}
