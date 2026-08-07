-- NPC name (String.wz/Npc.img.xml): 구옹

local pq = require("script/lib/party_quest")

local ENTRY_MAP = 251010404
local EXIT_MAP = 925100700
local REWARD_MAP = 925100600
local KEY_ID = 4001117
local SEAL_A = 4001120
local SEAL_B = 4001121
local SEAL_C = 4001122
local TREASURE_REACTOR = 2512001

local function strip_pq_items(me)
	pq.remove_all(KEY_ID, me)
	pq.remove_all(SEAL_A, me)
	pq.remove_all(SEAL_B, me)
	pq.remove_all(SEAL_C, me)
end

local function open_treasure_if_empty(map)
	if map == nil or pq.mob_count(map) > 0 then
		return
	end
	local reactor = map:reactor(TREASURE_REACTOR)
	if reactor ~= nil and reactor:state() == 0 then
		reactor:hit(1)
	end
end

local function handle_stage2(me, npc, sm, map)
	local stage = sm:get_property("stage2")
	if stage == "" then
		sm:set_property("stage2", "0")
		stage = "0"
	end
	if stage == "0" then
		if pq.has_item(me, SEAL_A, 20) then
			me:exchange({ item = { [SEAL_A] = 20 } }, nil)
			sm:set_property("stage2", "1")
			map:block_gen(false, 9300114)
			map:kill_all_mobs()
			sm:notice("구옹이 포탈의 첫번째 봉인을 해제했습니다.")
			me:dialog(npc, "#b#t4001120##k를 모두 모아오셨군요. 다음 문제를 해결할 준비가 되면 제게 다시 말을 걸어주세요.")
		else
			map:block_gen(true, 9300114)
			me:dialog(npc, "나타나는 해적을 잡고 #b#t4001120##k 20개를 제게 모아오시면 됩니다. 해적이 바로 나타나지 않더라도 잠시만 기다려 보세요. 행운을 빌어요!")
		end
	elseif stage == "1" then
		if pq.has_item(me, SEAL_B, 20) then
			me:exchange({ item = { [SEAL_B] = 20 } }, nil)
			sm:set_property("stage2", "2")
			map:block_gen(false, 9300115)
			map:kill_all_mobs()
			sm:notice("구옹이 포탈의 두번째 봉인을 해제했습니다.")
			me:dialog(npc, "#b#t4001121##k를 모두 모아오셨군요. 다음 문제를 해결할 준비가 되면 제게 다시 말을 걸어주세요.")
		else
			map:block_gen(true, 9300115)
			me:dialog(npc, "나타나는 해적을 잡고 #b#t4001121##k 20개를 제게 모아오시면 됩니다. 해적이 바로 나타나지 않더라도 잠시만 기다려 보세요. 행운을 빌어요!")
		end
	elseif stage == "2" then
		if pq.has_item(me, SEAL_C, 20) then
			me:exchange({ item = { [SEAL_C] = 20 } }, nil)
			sm:set_property("stage2", "3")
			map:block_gen(false, 9300116)
			map:kill_all_mobs()
			sm:notice("구옹이 포탈의 마지막 봉인을 해제했습니다.")
			me:dialog(npc, "#b#t4001122##k를 모두 모아오셨군요. 우측에 있는 포탈의 봉인이 풀렸으니 포탈을 통해 다음 맵으로 이동해 주시기 바랍니다.")
		else
			map:block_gen(true, 9300116)
			me:dialog(npc, "나타나는 해적을 잡고 #b#t4001122##k 20개를 제게 모아오시면 됩니다. 해적이 바로 나타나지 않더라도 잠시만 기다려 보세요. 행운을 빌어요!")
		end
	else
		me:dialog(npc, "다음 스테이지로 가는 포탈이 열렸습니다.")
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
		if map_id == EXIT_MAP then
			strip_pq_items(me)
			me:map(ENTRY_MAP)
			return
		end
		local sm = me:state_machine()
		if sm == nil then
			me:map(EXIT_MAP)
			return
		end
		if not pq.is_leader(me) then
			me:dialog(npc, "흠.. 당신은 파티장이 아니신 것 같군요.")
			return
		end
		if map_id == 925100000 or map_id == 925100200 or map_id == 925100300 then
			me:dialog(npc, "저는 도라지 왕 우양의 하인, 구옹이라고 합니다. 해적을 물리치고 도라지 왕 우양님을 구해주세요!")
		elseif map_id == 925100100 then
			handle_stage2(me, npc, sm, map)
		elseif map_id == 925100201 then
			if pq.mob_count(map) == 0 then
				me:dialog(npc, "해적들에게 넘어간 도라지들을 멋지게 혼내주셨군요. 훌륭합니다.")
				if sm:get_property("stage2a") == "0" or sm:get_property("stage2a") == "" then
					open_treasure_if_empty(map)
					sm:set_property("stage2a", "1")
				end
			else
				me:dialog(npc, "이 도라지들은 해적들의 꾀임에 넘어가 도라지 왕 우양님을 배신했습니다. ")
			end
		elseif map_id == 925100301 then
			if pq.mob_count(map) == 0 then
				me:dialog(npc, "해적들에게 넘어간 도라지들을 멋지게 혼내주셨군요. 훌륭합니다.")
				if sm:get_property("stage3a") == "0" or sm:get_property("stage3a") == "" then
					open_treasure_if_empty(map)
					sm:set_property("stage3a", "1")
				end
			else
				me:dialog(npc, "이 도라지들은 해적들의 꾀임에 넘어가 도라지 왕 우양님을 배신했습니다. ")
			end
		elseif map_id == 925100202 or map_id == 925100302 then
			me:dialog(npc, "해적들을 모두 물리치세요!")
		elseif map_id == 925100400 then
			me:dialog(npc, "해적들이 계속 나오는 문을 오래된 쇠 열쇠로 잠가버리세요!")
		elseif map_id == 925100500 then
			if pq.mob_count(map) == 0 then
				pq.party_warp(sm, REWARD_MAP)
			else
				me:dialog(npc, "해적들을 모두 물리치세요!")
			end
		end
	end
}
