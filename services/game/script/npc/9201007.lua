-- NPC name (String.wz/Npc.img.xml): 안나 수녀님

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		if sm:get_property("state") ~= "waiting" then
			me:dialog(npc, "정말 너무나 잘 어울리는 커플이군요~ 식이 끝났다면 발렌티나 수녀님에게 말을 걸어 결혼식 이벤트를 진행할 수 있습니다.")
			return
		end
		local groom_id = tonumber(sm:get_property("groom_id"))
		local bride_id = tonumber(sm:get_property("bride_id"))
		if me:id() ~= groom_id and me:id() ~= bride_id then
			me:dialog(npc, "결혼식이 곧 시작되니 잠시만 기다리세요.")
			return
		end
		local hall = me:map():characters()
		if hall[groom_id] == nil or hall[bride_id] == nil then
			me:dialog(npc, "신랑과 신부가 웨딩홀에 모여 있는지 확인해주세요.")
			return
		end
		if not me:dialog_yes_no(npc, "두분 정말 너무나 잘 어울려요. 지금 결혼식을 시작하시겠어요? 결혼식을 시작하게 되면 하객분들은 더 이상 입장하실 수 없으니 하객분들이 모두 입장하신 후 시작하시기 바래요.") then
			me:dialog(npc, "아직 하객분들이 다 입장하지 못한 모양이지요?")
			return
		end
		sm:call_hook("on_ceremony")
	end
}
