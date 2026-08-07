-- NPC name (String.wz/Npc.img.xml): 유레테

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz().id
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		if map_id == 926100203 or map_id == 926110203 then
			if sm:get_property("stage") == "1" and sm:get_property("stage5") == "0" then
				sm:set_property("stage", "2")
				map:message("유레테의 중얼거림을 들었다.")
			end
			return
		end
		if map_id ~= 926100401 and map_id ~= 926110401 then
			return
		end
		if sm:get_property("summoned") == "1" then
			return
		end
		local mob_id = 9300139
		if sm:get_property("stage") == "2" then
			local sel = me:dialog_list(npc, "큭큭큭.. 이것도 괜찮지. 내 연구의 희생물이 되기에 딱 좋아. 영광으로 생각하게나. 자네들이야 말로 최고의 기계공학과 연금술이 결합되는 것을 보는 거니까 말이야!!\r\n#b", {
				"멈춰요! 당신때문에 마가티아는 전쟁 일보 직전이에요!!",
				"멈춰요! 당신이 인정받도록 도와주겠어요.",
			})
			if sel == 2 then
				me:dialog(npc, "나를 인정받도록 도와주겠다고? 너희들이? 후후.. 웃기는 소리 하지마.")
				map:message(me:name() .. "님의 설득에 유레테는 마음이 흔들리는 모습이다.")
				sm:set_property("persuade_urete", "1")
				mob_id = 9300140
			else
				me:dialog(npc, "전쟁따위, 내 알바 아니지! 너희들은 그저 나의 연구의 결과물을 눈 똑똑히 뜨고 보고 희생되어 주면 되는거야!")
			end
			me:dialog(npc, "내 연구 결과물을 똑똑히 보아라! 가라! 프랑켄슈타인!")
		end
		sm:set_property("summoned", "1")
		map:remove_npc(npc)
		map:spawn_mob(mob_id, 240, 150)
		map:block_gen(false)
	end
}
