function on_start(me)
	local npc = 2010008
	local selected = me:dialog_list(npc, '골라보세요', {'채승현', '채진영'})
	if selected == nil then
		me:dialog(npc, '그런건 없어')
		return
	end

	if selected == 0 then
		me:dialog(npc, '채승현이 짱이긴 하지')
	else
		me:dialog(npc, '채진영은 좀 그렇긴 하지')
	end

	me:chat('ㅇㅈ')
end