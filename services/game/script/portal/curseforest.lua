return {
	on_enter = function(me)
		local hour = datetime().hour
		local night = hour >= 17 or hour <= 7
		local dest = nil
		if me:quest(2224):started() or me:quest(2226):started() then
			dest = 910100000
		elseif me:quest(2227):completed() then
			dest = 910100001
		end
		if dest == nil or night == false then
			me:message("무언가 알 수 없는 힘으로 가로막혀 있다.")
			return
		end
		me:play_portal_sound()
		me:message("검은 소용돌이와 함께 저주받은 숲으로 들어섰다.")
		me:map(dest, 1)
	end
}
