-- NPC name (String.wz/Npc.img.xml): 바르토스

local CANCEL = "그런가? 조금 더 생각해 보고 말을 걸라구."
local WRONG = "틀렸어! 넌 아직 애완동물에 대해서 많이 알지 못하는군.. 네놈이 데리고 있는 펫이 불쌍하다."

local QUESTIONS = {
	{
		prompt = "문제 1) 펫의 먹이를 팔고 있는 큐트는 어느 마을에 있을까?\r\n",
		options = { "리스항구", "헤네시스", "페리온", "엘리니아", "커닝시티", "슬리피우드" },
		answer = 2,
	},
	{
		prompt = "문제 2) 좋아. 문제를 잘 맞추었군. 자, 다음 문제 간다.\r\n이 사람들 중 펫과 관계가 없는 사람은 누구일까?\r\n",
		options = { "#p1032102#", "#p1012005#", "#p1012101#" },
		answer = 3,
	},
	{
		prompt = "문제 3) 잘 하고 있다. 그렇다면 이번엔 아주 쉬운 문제를 하나 내도록 하지. 다음 중 옳지 않은 것은 무엇인가?\r\n",
		options = {
			"펫은 아이템을 장비할 수 있습니다.",
			"펫에게 명령을 내리고, 가끔 친밀도가 성장합니다.",
			"뿐만 아니라, 먹이를 주어도 친밀도가 성장합니다.",
			"펫은 자신의 주인과 함께 몬스터를 공격합니다.",
		},
		answer = 4,
	},
	{
		prompt = "문제 4) 좋아. 음.. 그렇다면 펫은 레벨 몇 부터 사람의 말을 할 수 있을까?\r\n",
		options = { "#e1. #n#b레벨 5 #k", "#e2. #n#b레벨 10 #k", "#e3. #n#b레벨 15 #k", "#e4. #n#b레벨 20#k" },
		answer = 2,
	},
	{
		prompt = "문제 5) 마지막 문제다! #m100000000#의 #p1012004#는 #t2120000#를 판매하고 있다. 이것은 펫의 포만감을 몇이나 올려주지?\r\n",
		options = { "10", "20", "30" },
		answer = 3,
	},
}

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "음.. 뭐지? 이곳은 펫과 함께 산책을 즐기며 펫과의 친밀도를 향상시킬 수 있는 곳이라네.", {
			"이곳에 대해 알려주세요.",
			"요정 마르에게 이야기를 듣고 찾아왔는데...",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			if next(me:item(4031035)) ~= nil then
				me:dialog(npc, "뭐지? 이미 편지를 받았잖아? 미니맵의 #rM#k키를 눌러서 산책로 위에 있는 '프로드'를 찾아봐. 그리고 프로드에게 그 편지를 가져다 줘.", false, true)
				return
			end
			if not me:dialog_yes_no(npc, "흐음.. 펫과 산책을 해보고 싶다고? 좋아. 내가 도와주도록 하지. 이 편지를 내 친구 프로드에게 가져다 주면 펫과의 친밀도를 올릴 수 있을거야.") then
				me:dialog(npc, CANCEL, false, true)
				return
			end
			local code = me:exchange(nil, { item = { [4031035] = 1 } })
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "음? 인벤토리 공간은 충분히 갖고 있는건가? 다시 한번 확인해보겠나?")
				return
			end
			me:dialog(npc, "내가 준 #b#t4031035##k를 조련사 프로드에게 가져다 주면 된다네.", false, true)
			return
		end

		if not me:quest(2049):started() then
			me:dialog(npc, "#b#p1032102##k에게 이야기를 듣고 찾아왔다..? 어떤 이야기인가? #b#p1032102##k를 먼저 찾아가 보긴 한건가?")
			return
		end
		if next(me:item(4031034)) ~= nil then
			me:dialog(npc, "이미 #b#t4031034##k를 받은 것 같은데? 한개면 족하지 않나?")
			return
		end
		if not me:dialog(npc, "자네 혹시 #b움직이지 않게 된 펫#k을 데리고 있는건 아닌가? 정말 안된 일이지... 음? #b요정 마르#k에게 얘기를 듣고 왔다고? 그렇군... #b생명의 주문서#k라는 건가... 이봐 이봐~ 내가 그런걸 가지고 있을리가... 어라? 주머니에 이건 뭐지?", false, true) then
			return
		end
		if not me:dialog_yes_no(npc, "그냥 줄 수는 없어! 자네가 얼마나 펫에 대해 잘 알고 있는지 테스트를 해 봐야 겠어. 애정도 없는 녀석이라면 그런 놈과 같이 있는 펫이 불쌍하니까 말야. 도중에 한 번 이라도 틀리면 이 주문서는 절대로 가질 수 없을거야. 어때? 테스트를 받아보겠어?") then
			me:dialog(npc, CANCEL, false, true)
			return
		end
		if not me:dialog(npc, "좋다. 문제를 주도록 하지. 문제는 총 5문제 이다. 그럼 시작하도록 하지.", false, true) then
			return
		end

		for _, question in ipairs(QUESTIONS) do
			local answer = me:dialog_list(npc, question.prompt, question.options)
			if answer == nil then
				me:dialog(npc, CANCEL, false, true)
				return
			end
			if answer ~= question.answer then
				me:dialog(npc, WRONG)
				return
			end
		end

		if not me:dialog(npc, "맞았어! 으흠... 정말 펫에 대해 많은 걸 알고 있군 그래? 좋았어. 그 정도로 아는게 많다니 기꺼이 이 주문서를 주지. 물론 내껀 아니긴 하지만 말야... 남의 옷을 함부로 입고 그 안에 이런 중요한 걸 넣어둔 바보가 나쁜 거라고! 아... 아무튼 받아!", false, true) then
			return
		end
		local code = me:exchange(nil, { item = { [4031034] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "자네.. 인벤토리 슬롯이 부족한것 같은데? 슬롯이 부족하다면 주문서를 줄 수 없다네.")
			return
		end
		me:dialog(npc, "좋아. 그렇다면 이제 #b생명의 물#k과 함께 #b#p1032102##k 에게 찾아가 보도록 하게나. 행운을 비네. 하하하!")
	end
}
