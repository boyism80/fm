-- NPC name (String.wz/Npc.img.xml): 궁수 전직교관

local LETTER = 4031010
local DARK_MARBLE = 4031013

return {
	on_click = function(me, npc)
		local letters = 0
		for _, it in pairs(me:item(LETTER)) do
			letters = letters + it:count()
		end
		if letters < 1 then
			me:dialog(npc, "더욱 더 강한 궁수로 전직하고 싶은가요?")
			return
		end

		local sel = me:dialog_list(npc, "더욱 더 강한 궁수로 전직하고 싶은가요?\r\n\r\n#fUI/UIWindow.img/UtilDlgEx/list0#\r\n#d", {
			"영웅의 자질",
		})
		if sel == nil then
			return
		end
		if not me:dialog(npc, "흠... 이건 틀림 없는 #b헬레나#k님의 편지로군요... 당신은 궁수 2차 전직 시험을 보러 이곳까지 날 찾아온 모양이군요? 좋아요... 그럼 2차 전직 시험에 대해 간단하게 설명해 드리죠. 별로 복잡하지 않으니 걱정하지 마세요.", false, true) then
			return
		end
		if not me:dialog(npc, "내가 당신을 숨어있는 어떤 맵으로 보내주겠어요. 그 곳에는 보통 필드에서는 볼 수 없는 전혀 다른 몬스터가 나타난답니다. 물론 모습은 같지만 성격은 전혀 다른 녀석들이죠. 그들은 경험치도 주지않고 일반 아이템도 주지 않아요.", false, true) then
			return
		end
		if not me:dialog(npc, "그 안에서 녀석들을 쓰러뜨리다 보면 #b검은 구슬#k이라는 구슬을 얻을 수 있을 거에요. 그건 녀석들의 추악하고 더러운 마음이 모여 만들어진 특별한 구슬이죠. 이 구슬을 30개 모아서 안에 있는 제 동료에게 말을 걸면 합격이에요.", false, true) then
			return
		end
		if not me:dialog_yes_no(npc, "한 번 안에 들어가면 임무를 완수할 때까지 밖으로 나올 수 없어요. 그리고 이 안에서 죽어도 역시 경험치가 떨어지기 때문에 조심하는 것이 좋을거에요. 단단히 준비를 해야 한다는 말이죠. 자... 지금 당장 시험을 보겠어요?") then
			me:dialog(npc, "흐음? 수련을 받지 않는건가요?")
			return
		end
		if not me:dialog(npc, "좋아요! 그럼 저 안으로 들여보내 드리죠! 안에서 몬스터를 쓰러뜨려 검은 구슬 30개를 모은 후 안에 있는 제 동료에게 말을 걸면 시험 합격의 증거품인 #b영웅의 증거#k를 받을 수 있을꺼에요.\r\n그럼 건투를 빌어요.", false, true) then
			return
		end

		local marbles = 0
		for _, it in pairs(me:item(DARK_MARBLE)) do
			marbles = marbles + it:count()
		end
		local cost = { [LETTER] = 1 }
		if marbles > 0 then
			cost[DARK_MARBLE] = marbles
		end
		local code = me:exchange({ item = cost }, nil)
		if code ~= ExchangeResult.OK then
			return
		end
		me:map(108000100 + math.random(0, 2))
	end
}
