function on_start(me)
	local npc = 2010007
	local selected = me:dialog_list(npc,
		'길드를 만들고 싶은가? 혹은 길드 관련 업무를 위해서 찾아왔는가? 원하는 것을 말해보게.',
		{
			'길드를 만들고 싶습니다.',
			'길드를 해체합니다.',
			'길드 최대인원을 늘리고 싶습니다. (최대 100명)',
			'길드 최대인원을 늘리고 싶습니다. (최대 200명)',
		})
	if selected == nil then
		return
	end

	if selected == 0 then
		if me:guild_id() then
			me:dialog(npc, '흐음.. 이미 길드에 가입되어 있는 것 같은데?')
			return
		end
		if not me:dialog_yes_no(npc, '길드 제작 수수료는 #b1,500,000 메소#k라네, 정말 만들어 보고 싶은가?') then
			return
		end
		if me:generic_guild_message(1) then
			me:dialog(npc, '길드가 성공적으로 생성되었습니다.')
		end
	elseif selected == 1 then
		me:dialog(npc, '아직 지원하지 않습니다.')
	elseif selected == 2 then
		me:dialog(npc, '아직 지원하지 않습니다.')
	elseif selected == 3 then
		me:dialog(npc, '아직 지원하지 않습니다.')
	end
end
