-- NPC name (String.wz/Npc.img.xml): 하르모니아

local INTRO = "항상 강한 존재가 되고 싶은가?? 그럼 전사가 되어보는 것은 어떤가? 항상 강인함을 누리고 다니는 유일한 직업이 바로 전사라네."

local ADVANCED = {
	[Class.Hero] = "자네는 얼마 전에 #b히어로#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Paladin] = "자네는 얼마 전에 #b팔라딘#k이 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.DarkKnight] = "자네는 얼마 전에 #b다크나이트#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
}

local ADVANCE = {
	[Class.Crusader] = {
		class = Class.Hero,
		option = "#r#h # : #b히어로로 전직하고 싶습니다.",
		done = "이제 자네는 최고의 전사인 #b히어로#k가 되었네. 그리고 자네에게 약간의 SP와 AP와 주었으니 확인해 보게나. 이로써 자네는 최고의 강한 전사가 되었네. 하지만 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요하다네...",
	},
	[Class.WhiteKnight] = {
		class = Class.Paladin,
		option = "#r#h # : #b팔라딘으로 전직하고 싶습니다.",
		done = "이제 자네는 최고의 전사인 #b팔라딘#k이 되었네.  그리고 자네에게 약간의 SP와 AP와 주었으니 확인해 보게나. 이로써 자네는 최고의 강한 전사가 되었네. 하지만 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요하다네... ",
	},
	[Class.DragonKnight] = {
		class = Class.DarkKnight,
		option = "#r#h # : #b다크나이트로 전직하고 싶습니다.",
		done = "이제 자네는 최고의 전사인 #b다크나이트#k가 되었네. 그리고 자네에게 약간의 SP와 AP와 주었으니 확인해 보게나. 이로써 자네는 최고의 강한 전사가 되었네. 하지만 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요하다네...",
	},
}

return {
	on_click = function(me, npc)
		if me:level() < 120 then
			me:dialog(npc, INTRO)
			return
		end
		local advanced = ADVANCED[me:class()]
		if advanced ~= nil then
			me:dialog(npc, advanced)
			return
		end
		local q = me:quest(6904)
		if q == nil or not q:completed() then
			me:dialog(npc, INTRO)
			return
		end
		local advance = ADVANCE[me:class()]
		if advance == nil then
			return
		end

		local sel = me:dialog_list(npc, "어느 클레스로 전직하기를 원하는가?\r\n자네의 영웅으로써의 자질은 증명되었네.\r\n", {
			advance.option,
			"#r#h # : #b잠시만 생각할 시간이 필요합니다.",
		})
		if sel == nil then
			return
		end
		if sel == 2 then
			me:dialog(npc, "마음을 정했다면, 나에게 다시 말을 걸어 주게..")
			return
		end
		if me:skill_point() > me:level() * 3 then
			me:dialog(npc, "흠... 당신은 너무 많은 #bSP#k를 갖고 있는 것 같군요. 최소한 레벨 120 이전에 얻은 스킬포인트를 모두 사용하셔야 4차 전직을 하실 수 있습니다.")
			return
		end

		local code = me:exchange(nil, { item = { [2280003] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흠.. 소비 인벤토리 공간이 부족한 건 아닌가? 다시 한번 확인해 보게나.")
			return
		end
		me:class(advance.class)
		if not me:dialog(npc, advance.done, false, true) then
			return
		end
		me:dialog(npc, "난 자네가 더욱 더 수련하여 지상 최고의 전사가 되는 것을 뒤에서 지켜봐주도록 하겠네.", false, true)
	end
}
