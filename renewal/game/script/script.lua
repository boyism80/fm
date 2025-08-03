function on_start(me)
	local npc = 9001000
	-- me:dialog_list("안녕하세요", {"hello1", "hello2", "hello3"})
	-- me:dialog_accept('안녕하세요', true)
	-- me:dialog_yes_no('안녕하세요', true, true)
	local text = me:dialog_input(npc,'안녕하세요')
	me:dialog(npc, string.format("반갑습니다. %s님", me:name()))
	for i=1, 10 do
		me:chat(tostring(i))
		sleep(100)
	end
	if text ~= nil then
		me:chat(text)
	else
		me:chat('cancel')
	end
end