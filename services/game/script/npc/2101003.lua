-- NPC name (String.wz/Npc.img.xml): 아딘

local QUEST = 3933
local CLONE = 9100013

return {
	on_click = function(me, npc)
		local quest = me:quest(QUEST)
		if quest:started() == false then
			me:dialog(npc, "나에게 도전하겠다고? 용기는 가상하지만 나는 지금 바쁜 일이 있어서 말이야.")
			return
		end
		if quest:mob_kills(CLONE) > 0 then
			me:dialog(npc, "내가 졌다! 너, 꽤 강한데?")
			return
		end
		if me:dialog(npc, "네가 이렇게 강할 줄 몰랐어. 너 정도면 모래그림단원이 될 수 있을지도 모르겠다는 생각이 드는걸? 모래그림단원에게 가장 중요한 건 강함이고, 넌 충분히 강한 것 같거든. 하지만 역시 한 번만 더 시험을 해보고 싶은데, 어때? 괜찮겠어?", false, true) == false then
			return
		end
		if me:dialog_accept(npc, "진짜 너의 강함을 확인하려면 역시 몸으로 부딪혀 보는 수밖에 없겠지? 너와 대련을 해보고 싶어. 걱정 말라구. 너를 해치고 싶지는 않아. 내 분신으로 널 상대해주지. 지금 당장 대련에 들어가도 괜찮겠어?") == false then
			me:dialog(npc, "마음의 준비가 필요한건가? 너무 긴장하지는 말라구.")
			return
		end
		if me:dialog(npc, "좋아. 자신만만하군.", false, true) == false then
			return
		end
		local sm, err = state_machine("adin"):start_solo(me)
		if sm == nil then
			log("adin start_solo:", err)
			me:dialog(npc, "아... 잠시만 기다려 주게. 지금 누군가가 대련장을 쓰고 있는 것 같아. 잠시 후에 다시 찾아와 주게.")
		end
	end
}
