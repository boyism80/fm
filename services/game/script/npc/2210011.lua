-- NPC name (String.wz/Npc.img.xml): 구와르

function on_click(me, npc)
	if not me:dialog(npc, "암벽 거인이라... 그런 터무니없는 것이 만들어질 줄이야. 모든 것이 나의 과오로다. #b(구와르의 전음이 들린다.)#k", false, true) then
		return
	end
	if not me:dialog(npc, "이런 일이 일어날 것을 짐작하지 못한것은 아니지. 수백년 전, 과거의 나는 분명 검은 마법사의 세력에 동참하였고 그 일원에게 배신을 당하여 힘을 흡수당하였다... 모든 일은 그로부터 시작된 것이다. 내가 오랜 기간동안 정령들에 대한 지배력을 상실했기 때문에, 그동안 이렇게 이상한 일이 일어나고 말았지.", true, true) then
		return
	end
	if not me:dialog(npc, "이러한 일은 본래 잘못을 저질렀던 내가 책임을 져야 하지만 지금의 나는 힘이 없는 상태... 부디 미나르 숲 남부에서 암벽거인이라는 존재를 조사해주게.", true, true) then
		return
	end
	if not me:dialog(npc, "일반적인 방법으로는 암벽거인과 대화할 수 없겠지. 하지만 내가 방금 내 힘의 일부를 나누어 주었으니 분명 암벽거인과 대화할 수 있을것이다. 나의 추측이 맞다면 말이다... #b(구와르의 신비로운 힘의 일부가 몸 속으로 스며들었다.)#k", true, true) then
		return
	end
	me:dialog(npc, "그럼 이만. 필요할 때 내가 다시 접촉할 것이다", false, false)
end
