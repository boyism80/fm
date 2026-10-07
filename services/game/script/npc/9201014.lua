-- NPC name (String.wz/Npc.img.xml): 필라

local DIVORCE_DELAY = 259200

local function divorce(me, npc)
	local result = me:request_divorce()
	if result == "requested" then
		me:dialog(npc, "이혼 신청이 접수되었습니다. 72시간 후 제게 다시 찾아오셔서 이혼을 진행하실 수 있습니다.")
	elseif result == "pending" then
		me:dialog(npc, "이미 이혼 신청이 접수되어 있습니다. 신청 후 72시간이 지나면 이혼을 진행하실 수 있습니다.")
	elseif result == nil then
		me:dialog(npc, "지금은 이혼 신청을 처리할 수 없습니다. 잠시 후 다시 찾아와 주세요.")
	end
end

-- TODO: 이혼 신청 취소 (#259)
return {
	on_click = function(me, npc)
		local text = "얼굴에 근심이 가득해 보이는군요. 저에게 상담하시고 싶은 일이 있으신가요?"
		local marriage = me:marriage()
		if marriage == nil or not marriage:married() then
			me:dialog(npc, text)
			return
		end
		if me:dialog_list(npc, text, { "이혼 상담을 하러 왔습니다." }) == nil then
			return
		end

		local requested = marriage:divorce_requested_at()
		if requested == 0 then
			if not me:dialog(npc, "이런.. 안타깝군요. 정말 이것이 최선인지 생각은 해 보셨나요? 후우.. 당신들을 처음 봤을때 부터 예감이 좋진 않았지만.. 정말 슬픈 일이군요.", false, true) then
				return
			end
			if not me:dialog_yes_no(npc, "이혼을 하고 싶다면, 3일 간의 유예 기간을 드리고, 그때까지 마음이 변하지 않으신다면 이혼을 하실 수 있습니다. 부디 마음을 바꾸고 돌아오시길 바랍니다.") then
				me:dialog(npc, "이혼은 신중하게 생각해주세요.")
				return
			end
			divorce(me, npc)
			return
		end
		if now() - requested < DIVORCE_DELAY then
			me:dialog(npc, "아직 이혼 신청을 접수하신지 72시간이 경과하지 않으신 것 같군요. 이혼은 정말 신중하게 생각하셔야 합니다.")
			return
		end
		if not me:dialog(npc, "이혼 신청을 접수하신지 72시간이 지나셨군요. 이혼을 진행하실 수 있겠군요.", false, true) then
			return
		end
		if not me:dialog_yes_no(npc, "다시 한번 묻겠습니다.. 정말 이혼을 진행하시겠습니까?") then
			me:dialog(npc, "이혼은 신중하게 생각해주세요.")
			return
		end
		divorce(me, npc)
	end
}
