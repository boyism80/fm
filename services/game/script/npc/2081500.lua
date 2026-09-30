-- NPC name (String.wz/Npc.img.xml): 새뮤엘

local INTRO = "항상 마음 속에 자유만 존재하는 자가 되고 싶은가? 그럼 해적이 되어보는 것은 어떤가? 항상 자유만을 누리고 다니는 유일한 직업이 바로 해적이지."

local ADVANCED = {
	[Class.Buccaneer] = "그대는 얼마 전에 #b바이퍼#k가 되었던 #r#h ##k이군, 수행은 잘 되시오? 진정한 힘을 찾기 위하여 더 정진하길 바라오..",
	[Class.Corsair] = "그대는 얼마 전에 #b캡틴#k이 되었던 #r#h ##k이군, 수행은 잘 되시오? 진정한 힘을 찾기 위하여 더 정진하길 바라오..",
}

local ADVANCE = {
	[Class.Marauder] = {
		class = Class.Buccaneer,
		option = "#r#h # : #b바이퍼로 전직하고 싶습니다.",
		done = "이제 당신은 최고의 해적인 #b바이퍼#k가 되었소. 그리고 당신에게 약간의 SP와 AP와 주었으니 확인해 보시오. 이로써 자네는 최고의 강한 해적이 되었소. 하지만 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요하오...",
	},
	[Class.Outlaw] = {
		class = Class.Corsair,
		option = "#r#h # : #b캡틴으로 전직하고 싶습니다.",
		done = "이제 당신은 최고의 해적인 #b캡틴#k이 되었소. 그리고 당신에게 약간의 SP와 AP와 주었으니 확인해 보시오. 이로써 자네는 최고의 강한 해적이 되었소. 하지만 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요하오...",
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
		local q = me:quest(6944)
		if q == nil or not q:completed() then
			me:dialog(npc, INTRO)
			return
		end
		local advance = ADVANCE[me:class()]
		if advance == nil then
			return
		end

		local sel = me:dialog_list(npc, "어느 클레스로 전직하기를 원하시오?\r\n그대의 영웅으로써의 자질은 증명되었소.\r\n", {
			advance.option,
			"#r#h # : #b잠시만 생각할 시간이 필요합니다.",
		})
		if sel == nil then
			return
		end
		if sel == 2 then
			me:dialog(npc, "그대는 그대의 영웅으로서의 자질을 증명했소. \r\n이제 남은 것은 해적 궁극의 길로 향하는 것 뿐.\r\n전직을 할 준비가 되면 다시 말을 걸어주시오.")
			return
		end
		if me:skill_point() > me:level() * 3 then
			me:dialog(npc, "흠... 당신은 너무 많은 #bSP#k를 갖고 있는 것 같소. 최소한 레벨 120 이전에 얻은 스킬포인트를 모두 사용하셔야 4차 전직을 하실 수 있소.")
			return
		end

		local code = me:exchange(nil, { item = { [2280003] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흠.. 소비 인벤토리 공간이 부족한 건 아닌가? 다시 한번 확인해 보게나.")
			return
		end
		me:class(advance.class)
		me:dialog(npc, advance.done, false, true)
	end
}
