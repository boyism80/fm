local quest_id = 3314

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

		if not me:has_debuff(DebuffFlag.Poison) then
			me:dialog(npc, "...아직도 약을 먹지 않은 모양이군. 이 러셀론을 믿지 못한다는 건가? 알카드노 선배로서 자네에게 모범만을 보였다고 생각해 왔는데...", false, false)
			return
		end

		local sel = me:dialog_list(npc, "후후후후후.... 안색이 창백해진 걸 보니 역시 효과가 있군. 이번 실험은 성공이야! 으하하하! 역시 로이드를 해치울 정도로 튼튼한 녀석에게는 써도 괜찮군!\r\n#b", {
			"(역시 인체실험이었나!)",
		})
		if sel == nil then
			return
		end

		sel = me:dialog_list(npc, "무척 놀란 표정인걸? 그렇게 걱정할 것 없어. 그리 위험한 약은 아니야... 아니, 위험한 약이지만 해독제는 있으니까... 후후후후....#b\r\n", {
			"(병 주고 약 줘봤자 소용 없어!)",
		})
		if sel == nil then
			return
		end

		sel = me:dialog_list(npc, "이것으로 임의로 인체의 상태를 변경할 수 있게 되었군. 이제 보다 생명연금이 쉬워질 거야. 이걸로, 이제 그 녀석의 바램을 이뤄줄 수 있을지도 모르겠군...#b\r\n", {
			"그 녀석이요?",
		})
		if sel == nil then
			return
		end

		local file = "#fUI/UIWindow.img/QuestIcon/"
		me:dialog(npc, "그래... 그 녀석. 생명연금 쪽에서는 최고인 녀석이었지. 누구보다 훌륭한 실력을 가진 녀석이었는데... 녀석만 있다면 이런 연구는 금방 해결했겠지. 하지만 어쩔 수 없어... 녀석은 이미 실종되어 버렸으니까...\r\n\r\n" .. file .. "5/0#\r\n\r\n" .. file .. "8/0# 12500 exp", false, true)

		local a = math.random(0, 20)
		local n_item
		if a == 0 then
			n_item = 2022199
		elseif a >= 1 and a < 6 then
			n_item = 2022224
		elseif a >= 6 and a < 11 then
			n_item = 2022225
		elseif a >= 11 and a < 16 then
			n_item = 2022226
		else
			n_item = 2022227
		end

		local code = me:exchange({}, { item = { [2050004] = 10, [n_item] = 20 }, exp = 12500 })
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "소비창이 가득찬 것은 아닌가? 확인해 보게.", false, false)
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)

		me:dialog(npc, "녀석이 왜 실종되었는지는 아무도 몰라. 언제부턴가 녀석은 조급해했고, 사람들 몰래 알 수 없는 연구를 하기 시작했어. 아무리 물어도 어떤 연구인지 말하지 않았어. 녀석은 반쯤 미친 듯했지. 연구, 연구, 연구... 쉴새없이 연구만 했지. 생명연금에 관한... 그리고 결국, #b그 사건#k이 벌어졌지...", false, true)
		me:dialog(npc, "연금술사들의 마을이라는 마가티아에서도, 그 정도의 대형 폭발은 단 한 번도 없었어... 녀석이 어떤 실험을 했는지, 짐작조차 가지 않아. 도대체 어떤 무시무시한 것을 연구한 것일까... 녀석의 집을 조사했으니 협회장은 뭔가 알고 있을 텐데 아무 것도 말해주지 않아...", false, true)
		me:dialog(npc, "이 연구도 처음에는 녀석과 합작한 것이었어. 하지만 그 녀석은 사라졌고 더 이상 연구를 진행하기는 어려워졌지. 약에는 자신이 있는 편이지만 그래도 역시 쉽지 않아. 녀석이 하던 것이니 계속 하고 있기는 하지만... 녀석은 도대체 왜 인체의 상태를 변경하는 연구를 한 것일까...?", false, true)
		me:dialog(npc, "녀석은 아직 살아있을 거야. 그 녀석에게는, 그래야 할 이유가 있으니까.", false, false)
	end
}
