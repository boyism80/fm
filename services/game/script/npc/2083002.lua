-- NPC name (String.wz/Npc.img.xml): 나무뿌리 수정

local pq = require("script/lib/party_quest")

local ENTRY_MAP = 240050000
local EXIT_MAP = 240050500
local CHOICE_MAP = 240050200
local BOSS_ENTRY = 240050400
local LIGHT_CAVE = 240050300
local DARK_CAVE = 240050310

local KEY_ITEMS = {
	4001087,
	4001088,
	4001089,
	4001090,
	4001091,
	4001092,
	4001093,
}

local function remove_keys(me)
	for _, id in ipairs(KEY_ITEMS) do
		pq.remove_all(id, me)
	end
end

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		if wz == nil then
			return
		end
		local map_id = wz.id

		local options = nil
		if map_id == ENTRY_MAP then
			options = {
				"수정에 떠오른 글귀를 자세히 들여다봅니다.",
				"수정을 만진다.",
			}
		elseif map_id >= 240050100 and map_id <= 240050105 then
			options = {
				"수정에 떠오른 글귀를 자세히 들여다봅니다.",
				"수정을 만진다.",
			}
		elseif map_id == CHOICE_MAP then
			options = {
				"수정에 떠오른 글귀를 자세히 들여다봅니다.",
				"수정을 만진다.",
			}
		elseif map_id == LIGHT_CAVE or map_id == DARK_CAVE then
			options = {
				"수정에 떠오른 글귀를 자세히 들여다봅니다.",
				"수정을 만진다.",
			}
		elseif map_id == BOSS_ENTRY then
			options = {
				"수정에 떠오른 글귀를 자세히 들여다봅니다.",
				"수정을 만진다.",
			}
		elseif map_id >= 240060000 and map_id <= 240060201 then
			options = {
				"수정에 떠오른 글귀를 자세히 들여다봅니다.",
				"수정을 만진다.",
			}
		elseif map_id == EXIT_MAP then
			options = {
				"수정을 만진다.",
			}
		else
			return
		end

		local selected = me:dialog_list(npc, "수정에 무언가 글귀가 떠올라 있습니다.", options)
		if selected == nil then
			return
		end

		if map_id == EXIT_MAP then
			remove_keys(me)
			me:map(ENTRY_MAP)
			return
		end

		if map_id == BOSS_ENTRY then
			if selected == 0 then
				me:dialog(npc, "이곳은 혼테일의 동굴 입구입니다. #b원정대의 표식#k을 클릭하여 원정대 퀘스트를 시작할 수 있습니다.")
				return
			end
			if not me:dialog_yes_no(npc, "생명의 동굴 입구로 나가시겠습니까?") then
				return
			end
			me:map(ENTRY_MAP)
			return
		end

		if map_id >= 240060000 and map_id <= 240060201 then
			if selected == 0 then
				me:dialog(npc, "이곳은 혼테일의 동굴 입구입니다. 혼테일을 물리치고 평화를 찾아오세요!")
				return
			end
			if not me:dialog_yes_no(npc, "이 곳에서 나가시겠습니까? 혼테일의 동굴에는 하루 2회 입장할 수 있으며, 퇴장 시 입장 가능 회수가 1회 남습니다.") then
				return
			end
			me:map(BOSS_ENTRY)
			return
		end

		if selected == 0 then
			if map_id == ENTRY_MAP then
				me:dialog(npc, "이곳은 혼테일의 동굴 입구로써, #b6명#k의 파티원이 퀘스트에 도전할 수 있습니다. #b6명#k의 파티원이 모이면, #b혼테일의 이정표#k를 눌러 퀘스트를 시작할 수 있습니다. #b명예 결사대원의 증표#k가 있다면 #b결사대원의 암호석판#k을 클릭하여 혼테일의 동굴 입구로 바로 이동할수도 있습니다.")
			elseif map_id >= 240050100 and map_id <= 240050105 then
				me:dialog(npc, "이곳은 미로방으로써, 열쇠를 이용해서 각 방을 열 수 있습니다. 각 미로방의 나무뿌리 수정에 열쇠를 넣으면 미로방의 나무뿌리 수정에서 열쇠가 나타나며, 나타난 열쇠를 그루터기에 넣으면 다음 미로방이 열리게 되며, 다섯번째 미로방까지 이것을 반복하시면 됩니다.")
			elseif map_id == CHOICE_MAP then
				me:dialog(npc, "이곳은 선택의 동굴로써, 전구를 침으로써 빛의 동굴과 어둠의 동굴을 선택할 수 있습니다. 동굴을 선택한 후 포탈을 타면 파티원 전체가 선택된 동굴로 이동됩니다.")
			elseif map_id == LIGHT_CAVE or map_id == DARK_CAVE then
				me:dialog(npc, "이곳은 빛과 어둠의 동굴로써, 몬스터를 잡아 푸른 열쇠를 획득하여 #b혼테일의 이정표#k에 가져다 주면 됩니다.")
			end
			return
		end

		if not me:dialog_yes_no(npc, "포기하고 나가시겠습니까? 파티장일시 파티 전체가 나가지게 됩니다.") then
			return
		end
		local party = me:party()
		local sm = me:state_machine()
		if party ~= nil and party:leader_id() == me:id() and sm ~= nil then
			pq.party_warp(sm, EXIT_MAP)
		else
			me:map(EXIT_MAP)
		end
	end
}
