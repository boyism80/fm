-- NPC name (String.wz/Npc.img.xml): 꿈속의 헝겊인형

return {
	on_click = function(me, npc)
		local map = me:map()
		local here = map ~= nil and map:template_id() == 450004000
		local ask = "혹시, 누나를 만난다면 안부를 전해주세요."
		if here then
			ask = "저는 루시드의 동생이에요. 누나가 있는 곳 을 찾아 해매고있어요. 누나를 만난다면 안부를 전해주세요."
		end
		if not me:dialog_yes_no(npc, ask) then
			return
		end
		if id2map(450004750) == nil then
			me:dialog(npc, "아직 갈 수 없는 곳입니다.")
			return
		end
		me:map(450004750)
	end
}
