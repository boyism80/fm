-- NPC name (String.wz/Npc.img.xml): 핑크 벌룬

return {
	on_click = function(me, npc)
		me:dialog(npc, "보너스 맵에 오신것을 환영합니다. #b알리샤르#k를 쓰러뜨리셨다니 정말 놀랍군요. 시간이 별로 없으니 본론만 간단히 해볼까요? 이곳에는 여러 가지 상자가 있습니다. 파티원 전원이 제한 시간안에 상자를 부수고 안에 있는 물건을 입수하시면 됩니다. 운이 좋으면 정말 좋은 아이템이 나올수도 있으니 열심히 할 만한 가치가 있을 거에요. 그럼 힘내 주세요!", false, false)
	end
}
