-- NPC name (String.wz/Npc.img.xml): 나가는 곳

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "이곳에서 나가 매표소로 돌아가고 싶으세요? 입장권은 반환되지 않습니다.") then
			return
		end
		me:map(103000000, 0)
	end
}
