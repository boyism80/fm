function on_start(me)
	local npc = 2010009
	local message = me:dialog_input(npc, '입력하세요.')
	if message == '' then
		me:dialog(npc, '입력하라니까')
	else
		me:dialog(npc, string.format('%s라고 입력하셨습니다.', message))
		me:chat(string.format('내가 %s 라고 입력했다고??', message), true)
	end
end