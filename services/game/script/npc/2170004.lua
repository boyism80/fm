-- NPC name (String.wz/Npc.img.xml): 데모스

function on_click(me, npc)
	if not me:dialog_yes_no(npc, "다시 오르비스로 돌아 가실 건가요? 그럼 저와 함께 가요.") then
		return
	end
	if not me:dialog(npc, "지체하지 말고 바로 떠나도록 해요.", false, true) then
		return
	end
	me:map(200000000)
end
