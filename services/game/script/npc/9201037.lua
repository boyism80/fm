-- NPC name (String.wz/Npc.img.xml): 게리와 샤티마

return {
	on_click = function(me, npc)
		if me:quest(130010):record() ~= "ing" then
			me:dialog(npc, "진실한 사랑.. 우리 이쁘죠?")
			return
		end
		if not me:dialog_yes_no(npc, "\r\n당신이 사랑의 서약을 해주신다면 그 서약을 물건에 담아드리겠어요.\r\n\r\n서약의 마음을 가다듬고 오세요.") then
			me:dialog(npc, "아직 서약의 마음을 가다듬지 못하셨나요? 준비가 되면 다시 말을 거세요.")
			return
		end

		local text = me:dialog_input(npc, "지금 이 자리에서 우리를 따라 맹세하세요.\r\n\r\n#b나의 사랑을 맹세합니다#k")
		if text == nil then
			return
		end
		if text ~= "나의 사랑을 맹세합니다" then
			me:dialog(npc, "서약이 틀렸어요. 제대로 마음을 가다듬고 다시 찾아오세요.")
			return
		end
		text = me:dialog_input(npc, "\r\n#b진실한 마음으로 영원히 사랑하겠습니다#k")
		if text == nil then
			return
		end
		if text ~= "진실한 마음으로 영원히 사랑하겠습니다" then
			me:dialog(npc, "서약이 틀렸어요. 제대로 마음을 가다듬고 다시 찾아오세요.")
			return
		end
		if not me:dialog(npc, "잘 했어요. 당신의 서약을 이 곳에 담았답니다. 가져가세요. 그리고 그 서약을 절대 잊지마세요.", false, true) then
			return
		end

		local code = me:exchange(nil, { item = { [4213000] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "기타 인벤토리 슬롯을 한칸 비우신 후 다시 찾아오세요.")
			return
		end
		me:show_effect(EffectType.QuestCompletion)
	end
}
