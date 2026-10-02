return {
	on_enter = function(me)
		if me:quest(31339):completed() == false then
			me:message("아직은 가실 수 없습니다.", Msg.PinkText)
			return
		end
		me:map(240091000)
	end
}
