-- NPC name (String.wz/Npc.img.xml): 탕윤

local QUEST = 2180
local MILK = 4031850

return {
	on_click = function(me, npc)
		if me:quest(QUEST):started() == false then
			me:dialog(npc, "나는 노틸러스호의 식량을 담당하는 탕윤이라고 하네.")
			return
		end
		if next(me:item(MILK)) ~= nil then
			me:dialog(npc, "이미 신선한 우유를 구해왔구만?")
			return
		end
		if me:dialog(npc, "자~ 그럼 나의 소중한 젖소들이 살고 있는 외양간으로 보내주지. 우유를 다 먹어치워버리는 아기젖소들을 조심하게. 노력이 헛수고가 되어버릴 수 있으니.", false, true) == false then
			return
		end
		if me:dialog(npc, "아기젖소와 어미젖소는 한눈에 구분가지 않을거야. 태어난지 얼마 안된 아기들이지만 엄청난 먹성으로 벌써 어미소만 하거든. 모습도 붕어빵처럼 똑같이 생겼으니... 가끔 나도 헷갈린다네. 그럼 잘 해보게.", false, true) == false then
			return
		end
		local sm, err = state_machine("nautilus_cow"):start_solo(me)
		if sm == nil then
			log("nautilus_cow start_solo:", err)
			me:dialog(npc, "이미 이 안에 다른 누군가가 들어가있는 것 같군. 나중에 다시 시도하게.")
		end
	end
}
