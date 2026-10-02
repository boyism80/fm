local PROOF = 4001261

return {
	on_enter = function(me)
		local exp = 260000
		if me:map():wz():id() == 105100401 then
			exp = 130000
		end
		if me:exchange({}, { exp = exp, item = { [PROOF] = 1 } }) ~= ExchangeResult.OK then
			me:message("기타 인벤토리 공간이 부족합니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(105100100, 0)
	end
}
