local pq = require("script/integration/lib/party_quest")

local M = {}

local DIALOGS = { resp.dialog, resp.dialog_yes_no, resp.dialog_list, resp.warp }

local MINAR = { from = 101010000, portal = "minar00", to = 240010100, item = 4031346 }
local KEROBEN = { from = 240040700, npc = 2081005, to = 240050000, item = 2210003, use = 2210003 }
local DOLPHIN = { from = 230000000, npc = 2060009, select = 1, to = 251000100, meso = 10000 }
local SPINEL_THAILAND = { from = 102000000, npc = 9000020, select = 0, to = 500000000, meso = 3000 }
local SPINEL_TAIWAN = { from = 102000000, npc = 9000020, select = 1, to = 740000000, meso = 3000 }
local SPINEL_SINGAPORE = { from = 102000000, npc = 9000020, select = 3, to = 540010000, meso = 3000 }
local SPINEL_SHANGHAI = { from = 102000000, npc = 9000020, select = 5, to = 701000000, meso = 3000 }
local SHANGHAI_POLICE = { from = 701010320, npc = 9310004, to = 701010321, quests = { [4103] = 1 }, item = 4031289 }
local TIME_GATE = { from = 802000101, npc = 9120026, select = 0, to = 802000300, quests = { [50398] = 2 }, item = 4001393 }
local HERA = { from = 102000000, npc = 9201023, to = 680000000 }
local SHANE = { from = 101000000, npc = 1032003, to = 101000100, quests = { [2050] = 2 } }
local SUBWAY_GATE = { from = 103000100, npc = 1052007, select = 1, ride = 103000300, to = 103000310, level = 20 }
local SPIEGELMANN = { from = 220000000, npc = 2042002, select = 0, to = 980000000, level = 30 }
local DOJO_NOTICE = { from = 101000000, npc = 2091006, select = 0, to = 925020000 }
local ZAKUM_ALTAR = { from = 211042300, portal = "ps00", to = 211042400 }
local ELLIN_FOREST = { from = 222020400, portal = "in00", to = 300000100 }
local ELLIN_PQ = { from = 300030000, portal = "east00", to = 300030100 }
local FREE_MARKET = { from = 100000100, portal = "market00", to = 910000000, level = 8 }
local YULETE_LAB = { from = 261000000, portal = "yurete00", to = 926130000 }
local HORNTAIL_SQUAD = { from = 240050000, npc = 2083000, to = 240050400, item = 4001086, party = true }
local RIAN = { from = 100000204, npc = 3003472, select = 1, to = 130000000, meso = 3000 }

local TEMPLE_QUESTS = {}
for _, quest in ipairs({ 3501, 3502, 3503, 3504, 3507, 3508, 3509, 3510, 3511, 3514, 3515, 3516, 3517, 3518, 3521 }) do
	TEMPLE_QUESTS[quest] = 2
