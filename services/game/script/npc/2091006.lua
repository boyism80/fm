-- NPC name (String.wz/Npc.img.xml): 무릉도장 공고문

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "#e< 공고 >#n\r\n무릉도장에 도전할 만용이 있는 젊은이는 무릉도장으로 찾아 오시오.  - 무공 -", {
			"무릉도장에 도전해 본다.",
			"공고문을 더 자세히 읽어 본다.",
		})
		if sel == 2 then
			if me:dialog(npc, "#e< 공고 : 도전하라! >#n\r\n나는 무릉도장의 주인 무공이다. 나는 오래 전에 무릉에서 선\r\n인이 되기 위해 수련을 시작했고, 이제 나의 내공은 경지의\r\n반열에 올라섰다. 무릉도장의 주인은 나약하기 그지없는 사\r\n람이었다. 그러므로 오늘부터 무릉도장은 내가 접수하기로\r\n했다. 무릉도장은 가장 강한 사람에게만 소유할 자격이 주어\r\n진다.\r\n나에게 가르침을 받고자 하는 사람이 있다면 언제든지 도전\r\n하라! 혹은 나에게 도전하고 싶은 자라도 상관없다. 자신\r\n의 나약함을 뼈저리게 느끼게 해주겠다.", false, true) == false then
				return
			end
			me:dialog(npc, "추신 : 혼자서 도전해도 좋다. 하지만 용기가 없는 자라면 여\r\n러명이 와도 상관없다.")
			return
		end
		if sel ~= 1 then
			return
		end
		if not me:dialog_yes_no(npc, "#b(공고문에 손을 대자, 신비한 기운이 나를 감싸기 시작했다.)#k\r\n\r\n이대로 무릉도장으로 이동하시겠습니까?") then
			return
		end
		me:save_location("MULUNG_TC")
		me:map(925020000)
	end
}
