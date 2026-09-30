-- NPC name (String.wz/Npc.img.xml): 슈리

local FREE_PASS = 4031134
local FLORINA = 110000000
local FARE = 1500

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "플로리나 비치로 여행을 떠나보고 싶지는 않은가요? 빅토리아 아일랜드 가까운 곳에는 #b플로리나 비치#k라는 환상적인 해변이 있답니다. #b1500 메소#k를 내거나 #b자유여행권#k이 있다면 언제든지 저를 통해 그곳으로 갈 수 있답니다.\r\n\r\n", {
			" #b1500 메소#k를 내겠습니다.",
			" #b자유여행권#k을 가지고 있습니다.",
			" #b자유여행권#k이 뭔가요?",
		})
		if sel == nil then
			me:dialog(npc, "플로리나 비치로 여행을 떠나보고 싶지는 않은가요? 빅토리아 아일랜드 가까운 곳에는 플로리나 비치라는 환상적인 해변이 있답니다.", false, true)
			return
		end

		if sel == 1 then
			if not me:dialog_yes_no(npc, "#b1500 메소#k를 내고 플로리나비치로 가겠다는 건가요? 좋아요~ 하지만 그곳에도 몬스터가 있는 모양이니 준비하는 것을 잊지 말도록 하세요. 그럼 출항할 준비를 해야 겠군요. 자... 지금 당장 플로리나비치로 떠나보겠어요?") then
				me:dialog(npc, "플로리나 비치로 여행을 떠나보고 싶지는 않은가요? 빅토리아 아일랜드 가까운 곳에는 플로리나 비치라는 환상적인 해변이 있답니다.", false, true)
				return
			end
			local code = me:exchange({ meso = FARE }, nil)
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "메소가 부족해요. 좀더 모아서 다시 오세요~ 필드에서 몬스터를 쓰러뜨려 보던지... 방법은 여러가지가 있답니다.", false, true)
				return
			end
			me:save_location("FLORINA")
			me:map(FLORINA, 0)
			return
		end

		if sel == 2 then
			if not me:dialog_yes_no(npc, "#b자유여행권#k을 가지고 있나요? 그것만 가지고 있다면 언제든지 플로리나비치로 갈 수 있지요. 좋아요~ 하지만 그곳에도 몬스터가 있는 모양이니 준비하는 것을 잊지 말도록 하세요. 그럼 출항할 준비를 해야 겠군 그래요. 자... 지금 당장 플로리나비치로 떠나보겠어요?") then
				return
			end
			local count = 0
			for _, it in pairs(me:item(FREE_PASS)) do
				count = count + it:count()
			end
			if count < 1 then
				me:dialog(npc, "음... #b자유여행권#k은 분명 제대로 갖고 계신건가요?", false, true)
				return
			end
			me:save_location("FLORINA")
			me:map(FLORINA, 0)
			return
		end

		if not me:dialog(npc, "자유여행권은 가지고 있기만 하면 평생 무료로 언제든지 플로리나비치로 갈 수 있는 아이템이랍니다. 워낙 특별한 티켓이라 우리들도 특별히 취급하고 있었는데 얼마전 휴가를 내고 지구방위본부라는 곳을 다녀오던 길에 잃어버렸답니다..", false, true) then
			return
		end
		me:dialog(npc, "아쉽게도 결국 찾지 못하고 돌아오긴 했답니다.. 만약 그곳에 있는 누군가가 잘 주워다 주었으면 좋겠지만 말이에요... 음.. 당신이 지구방위본부로 가서 모험을 하다 보면 찾을 수도 있지 않을까요? 설명은 여기까지랍니다. 궁금한 점이 있다면 언제든지 제게 물어봐 주세요.", true, true)
	end
}
