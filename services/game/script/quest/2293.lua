-- 마지막 연주곡 (Quest.wz/QuestData/2293.img): 마지막 연주곡

local quest_id = 2293

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "락 스피릿이 마지막으로 연주한 곡을 기억 하시나요? 대충 짐작이 가는 곡들이 몇 곡 있는데 제 신곡을 들려드릴 테니 어느 곡인지 골라주세요. #b기회는 단 한 번 뿐#k이니, 신중하게 선택하세요.", false, true)

		while true do
			local labels = {}
			local opts = {}
			if me:records():get("rock_spirit.song1") == 0 then
				labels[#labels + 1] = "1번 노래 듣기"
				opts[#opts + 1] = 1
			end
			if me:records():get("rock_spirit.song2") == 0 then
				labels[#labels + 1] = "2번 노래 듣기"
				opts[#opts + 1] = 2
			end
			if me:records():get("rock_spirit.song3") == 0 then
				labels[#labels + 1] = "3번 노래 듣기"
				opts[#opts + 1] = 3
			end
			labels[#labels + 1] = "정답 입력하러 가기."
			opts[#opts + 1] = 4

			local sel = me:dialog_list(npc, "자 그럼 보기부터 드릴테니 원하시면 들어보시고 이 중에 골라 주세요. #b한 번#k만 들려드리니 잘 듣고 골라 주세요.", labels)
			if sel == nil then
				return
			end
			local selected = opts[sel]

			if selected == 1 then
				if me:records():get("rock_spirit.song1") > 0 then
					return
				end
				me:dialog(npc, "들려드리는 노래를 잘 들어보세요. 준비 되셨다면 #b다음#k버튼을 누르세요.", false, true)
				me:dialog(npc, "노래를 들려드릴게요.", false, true)
				me:records():set("rock_spirit.song1", 1)
				me:play_sound("quest2288/1", false)
			elseif selected == 2 then
				if me:records():get("rock_spirit.song2") > 0 then
					return
				end
				me:dialog(npc, "들려드리는 노래를 잘 들어보세요. 준비 되셨다면 #b다음#k버튼을 누르세요.", false, true)
				me:dialog(npc, "노래를 들려드릴게요.", false, true)
				me:records():set("rock_spirit.song2", 1)
				me:play_sound("quest2288/2", false)
			elseif selected == 3 then
				if me:records():get("rock_spirit.song3") > 0 then
					return
				end
				me:dialog(npc, "들려드리는 노래를 잘 들어보세요. 준비 되셨다면 #b다음#k버튼을 누르세요.", false, true)
				me:dialog(npc, "노래를 들려드릴게요.", false, true)
				me:records():set("rock_spirit.song3", 1)
				me:play_sound("quest2288/3", false)
			elseif selected == 4 then
				local _ = me:dialog_input(npc, "자. 그럼 이제 정답을 말해 주세요. #b기회는 단 한번#k이니 꼭 신중하게 답변해 주세요. 아래 창에 #b1, 2, 3#k 숫자 중 한글자만 입력해 주세요.")
				me:dialog(npc, "그가 연주한 곡은 바로 이거였군요. 뭐 제 곡은 아니었지만, 이제야 모든 궁금증이 풀렸어요. 정말 감사해요.\r\n\r\n#fUI/UIWindow.img/QuestIcon/4/0#\r\n\r\n#fUI/UIWindow.img/QuestIcon/8/0# 50500 exp", false, true)
				me:records():remove("rock_spirit.song1")
				me:records():remove("rock_spirit.song2")
				me:records():remove("rock_spirit.song3")
				q:force_complete(npc)
				me:show_effect(EffectType.QuestCompletion)
				me:exp(me:exp() + 50500)
				return
			end
		end
	end
}
