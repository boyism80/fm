-- NPC name (String.wz/Npc.img.xml): 리아

return {
	on_click = function(me, npc)
		me:dialog(npc, "휴대폰으로 매일매일 메이플스토리~ 메이플 멤버샵에 관심이 있으세요?", false, false)
	end
}
