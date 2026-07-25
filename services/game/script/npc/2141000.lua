-- NPC name (String.wz/Npc.img.xml): 키르스턴

return {
	on_click = function(me, npc)
		if not me:dialog_accept(2141000, "여신의 거울만 있으면... 다시 검은 마법사를 불러낼 수 있어!... 이, 이상해... 왜 검은 마법사를 불러내지 않는 거지? 이 기운은 뭐지? 검은 마법사와는 전혀 다른... 크아아악!\r\n\r\n#b(키르스턴의 어깨에 손을 댄다.)#k") then
			return
		end
		npc:show_effect("magic")
		sleep(2000)
		local map = me:map()
		if map == nil then
			return
		end
		map:remove_npc(npc)
		local reactor = map:reactor(2709000)
		if reactor == nil then
			return
		end
		reactor:hit()
	end
}
