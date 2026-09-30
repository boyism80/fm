-- NPC name (String.wz/Npc.img.xml): 페이슨

local TICKET = 4031134
local FARE = 1500

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "리스항구에 약간 떨어진 곳에 #b플로리나비치#k라는 환상적인 해변이 있다는 말은 들어본 적이 있는가? #b1500 메소#k를 내거나 #b자유여행권#k이 있다면 언제든지 나를 통해 그곳으로 갈 수 있다네.", {
			"#b1500 메소#k를 내겠습니다.",
			"#b자유여행권#k을 가지고 있습니다.",
			"#b자유여행권#k이 뭔가요?",
		})
		if sel == nil then
			me:dialog(npc, "플로리나 비치로 여행을 떠나보고 싶지는 않은가? 빅토리아 아일랜드 가까운 곳에는 플로리나 비치라는 환상적인 해변이 있지.", false, true)
			return
		end
		if sel == 3 then
			if me:dialog(npc, "후후... #b자유여행권#k이 뭔지 궁금해진 모양이군 그래. 자유여행권은 가지고 있기만 하면 평생 무료로 언제든지 플로리나비치로 갈 수 있는 아이템이야. 워낙 특별한 티켓이라 우리들도 특별히 취급하고 있었는데 얼마전 휴가를 내고 지구방위본부라는 곳을 다녀오던 길에 잃어버렸지 뭔가.", false, true) == false then
				return
			end
			me:dialog(npc, "아쉽게도 결국 찾지 못하고 돌아오긴 했는데 여간 찜찜한게 아냐. 그곳에 있는 누군가가 잘 주워다 주었으면 좋겠지만 말야... 아무튼 뭐 대충 이런거고 자네가 지구방위본부로 가서 모험을 하다 보면 찾을 수도 있지 않을까? 설명은 여기까지라네. 궁금한 점이 있다면 언제든지 나에게 물어봐 주길 바라네.", true, true)
			return
		end
		if sel == 1 then
			if not me:dialog_yes_no(npc, "#b1500 메소#k를 내고 플로리나비치로 가겠다는 건가? 좋아~ 하지만 그곳에도 몬스터가 있는 모양이니 준비하는 것을 잊지 말도록 하게나. 그럼 출항할 준비를 해야 겠군 그래. 자... 지금 당장 플로리나비치로 떠나볼텐가?") then
				me:dialog(npc, "플로리나 비치로 여행을 떠나보고 싶지는 않은가? 빅토리아 아일랜드 가까운 곳에는 플로리나 비치라는 환상적인 해변이 있지.", false, true)
				return
			end
			if me:exchange({ meso = FARE }, nil) ~= ExchangeResult.OK then
				me:dialog(npc, "메소가 부족한걸? 좀더 모아서 가져오라구~ 입구 있던걸 판다든지... 필드에서 몬스터를 쓰러뜨려 보던지... 방법은 여러가지가 있으니까 말야.", false, true)
				return
			end
			me:save_location("FLORINA")
			me:map(110000000, 0)
			return
		end
		if not me:dialog_yes_no(npc, "#b자유여행권#k을 가지고 있다는 건가? 그것만 가지고 있다면 언제든지 플로리나비치로 갈 수 있지. 좋아~ 하지만 그곳에도 몬스터가 있는 모양이니 준비하는 것을 잊지 말도록 하게나. 그럼 출항할 준비를 해야 겠군 그래. 자... 지금 당장 플로리나비치로 떠나볼텐가?") then
			me:dialog(npc, "플로리나 비치로 여행을 떠나보고 싶지는 않은가? 빅토리아 아일랜드 가까운 곳에는 플로리나 비치라는 환상적인 해변이 있지.", false, true)
			return
		end
		if item_count(me, TICKET) < 1 then
			me:dialog(npc, "흐음... #b자유여행권#k이 어디에 있다는 건가? 분명히 가지고 있는거야? 다시 한 번 확인해 달라고.", false, true)
			return
		end
		me:save_location("FLORINA")
		me:map(110000000, 0)
	end
}
