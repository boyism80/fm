-- NPC name (String.wz/Npc.img.xml): 콜로세움

return {
	on_click = function(me, npc)
		local npc_id = npc:id()
		if not me:dialog_yes_no(npc, "전투를 마치고 이동합니다.") then
			return
		end
		if next(me:map():mobs()) == nil and me:quest(31013):started() then
			me:quest(31018):start(npc_id, "1")
		end
		me:map(200101400)
	end
}
