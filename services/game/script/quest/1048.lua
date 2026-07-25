local quest_id = 1048

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		local name = me:name()
		local score = 0

		local sel = me:dialog_list(npc,
			string.format(
				'안녕하세요? #b%s#k님. 이제 곧 전직을 하게 되셨네요. 제가 #b%s#k님에게 좋은 직업을 추천해 드리기 위해서 몇 가지 알아야 할 점이 있어서 질문을 드리려고 합니다. #b혹시 메이플스토리 세계를 처음 접해보셨나요?',
				name, name),
			{
				'네. 처음이에요.',
				'몇 번 해보긴 했지만 아직 잘 모르겠어요.',
				'조금은 익숙해진 것 같아요.',
				'이미 모든 걸 다 알고 있어요.',
			})
		if sel == nil then
			return
		end

		local str = ''
		if sel == 0 then
			str = '메이플스토리를 처음 접하셨군요.'
			score = score + 1
		elseif sel == 1 then
			str = '메이플스토리를 몇번 해 보신 경험이 있으시군요.'
			score = score + 3
		elseif sel == 2 then
			str = '메이플스토리를 몇번 해 보신 경험이 있으시군요.'
			score = score + 5
		else
			str = '메이플스토리를 오랫동안 즐겨주셨군요.'
			score = score + 7
		end
		str = str .. ' 앞으로 게임을 진행하시면 다양한 경험을 하실 수 있을 거에요. 좋은 일도 많지만 간혹 힘든 일이 생길수도 있어요. 혹시 #b게임을 진행하다가 어려운 일이 생겼다면 어떻게 하시겠어요?'

		sel = me:dialog_list(npc, str, {
			'전 스스로 다 해결할 수 있어요.',
			'세상은 도우면서 사는 것 아닌가요?',
			'도움을 청하고 싶지만 수줍음을 많이 타서...',
		})
		if sel == nil then
			return
		end

		str = ''
		if sel == 0 then
			str = '어려운 일은 스스로 해결하는 것을 좋아하시나봐요.'
			score = score + 1
		elseif sel == 1 then
			str = '그렇죠. 어려운 일이 생기면 서로 돕는것이 좋지요.'
			score = score + 3
		else
			str = '어려운 일일수록 부탁하기가 조금 부담스러운 경우가 많지요.'
			score = score + 5
		end
		str = str .. ' 이곳 저곳 여행을 하다보면 종종 몬스터가 나타난답니다. #b몬스터가 나타난다면 어떻게 하시겠어요?'

		sel = me:dialog_list(npc, str, {
			'몬스터를 피해서 멀리 도망갈 거에요.',
			'멀리서 동료들을 도와 공격할 거에요.',
			'피하지 않고 싸울거에요.',
		})
		if sel == nil then
			return
		end

		str = ''
		if sel == 0 then
			score = score + 1
			str = '가끔은 몬스터를 피하고 싶을 때도 있답니다.'
		elseif sel == 1 then
			str = '몬스터를 공격할 때는 물약을 잊지 말고 꼭 챙기세요.'
			score = score + 3
		else
			str = '몬스터를 공격할 때는 물약을 잊지 말고 꼭 챙기세요.'
			score = score + 5
		end
		str = str .. ' 이제 곧 전직을 하시게 됩니다. 전직을 하고 나면 다양한 스킬을 사용하실 수 있어요. 많은 스킬 중에 #b어떤 스킬을 좋아하세요?'

		sel = me:dialog_list(npc, str, {
			'화려하고 멋있는 스킬이 최고죠.',
			'무조건 강력한 스킬이 좋아요.',
			'별 관심 없어요.',
		})
		if sel == nil then
			return
		end

		if sel == 0 then
			score = score + 1
		elseif sel == 1 then
			score = score + 3
		else
			score = score + 5
		end

		local class = '마법사'
		local tracker_record = '200'
		if score <= 5 then
			class = '마법사'
			tracker_record = '200'
		elseif score <= 9 then
			class = '전사'
			tracker_record = '100'
		elseif score <= 13 then
			class = '도적'
			tracker_record = '400'
		elseif score <= 17 then
			class = '궁수'
			tracker_record = '300'
		elseif score <= 22 then
			class = '해적'
			tracker_record = '500'
		end

		local msg = string.format(
			'질문에 답해주셔서 감사합니다. #b%s#k님에게 추천해 드릴 직업은 %s 입니다. 전직 후 더 강한 모습으로 뵈었으면 합니다.',
			name, class)
		if not me:dialog(npc, msg, false, true) then
			return
		end

		q:start(npc, true)
		q:force_complete(npc)
		me:quest(7631):start(tracker_record)
	end
}