end
local DEEP_TEMPLE = { RIAN }
for _, gate in ipairs({
	{ 270010100, 270010110 }, { 270010200, 270010210 }, { 270010300, 270010310 }, { 270010400, 270010410 }, { 270010500, 270020000 },
	{ 270020100, 270020110 }, { 270020200, 270020210 }, { 270020300, 270020310 }, { 270020400, 270020410 }, { 270020500, 270030000 },
	{ 270030100, 270030110 }, { 270030200, 270030210 }, { 270030300, 270030310 }, { 270030400, 270030410 }, { 270030500, 270040000 },
}) do
	DEEP_TEMPLE[#DEEP_TEMPLE + 1] = { from = gate[1], portal = "in00", to = gate[2], quests = TEMPLE_QUESTS }
end
DEEP_TEMPLE[#DEEP_TEMPLE + 1] = { from = 270040000, portal = "in00", to = 270040100, item = 4032002 }

local function platform(select, to)
	return { from = 200000100, npc = 2012006, select = select, to = to }
end

local ROUTES = {
	{ MINAR },
	{ DOLPHIN },
	{ SPINEL_THAILAND },
	{ SPINEL_TAIWAN },
	{ SPINEL_SINGAPORE },
	{ SPINEL_SHANGHAI },
	{ HERA },
	{ SHANE },
	{ SUBWAY_GATE },
	{ SPIEGELMANN },
	{ DOJO_NOTICE },
	{ ZAKUM_ALTAR },
	{ TIME_GATE },
	{ FREE_MARKET },
	{ YULETE_LAB },
	{ RIAN },
	{ platform(0, 200000111) },
	{ platform(1, 200000121) },
	{ platform(2, 200000131) },
	{ platform(4, 200000151) },
	{ MINAR, KEROBEN },
	{ MINAR, KEROBEN, HORNTAIL_SQUAD },
	{ SPINEL_SHANGHAI, SHANGHAI_POLICE },
	{ ELLIN_FOREST, ELLIN_PQ },
	DEEP_TEMPLE,
}

local UNREACHABLE = {
	[123256780] = "보스 대기실 입장 경로 없음 #370",
}

local function answer(hop, name)
	if name == resp.dialog_yes_no then
		return req.dialog { dialog_type = 1, next = true }
	elseif name == resp.dialog_list then
		return req.dialog { dialog_type = 4, next = true, selected = hop.select or 0 }
	end
	return req.dialog { dialog_type = 0, next = true }
end

function M.talk(ctx, bot, hop)
	local oid = pq.npc(ctx, bot, hop.npc)
	if oid == false then
		return false
	end
	local arrival = hop.ride or hop.to
	local p, name = bot:request(DIALOGS, req.npc_click { oid = oid })
	for _ = 1, 10 do
		if p == false then
			return ctx:fail(string.format("%s NPC %d 대화 중 응답 없음 (맵 %d)", bot:name(), hop.npc, bot:map()))
		end
		if name == resp.warp then
			if p.character.map ~= arrival then
				return ctx:fail(string.format("%s NPC %d 대화 뒤 다른 맵으로 이동: %d", bot:name(), hop.npc, p.character.map))
			end
			return pq.wait_map(ctx, bot, hop.to, 30000)
		end
		p, name = bot:request(DIALOGS, answer(hop, name))
	end
	return ctx:fail(string.format("%s NPC %d 대화가 끝나지 않음", bot:name(), hop.npc))
end

local function enter(ctx, bot, hop)
	local x, y = bot:portal_position(hop.portal)
	if x == nil then
		return ctx:fail(string.format("%s 맵 %d에 포탈이 없음: %s", bot:name(), bot:map(), hop.portal))
	end
	if pq.move(bot, x, y) == false then
		return ctx:fail(bot:name() .. " 포탈 앞으로 이동 실패: " .. hop.portal)
	end
	return pq.portal(ctx, bot, hop.portal, hop.to)
end

local function use(ctx, bot, item_id)
	local slot = bot:slot(item_id)
	if slot == nil then
		return ctx:fail(bot:name() .. " 사용할 아이템이 없음: " .. item_id)
	end
	if bot:request(resp.inventory_operation, req.use_item { slot = slot, item_id = item_id }) == false then
		return ctx:fail(bot:name() .. " 아이템 사용 응답 없음: " .. item_id)
	end
	return true
end

-- A quest requirement resets the bot (/봇초기화) unless every quest is already in its state.
local function prepare(ctx, bot, hop)
	if hop.quests ~= nil then
		local have = bot:quests()
		local states = {}
		local met = true
		for id, state in pairs(hop.quests) do
			states[#states + 1] = id .. ":" .. state
			met = met and have[id] == state
		end
		if met == false then
			local text = string.format("/봇초기화 %d %d %d - %s", bot:level(), bot:class(), bot:meso(), table.concat(states, ","))
			if pq.command(bot, text, "봇초기화 완료") == false then
				return ctx:fail(bot:name() .. " 경로 퀘스트 준비 실패")
			end
		end
	end
	if hop.level ~= nil and bot:level() < hop.level and pq.command(bot, "/레벨바꾸기 " .. hop.level, "레벨 설정") == false then
		return ctx:fail(bot:name() .. " 경로 레벨 준비 실패: " .. hop.level)
	end
	if hop.meso ~= nil and bot:meso() < hop.meso and pq.command(bot, "/메소얻기 " .. hop.meso, "획득") == false then
		return ctx:fail(bot:name() .. " 경로 메소 준비 실패: " .. hop.meso)
	end
	if hop.item ~= nil and (bot:items()[hop.item] or 0) == 0 and pq.command(bot, "/아이템생성 " .. hop.item, "아이템 생성") == false then
		return ctx:fail(bot:name() .. " 경로 아이템 준비 실패: " .. hop.item)
	end
	return true
end

local function pass(ctx, bot, hop)
	if bot:map_move(hop.from) == false then
		return false
	end
	if hop.use ~= nil and use(ctx, bot, hop.use) == false then
		return false
	end
	if hop.portal ~= nil then
		return enter(ctx, bot, hop)
	end
	return M.talk(ctx, bot, hop)
end

local function hop_to(ctx, bot, hop)
	if prepare(ctx, bot, hop) == false then
		return false
	end
	local formed = false
	if hop.party then
		local p, name = bot:request({ resp.party_created, resp.party_status_message }, req.party_operation { operation = PARTY.Create })
		if p == false then
			return ctx:fail(bot:name() .. " 경로 파티 생성 응답 없음")
		end
		formed = name == resp.party_created
	end
	local passed = pass(ctx, bot, hop)
	if formed and bot:request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }) == false then
		return ctx:fail(bot:name() .. " 경로 파티 해산 실패")
	end
	return passed
end

local function find(bot, map_id)
	for _, hops in ipairs(ROUTES) do
		if bot:walkable(map_id, hops[#hops].to) then
			return hops
		end
	end
	return nil
end

function M.travel(ctx, bot, map_id)
	if bot:map() == map_id then
		return true
	end
	if UNREACHABLE[map_id] ~= nil then
		return bot:force_move(map_id, UNREACHABLE[map_id])
	end
	if bot:walkable(map_id) then
		return bot:map_move(map_id)
	end

	local hops = find(bot, map_id)
	if hops == nil then
		return ctx:fail(string.format("%s 맵 %d로 가는 경로가 없음", bot:name(), map_id))
	end
	local first = 1
	for i = #hops, 1, -1 do
		if bot:walkable(hops[i].to) then
			first = i + 1
			break
		end
	end
	for i = first, #hops do
		if hop_to(ctx, bot, hops[i]) == false then
			return false
		end
	end
	return bot:map_move(map_id)
end

return M
