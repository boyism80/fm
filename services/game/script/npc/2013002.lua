-- NPC name (String.wz/Npc.img.xml): 여신 미네르바

local pq = require("script/lib/party_quest")

local RANKING_QUEST = 1203
local REPEAT_QUEST = 199602
local FEATHER = 4001158
local DIARY_BOOK = 4161014

local function clear_fx(map)
	if map == nil then
		return
	end
	map:show_effect("quest/party/clear")
	map:play_sound("Party1/Clear")
end

local function has_full_diary(me)
	for i = 0, 9 do
		if not pq.has_item(me, 4001064 + i, 1) then
			return false
		end
	end
	return true
end

local function exchange_diary(me, npc)
	local cost = {}
	for i = 0, 9 do
		cost[4001064 + i] = 1
	end
	local code = me:exchange({ item = cost }, { item = { [DIARY_BOOK] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "인벤토리 공간을 확보한 뒤 다시 말을 걸어주세요.")
		return
	end
	if code ~= ExchangeResult.OK then
		me:dialog(npc, "일기장 조각을 모두 모아오신 게 맞나요?")
	end
end

local function pick_reward()
	local rnum = math.random(0, 250)
	if rnum == 0 then
		return 2000004, 10
	elseif rnum == 1 then
		return 2000002, 100
	elseif rnum == 2 then
		return 2000003, 100
	elseif rnum == 3 then
		return 2000006, 50
	elseif rnum == 4 then
		return 2022000, 50
	elseif rnum == 5 then
		return 2022003, 50
	elseif rnum <= 11 then
		local scrolls = { 2040002, 2040402, 2040502, 2040505, 2040602, 2040802 }
		return scrolls[rnum - 5], 1
	elseif rnum == 12 then
		return 4003000, 70
	elseif rnum <= 18 then
		return 4010000 + (rnum - 13), 20
	elseif rnum == 19 then
		return 4010006, 15
	elseif rnum <= 26 then
		return 4020000 + (rnum - 20), 20
	elseif rnum == 27 then
		return 4020007, 10
	elseif rnum == 28 then
		return 4020008, 10
	elseif rnum == 29 then
		return 1032013, 1
	elseif rnum == 30 then
		return 1032011, 1
	elseif rnum == 31 then
		return 1032014, 1
	elseif rnum <= 35 then
		return 1102021 + (rnum - 32), 1
	elseif rnum == 36 then
		return 2040803, 1
	elseif rnum == 37 then
		return 2070011, 1
	elseif rnum <= 51 then
		local weapons = {
			2043001, 2043101, 2043201, 2043301, 2043701, 2043801,
			2044001, 2044101, 2044201, 2044301, 2044401, 2044501, 2044601, 2044701,
		}
		return weapons[rnum - 37], 1
	elseif rnum == 52 then
		return 2000004, 35
	elseif rnum == 53 then
		return 2000002, 80
	elseif rnum == 54 then
		return 2000003, 80
	elseif rnum == 55 then
		return 2000006, 35
	elseif rnum == 56 then
		return 2022000, 35
	elseif rnum == 57 then
		return 2022003, 35
	elseif rnum == 58 then
		return 4003000, 75
	elseif rnum <= 64 then
		return 4010000 + (rnum - 59), 18
	elseif rnum == 65 then
		return 4010006, 12
	elseif rnum <= 72 then
		return 4020000 + (rnum - 66), 18
	elseif rnum == 73 then
		return 4020007, 7
	elseif rnum == 74 then
		return 4020008, 7
	elseif rnum <= 95 then
		local mid = {
			2040001, 2040004, 2040301, 2040401, 2040501, 2040504, 2040601, 2040601,
			2040701, 2040704, 2040707, 2040801, 2040901, 2041001, 2041004, 2041007,
			2041010, 2041013, 2041016, 2041019, 2041022,
		}
		return mid[rnum - 74], 1
	elseif rnum <= 130 then
		return 2000004, 20
	elseif rnum <= 150 then
		return 2000005, 10
	elseif rnum <= 180 then
		return 2000002, 100
	elseif rnum <= 200 then
		return 2000006, 50
	end
	return 2000003, 100
end

local function handle_hub(me, npc, sm, map)
	if sm:get_property("clearall") == "" then
		clear_fx(map)
		sm:set_property("clearall", "clear")
		pq.party_exp(sm, 23000)
		for _, p in ipairs(sm:players()) do
			if p ~= nil then
				p:end_party_quest(RANKING_QUEST)
			end
		end
	end
	me:dialog(npc, "제 석상을 복원하고 석상에 갇힌 저, 미네르바를 구해주셔서 정말 감사합니다. 그대들에게 여신의 축복이 함께 하기를...")
	pq.party_warp(sm, 920011100)
end

local function handle_reward(me, npc)
	if has_full_diary(me) then
		me:dialog(npc, "제 일기장 조각들을 모아오셨군요. 일기장 한권으로 바꿔드리도록 하지요.")
		exchange_diary(me, npc)
		return
	end
	me:dialog(npc, "저를 구해주시고 여기까지 오신 여러분들께 다시 한번 진심으로 감사드립니다. 여러분들의 여행에 도움이 되도록 작은 선물을 준비했으니 인벤토리에 빈 공간이 있는지 확인해 주세요.")
	local item_id, count = pick_reward()
	local code = me:exchange({}, { item = { [item_id] = count, [FEATHER] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "인벤토리 공간은 최소 두칸씩 비우신 후 다시 말을 걸어주세요.")
		return
	end
	if code ~= ExchangeResult.OK then
		me:dialog(npc, "인벤토리 공간은 최소 두칸씩 비우신 후 다시 말을 걸어주세요.")
		return
	end
	local q = me:quest(REPEAT_QUEST)
	if q ~= nil then
		local n = tonumber(q:record()) or 0
		q:record(tostring(n + 1))
	end
	me:map(920011200)
end

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz():id()
		if map_id == 920010100 then
			local sm = me:state_machine()
			if sm == nil then
				return
			end
			handle_hub(me, npc, sm, map)
			return
		end
		handle_reward(me, npc)
	end
}
