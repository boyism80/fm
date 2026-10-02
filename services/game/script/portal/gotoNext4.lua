return {
	on_enter = function(me)
		local surveyed = me:quest(3141):completed()
		if me:quest(3143):started() then
			if surveyed then
				me:message("사자왕의 성 조사 퀘스트가 완료 되었습니다. 비스트를 찾아 가보세요.", Msg.PinkText)
			else
				me:message("길이 막혀있습니다.", Msg.PinkText)
			end
			return
		end
		if surveyed == false then
			me:message("길이 막혀있습니다.", Msg.PinkText)
			return
		end
		me:map(211060800, 0)
	end
}
