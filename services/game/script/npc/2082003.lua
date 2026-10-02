-- NPC name (String.wz/Npc.img.xml): 코르바

local MORPH = 2210016

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "날개가 있다면 그 곳에 갈 수도 있겠지. 허나 그것만으로는 부족해. 검보다 날카로운 바람 사이를 날려면 단단한 비늘 역시 필요하거든. 돌아오는 방법까지 아는 하프링은 이제 나뿐이지... 그곳에 가겠다면 변신시켜 주겠네. 자네의 지금 모습이 무엇이든 이 순간만큼은 #b드래곤#k이 되는 걸세...", {
			"드래곤으로 변신하고 싶어요.",
		})
		if sel == nil then
			return
		end
		if me:level() < 90 then
			me:dialog(npc, "자네는 아직 그 곳으로 가기에는 너무 약한 모양이네. 좀 더 수련을 쌓고 다시 찾아오시게.")
			return
		end
		me:use_item(MORPH)
		me:map(200090500, 0)
	end
}
