-- NPC name (String.wz/Npc.img.xml): 강경찰

function on_click(me, npc)
	if not me:dialog_yes_no(npc, "이곳에서 나가 #b#m701010320##k 지역으로 돌아가고 싶으세요?") then
		return
	end
	me:map(701010320)
end
