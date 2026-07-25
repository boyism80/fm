-- NPC name (String.wz/Npc.img.xml): 목도리 눈사람

return {
	on_click = function(me, npc)
		if me:dialog_yes_no(npc, "트리는 예쁘게 잘 꾸며보셨나요? 여러 사람들과 함께 트리를 꾸미는 것도 한 번쯤 해볼만 하다니까요~ 아참.. 이곳에서 정말 나가시겠어요?") then
			me:map(209000000)
		else
			me:dialog(npc, "나무를 꾸미는 데 더 많은 시간을 보내고 싶으세요? 만약에 이 장소에서 나가고 싶다면 저에게 말을 걸어 주세요.", false, false)
		end
	end
}
