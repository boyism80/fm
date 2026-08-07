-- NPC name (String.wz/Npc.img.xml): 강경찰

local OUT_MAP = 701010320

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "이곳에서 나가 #b#m701010320##k 지역으로 돌아가고 싶으세요?") then
			return
		end
		me:map(OUT_MAP)
	end
}
