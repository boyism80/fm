-- NPC name (String.wz/Npc.img.xml): 발렌티나 수녀님

local LOBBY = 680000200

local function couple(sm, me)
	local id = tostring(me:id())
	return id == sm:get_property("groom_id") or id == sm:get_property("bride_id")
end

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			return
		end

		if me:map():wz():id() == LOBBY then
			if not couple(sm, me) then
				me:dialog(npc, "신랑, 혹은 신부가 웨딩홀로 이동할 때 까지 기다려주세요.")
				return
			end
			if not me:dialog_yes_no(npc, "이만 결혼식장으로 입장하시겠어요? 결혼식이 시작되기 전 제한시간이 다 되기전에 입장해 주시기 바랍니다. 또한, 입장하시면 로비에 입장한 하객 전원이 함께 이동됩니다.") then
				me:dialog(npc, "아직 남은 일이 있나요? 결혼식이 시작되기 전에 입장하셔야 합니다.")
				return
			end
			sm:call_hook("on_enter_cathedral")
			return
		end

		local state = sm:get_property("state")
		if state == "ceremony" or (state == "kissed" and not couple(sm, me)) then
			me:dialog(npc, "결혼식이 진행중일때는 나가실 수 없습니다.")
			return
		end
		if state == "kissed" then
			if sm:get_property("finale") ~= "0" then
				me:dialog(npc, "이미 이벤트 피날레가 시작되었습니다.")
				return
			end
			if not me:dialog_yes_no(npc, "정말 다정한 한쌍의 달팽이 같군요! 이만 결혼식을 마치고 퇴장하시겠어요?") then
				me:dialog(npc, "시간이 다 되면 자동으로 퇴장되니 그 이전에 나가고 싶으시다면 언제든지 제게 말을 걸어주세요.")
				return
			end
			sm:call_hook("on_finale")
			return
		end
		if not me:dialog_yes_no(npc, "웨딩 홀로 돌아가시고 싶으세요? 결혼식이 시작하기 전까지는 언제든지 이곳으로 돌아오실 수 있습니다.") then
			me:dialog(npc, "결혼식이 곧 시작되니 조금만 기다려 주세요.")
			return
		end
		me:map(sm:map(LOBBY))
	end
}
