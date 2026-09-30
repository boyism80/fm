-- NPC name (String.wz/Npc.img.xml): 헤라

return {
	on_click = function(me, npc)
		if me:dialog(npc, "아~ 오늘은 정말 멋진 날이야! 세상은 정말 아름다워~! 이 세상에 사랑이 가득한 것 같지 않아? 웨딩빌리지에 가득한 사랑의 기운이 이 곳까지 흘러넘치고 있는 것 같아~!", false, true) == false then
			return
		end
		if not me:dialog_yes_no(npc, "웨딩빌리지에 가본 적 있어? 그 곳은 사랑이 넘쳐가는 곳이지~. 그곳에서는 사랑하는 사람과 결혼도 할 수 있대. 정~말 낭만적이지 않아? 네가 그 곳에 가고 싶다면 내가 보내줄 수 있어. 어때 한 번 가볼래?") then
			me:dialog(npc, "가지 않는거야? 언제라도 가보고 싶다면 나를 찾아오도록 해.")
			return
		end
		if me:dialog(npc, "정말 잘 생각했어. 웨딩빌리지에서 사랑의 기운을 만끽하고 오라구~. 돌아올 때는 다시 이곳으로 돌아올테니까 걱정말고~", false, true) == false then
			return
		end
		me:save_location("AMORIA")
		me:map(680000000, 0)
	end
}
