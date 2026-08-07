-- NPC name (String.wz/Npc.img.xml): 아우라

local pq = require("script/lib/party_quest")

local EXIT_MAP = 280090000
local FIRE_ORE = 4001018
local PAPER = 4001015
local KEY = 4001016

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm ~= nil and sm:get_property("clear") == "1" then
			me:dialog(npc, "저 쪽에 생긴 포탈을 통해 아도비스가 있던 맵으로 돌아가실 수 있습니다. 포탈을 지나는 도중에 저에게 주신 #b불의 원석#k을 #b불의 원석 조각#k으로 만들어 파티원 각각에게 나누어 드립니다. 그럼 1단계를 클리어 하신것을 축하드립니다.")
			return
		end
		if not me:dialog(npc, "폐광 동굴을 조사하러 오신 분들이군요. 여러분들은 최종 목표인 자쿰던전의 보스를 만나기 위해 필요한 물건을 입수하셔야 합니다. 그 물건을 얻으려면 우선은 물건의 재료를 얻는 것이 우선이겠지요. 재료중에 하나인 #b불의 원석#k을 바로 이곳에서 얻을 수 있습니다. 물론 쉬운 일은 아니지만 말이죠.", true, true) then
			return
		end
		if not me:dialog(npc, "이곳에는 수 많은 동굴로 통하는 입구가 있습니다. 동굴 안에는 상자들이 있는 상자를 파괴한 후 #b열쇠 7개#k를 모으셔야 합니다. 상자는 스킬 공격으로는 부술 수 없으며 오로지 일반 공격으로만 타격을 줄 수 있으니 주의해 주세요. 그 후 7개의 열쇠를 가장 안쪽 방에 있는 커다란 보물상자에 떨어뜨리면 #b불의 원석#k을 얻을 수 있을 겁니다. 떨어뜨린후 얼마간의 시간이 지나야 얻을 수 있으니 기다려 보세요.", true, true) then
			return
		end
		if not me:dialog(npc, "물론 모든 상자에 열쇠가 있는건 아닙니다. 생각치 못했던 일이 벌어질 수도 있으니 조심해야 하죠. 상자를 조사하다 보면 가끔 종이 문서가 나오는데 이것도 모아와 주시면 틀림없이 좋은 일이 있을 겁니다. 종이 문서는 30장 이상 모아오시면 됩니다. 제가 설명해 드릴수 있는건 여기까지군요.", true, true) then
			return
		end
		local sel = me:dialog_list(npc, "더 궁금한 점이 있으신가요?", {
			"#b불의 원석#k을 가져왔습니다.",
			"퀘스트를 포기하고 이 곳에서 나갑니다.",
		})
		if sel == 1 then
			if not pq.has_item(me, FIRE_ORE) then
				if not pq.is_leader(me) then
					me:dialog(npc, "파티장이 아니시군요. 파티장이 제게 말을 걸어 진행할 수 있습니다.")
					return
				end
				me:dialog(npc, "파티원들이 동굴에서 얻은 물건은 파티장이 모두 모아 저에게 주셔야 합니다. 다시 한 번 확인해 주세요.")
				return
			end
			local paper_count = pq.item_count(me, PAPER)
			local scrolls = paper_count >= 30
			local confirm
			if paper_count < 1 then
				confirm = "#b불의 원석 1개#k는 무사히 가져오셨지만, #b종이 문서#k는 하나도 모으지 못하셨군요. 파티원들이 모으신 물건은 이게 전부가 맞습니까?"
			else
				confirm = "#b불의 원석1개#k와 #b종이 문서 "
					.. tostring(paper_count)
					.. "개#k를 모아와 주셨군요. 파티원들이 모으신 물건은 이게 전부가 맞습니까?"
			end
			if not me:dialog_yes_no(npc, confirm) then
				me:dialog(npc, "한 번 더 신중히 생각해 보시고 다시 말을 걸어 주세요.")
				return
			end
			if not pq.party_all_here(me) then
				me:dialog(npc, "파티원이 아직 이곳에 모두 모이지 않으신 것 같군요. 파티원이 모두 모였는지 다시 한번 확인해 보세요.")
				return
			end
			if sm == nil then
				me:dialog(npc, "오류가 발생했습니다.")
				return
			end
			me:exchange({ item = { [FIRE_ORE] = 1 } }, nil)
			pq.remove_all(PAPER, me)
			pq.remove_all(KEY, me)
			sm:set_property("clear", "1")
			if scrolls then
				sm:set_property("paper", "1")
			end
			me:dialog(npc, "좋습니다. 저 쪽에 생긴 포탈을 통해 아도비스가 있던 맵으로 돌아가실 수 있습니다. 포탈을 지나는 도중에 저에게 주신 #b불의 원석#k을 #b불의 원석 조각#k으로 만들어 파티원 각각에게 나누어 드립니다. 그럼 1단계를 클리어 하신것을 축하드립니다. 안녕히...")
		elseif sel == 2 then
			if not me:dialog_yes_no(npc, "포기하신다면 처음부터 다시 도전하셔야 할텐데... 게다가 파티가 모두 다 함께 진행하는 퀘스트이기 때문에 한 명이라도 도중에 나간다면 클리어가 어려울지도 모르는 일... 정말 나가시겠어요?") then
				me:dialog(npc, "한 번 더 신중히 생각해 보시고 다시 말을 걸어 주세요.")
				return
			end
			if sm ~= nil then
				if pq.is_leader(me) then
					sm:finish(EXIT_MAP)
				else
					sm:unregister(me)
					me:map(EXIT_MAP)
				end
			else
				me:map(EXIT_MAP)
			end
		end
	end
}
