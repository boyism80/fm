-- NPC name (String.wz/Npc.img.xml): 할리

local FEE = 40000
local MAX_CAPACITY = 100

return {
	on_click = function(me, npc)
		local refuse = "이봐, 자네! 친구 목록이 부족한 건 아닌가? 친구 목록이 부족하다면 나와 거래해 보는건 어떤가?"
		if me:dialog_yes_no(npc, "오늘도 손님이 많아야 할텐데... 엇! 잠깐만! 자네 친구 목록을 늘리고 싶지는 않은가? 자네 얼굴을 보니 한 눈에 친구가 굉장히 많을것 같다는 예감이 들었다네. 어떤가... 메소만 조금 준다면 자네의 친구 목록을 늘려주겠네. 대신 같은 계정에 다른 캐릭터에게는 적용되지 않으니까 신중히 하라구. 해볼텐가?") == false then
			me:dialog(npc, refuse, false, true)
			return
		end
		if me:dialog_yes_no(npc, "좋았어! 잘생각했네. 가격은 그다지 비싸지 않아. 내가 큰 맘 먹고 #r대폭할인#k을 시작했기 때문이지. #b친구 목록 5명을 추가 하는 데 4만 메소#k라네. 물론 낱개로는 팔지 않는다네. 한 번 사두기만 하면 영원히 목록이 늘어나는 셈이니까 친구 목록이 부족한 사람이라면 해보는 것도 나쁘지 않아. 어때? 4만 메소 내겠나?") == false then
			me:dialog(npc, refuse, false, true)
			return
		end

		local lack = "자네... #b4만 메소#k는 확실히 가지고 있는 건가? 만일 그렇다면 친구 목록이 최대로 늘어나 있는건 아닌지 확인해 주게. 아무리 돈을 낸다 하더라도 친구 목록은 기본 개수를 포함해서 #b100개 이상#k 늘어날 수 없다네."
		local capacity = me:buddy_capacity()
		if capacity >= MAX_CAPACITY then
			me:dialog(npc, lack, false, true)
			return
		end
		if me:exchange({ meso = FEE }, nil) ~= ExchangeResult.OK then
			me:dialog(npc, lack, false, true)
			return
		end
		me:buddy_capacity(capacity + 5)
		me:dialog(npc, "좋았어! 자네의 친구 목록이 5개 늘어나 있을 것일세. 한 번 확인해 보게나. 그리고 만일 그래도 친구 목록이 부족하다면 언제든지 날 찾아와 주게. 몇 번이라도 늘려줄테니까 말야. 물론 공짜로는 안되지만... 그럼 잘가게~")
	end
}
