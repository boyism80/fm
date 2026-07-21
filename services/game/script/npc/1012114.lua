-- NPC name (String.wz/Npc.img.xml): 어흥이

local pq = require("script/lib/party_quest")

local RICE_CAKE = 4001101
local EXIT_MAP = 910010300
local BONUS_SHORTCUT = 910010100
local RANKING_QUEST = 1200
local CLEAR_EXP = 1500

function on_click(me, npc)
	local selected = me:dialog_list(npc, "어흥! 나는 몹시 배가 고파.. 월묘가 만든 #b월묘의 떡#k을 내게 가져 와.", {
		"떡을 가져왔어요!",
		"여기서 무얼 하죠?",
		"퀘스트를 포기하고 나갑니다.",
	})
	if selected == nil then
		return
	end

	if selected == 1 then
		me:dialog(npc, "이곳은 보름달이 차면 월묘들이 떡을 만드는 달맞이 꽃이라네. 보름달이 뜨기 위해서는 달맞이꽃 씨앗을 종류별로 6개를 모아, 올바른 위치에 심으면 된다네. 보름달이 뜨고 월묘가 등장하면 다른 몬스터들로 부터 월묘를 보호하고 월묘가 만들어 내는 떡을 파티장이 내게 #b10개#k 가져오면 된다네. 만약 월묘를 지키지 못하면 퀘스트는 실패하고 나는 몹시 배가 고파지고.. 너희를 잡아먹을지도 모르지! 어흥!")
		return
	end

	if selected == 2 then
		if not me:dialog_yes_no(npc, "어흥! 정말 이곳에서 나가고 싶어?") then
			return
		end
		local party = me:party()
		local sm = me:state_machine()
		if party ~= nil and party:leader_id() == me:id() and sm ~= nil then
			pq.party_warp(sm, EXIT_MAP)
		else
			me:map(EXIT_MAP)
		end
		return
	end

	local party = me:party()
	if party == nil or party:leader_id() ~= me:id() then
		me:dialog(npc, "어흥! 넌 파티장이 아니잖아?")
		return
	end
	if not pq.has_item(me, RICE_CAKE, 10) then
		me:dialog(npc, "어흥! 어서 월묘가 만든 떡을 내게 가져와!")
		return
	end

	me:dialog(npc, "오 이것은 월묘가 만든 떡이 아닌가? 어서 내게 떡을 주게나.")
	me:dialog(npc, "냠냠 정말 맛있군. 그럼 다음에도 나를 찾아와 #b월묘의 떡#k을 구해주게. 그럼 잘 가게.")

	local sm = me:state_machine()
	local group = nil
	if sm ~= nil then
		group = sm:group()
	end
	local clear = nil
	if group ~= nil then
		clear = group:get_property("clear")
	end
	if clear ~= "1" then
		local map = me:map()
		if map ~= nil then
			map:show_effect("quest/party/clear")
			map:play_sound("Party1/Clear")
		end
		if sm ~= nil then
			pq.party_exp(sm, CLEAR_EXP)
		end
		pq.remove_all(RICE_CAKE, me)
		if group ~= nil then
			group:set_property("clear", "1")
		end
		me:end_party_quest(RANKING_QUEST)
	end

	if sm ~= nil then
		pq.party_warp(sm, BONUS_SHORTCUT)
	else
		me:map(BONUS_SHORTCUT)
	end
end
