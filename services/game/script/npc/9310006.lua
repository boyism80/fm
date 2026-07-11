-- NPC name (String.wz/Npc.img.xml): 허경찰

function on_click(me, npc)
	if not me:dialog(npc, "드디어 여기까지 오셨군요. 이제부터가 본격적인 시작입니다. 제가 들여보내드리는 곳에는 대왕지네가 나타날 것입니다. 그곳에서 부디 대왕지네를 물리치고 #b지네의 붉은 구슬#k을 빼앗아 돌아오시기 바랍니다.", false, true) then
		return
	end
	me:map(701010323)
end
