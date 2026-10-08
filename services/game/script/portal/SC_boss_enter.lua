return {
	on_enter = function(me)
		if me:quest(4325):started() == false then
			me:message("두꺼비 요괴 퇴치 퀘스트가 없으면 들어갈 수 없습니다.", Msg.PinkText)
			return
		end
		me:map(231050000)
		me:message("영주가 등장하지 않는다면 채널을 옮겨보세요.", Msg.PinkText)
	end
}
