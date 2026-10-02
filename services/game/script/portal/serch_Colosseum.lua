return {
	on_enter = function(me)
		local quest = me:quest(31010)
		if quest:started() then
			quest:record("0")
			me:message("맘무트가 보인다. 흠.. 왠지 강해 보이는걸? 미카엘에게 돌아가 이 사실을 알려주자.", Msg.ScrollingTop)
			me:message("몬스터 맘무트를 확인 했습니다. 미카엘에게 돌아가 이야기를 전달해 주세요.", Msg.PinkText)
		elseif quest:completed() == false then
			me:message("들어가기엔 겁이 난다. 다음에 가도록 하자.")
			return
		end
		me:play_portal_sound()
		me:map(200101100, 1)
	end
}
