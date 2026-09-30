-- NPC name (String.wz/Npc.img.xml): 레고르

local INTRO = "자유로움을 누리고 다니는 궁수는 항상 자유롭게 활동하는 직업입니다."

local ADVANCED = {
	[Class.Bowmaster] = "자네는 얼마 전에 #b보우마스터#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 자유로운 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Crossbowmaster] = "자네는 얼마 전에 #b신궁#k이 되었던 #r#h ##k이군, 수행은 잘 되는가? 자유로운 존재가 되기 위하여 더 정진하길 바라네..",
}

local ADVANCE = {
	[Class.Ranger] = {
		class = Class.Bowmaster,
		option = "#r#h # : #b보우마스터로 전직하고 싶습니다.",
		done = "이제 자네는 최고의 궁수인 #b보우마스터#k이 되었다네. 축하 선물로 자네에게 약간의 SP와 AP와 주었으니 확인해 보게나. 이로써 자네는 최고의 강한 궁수가 되었네. 하지만 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요하다네...",
	},
	[Class.Sniper] = {
		class = Class.Crossbowmaster,
		option = "#r#h # : #b신궁으로 전직하고 싶습니다.",
		done = "이제 자네는 최고의 궁수인 #b신궁#k이 되었다네. 축하 선물로 자네에게 약간의 SP와 AP와 주었으니 확인해 보게나. 이로써 자네는 최고의 강한 궁수가 되었네. 하지만 앞으로 더 진보하기 위해서는 끊임없는 수련이 필요하다네...",
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
		local q = me:quest(6924)
		if q == nil or not q:completed() then
			me:dialog(npc, INTRO)
			return
		end
		local advance = ADVANCE[me:class()]
		if advance == nil then
			return
		end

		local sel = me:dialog_list(npc, "어느 클레스로 전직하기를 원하는가?\r\n그대의 영웅으로써의 자질은 증명되었다.\r\n", {
			advance.option,
			"#r#h # : #b잠시만 생각할 시간이 필요합니다.",
		})
		if sel == nil then
			return
		end
		if sel == 2 then
			me:dialog(npc, "그대는 그대의 영웅으로서의 자질을 증명했습니다. \r\n이제 남은 것은 궁수 궁극의 길로 향하는 것 뿐.\r\n전직을 할 준비가 되면 다시 말을 걸어주세요.")
			return
		end
		if me:skill_point() > me:level() * 3 then
			me:dialog(npc, "흠... 당신은 너무 많은 #bSP#k를 갖고 있는 것 같군요. 최소한 레벨 120 이전에 얻은 스킬포인트를 모두 사용하셔야 4차 전직을 하실 수 있습니다.")
			return
		end

		local code = me:exchange(nil, { item = { [2280003] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흠.. 소비 인벤토리 공간이 부족한 건 아닌가요? 다시 한번 확인해 보세요.")
			return
		end
		me:class(advance.class)
		me:dialog(npc, advance.done, false, true)
	end
}
