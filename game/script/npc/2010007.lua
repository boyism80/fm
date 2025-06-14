function on_start(me)
	local npc = 2010007
	local function reject()
		me:dialog(npc, '실망이에요')
	end

	if not me:dialog(npc, string.format('안녕하세요 %s님', me:name())) then
		return
	end

	if not me:dialog(npc, '제가 할말이 있는데요') then
		return
	end

	local selected = me:dialog_list(npc, '음...', {'뭔데 뜸을 들여?', '이녀석 굉장히 소심한 녀석이군'})
	if selected == nil then
		return
	end

	if selected == 0 then
		if not me:dialog_accept(npc, '저의 부탁을 들어주세요.') then
			return reject()
		end
	elseif selected == 1 then
		me:dialog(npc, '날 함부로 판단하지마')
		return
	else
		return
	end

	me:dialog(npc, '당신은 정말 착한 사람이군요...')
	selected = me:dialog_list(npc, '당신은 절 위해 100만메소를 기꺼이 사용하실 수 있나요?', {'물론이지', '그건 좀...'})
	if selected == nil then
		return reject()
	end

	if selected ~= 0 then
		return reject()
	end

	me:dialog(npc, '좋아요. 당신의 주머니를 뒤져볼게요. 돈이 있는지 검사를 좀...')
	local meso = me:meso()
	if meso >= 1000000 then
		me:dialog(npc, '돈은 충분하군요..')
		local ok, _, removed = me:remove_meso(-1)
		if not ok then
			me:dialog(npc, string.format('돈이 부족해요..'))
			return
		end

		selected = me:dialog_list(npc, '에라 모르겠다 돈 다 훔쳐버리기 ㅋㅋ', {'...!!!', '뭐하는새끼야 이거'})
		if selected == nil then

		elseif selected == 0 then
			me:dialog(npc, string.format('놀라셨죠.. 장난이었어요.. 다시 %d메소를 돌려드릴게요..', meso))
		else
			me:dialog(npc, '말하는 뽄새보소? 장난이었어 임마 ㅋ')
		end
		
		while true do
			local ok, _, cap = me:add_meso(removed)
			if ok then
				break
			end

			me:dialog(npc, string.format('돈이 너무 많아요.. 일단 %d메소를 돌려드릴게요..', cap))
			removed = cap
		end
	else
		me:dialog(npc, '당신은 너무 가난해서 대화를 그만두고 싶군요.')
	end
	me:dialog(npc, '그럼 이만...')
	me:chat('별 미친놈을 다 보겠네')
end