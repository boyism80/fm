-- NPC name (String.wz/Npc.img.xml): 로나

return {
	on_click = function(me, npc)
		if me:records():get("gift.officiant_letter") > 0 then
			me:dialog(npc, "이미 너의 주례승락서는 가져갔다.")
			return
		end
		local count = 0
		for _, it in pairs(me:item(4213001)) do
			count = count + it:count()
		end
		if count >= 1 then
			me:exchange({ item = { [4213001] = 1 } }, nil)
		end
		me:records():set("gift.officiant_letter", 1)
		me:dialog(npc, "결혼 못해 ㅎ")
	end
}
