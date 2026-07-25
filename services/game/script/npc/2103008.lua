-- NPC name (String.wz/Npc.img.xml): 기묘한 목소리

return {
	on_click = function(me, npc)
		local q3925 = me:quest(3925)
		if q3925 == nil or not q3925:completed() then
			me:notice("동굴 문은 꿈쩍도 하지 않는다.")
			return
		end
		local q3926 = me:quest(3926)
		local q3946 = me:quest(3946)
		local q3926_open = q3926 == nil or (not q3926:started() and not q3926:completed())
		local q3946_active = q3946 ~= nil and q3946:started()
		if not q3926_open and not q3946_active then
			me:notice("동굴 문은 꿈쩍도 하지 않는다.")
			return
		end
		local text = me:dialog_input(npc, "동굴의 문을 열고 싶다면 암호를 말해라...")
		if text == "열려라참깨" or text == "열려라 참깨" then
			me:notice("암호를 말하자 신비한 힘이 동굴 안으로 인도한다.")
			me:map(260010402, 1)
		else
			me:notice("동굴 문은 꿈쩍도 하지 않는다.")
		end
	end
}
