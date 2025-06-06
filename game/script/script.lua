function on_start(me)
	-- me:dialog_list("안녕하세요", {"hello1", "hello2", "hello3"})
	-- me:dialog_accept('안녕하세요', true)
	-- me:dialog_yes_no('안녕하세요', true, true)
	local text = me:dialog_input('안녕하세요')
	if text ~= nil then
		me:chat(text)
	else
		me:chat('cancel')
	end
end