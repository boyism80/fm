-- NPC name (String.wz/Npc.img.xml): 몽롱

return {
	on_click = function(me, npc)
		if me:dialog_yes_no(npc, "자네.. 혹시 PC방에서 접속한건 아닌가? 후후.. 그렇다면 이곳으로 들어가 보게. 익숙한 곳을 보게 될거야.") then
			me:map(193000000)
		else
			me:dialog(npc, "흠.. 그렇다면 어쩔 수 없지.", false, false)
		end
	end
}
