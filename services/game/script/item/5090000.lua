-- Item name (String.wz/Cash.img.xml): 쪽지

return {
	on_cash = function(me, item_id, text, ear, pet_sn, input)
		return me:send_note(input.recipient, text)
	end,
}
