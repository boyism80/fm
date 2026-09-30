-- NPC name (String.wz/Npc.img.xml): 카산드라의 팬

return {
	on_click = function(me, npc)
		if me:dialog(npc, "뭐? 나에게 볼일 있어? 음...마법사라면 모를까...", false, true) == false then
			return
		end
		if me:level() < 120 then
			me:dialog(npc, "혹시 안에 들어갈 작정이라면 생각을 바꾸는 것이 좋을 거야. 하지만 정 들어가고 싶다면... 안에서도 살아남을 수 있을 정도로 강한 자들만이 내 허락하에 들어갈 수가 있어. 더 이상의 피는 보고 싶지 않단 말야. 어디 보자... 흠! 자네는 아직 레벨 120을 넘지 못했군, 용기는 가상하지만 단념해.", true, false)
			return
		end
		if not me:dialog_yes_no(npc, "혹시 안에 들어갈 작정이라면 생각을 바꾸는 것이 좋을 거다. 마법사가 아니라면 말이다.") then
			me:dialog(npc, "마법사 인가?", false, true)
			return
		end
		me:map(921100100)
	end
}
