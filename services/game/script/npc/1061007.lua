-- NPC name (String.wz/Npc.img.xml): 부서지고 있는 석상

return {
	on_click = function(me, npc)
		if me:dialog_yes_no(npc, "석상에 손을 대자 어디론가 빨려드는 듯한 느낌이 듭니다. 이대로 #b슬리피우드#k로 돌아가시겠습니까?") then
			me:map(105040300)
		else
			me:dialog(npc, "석상에서 손을 떼자 아무 일도 없던 것처럼 원래대로 돌아왔습니다.", false, false)
		end
	end
}
