local pq = require("script/integration/lib/party_quest")

local LAKELIS = 9020000
local CLOTO = 9020001
local KERNING = 103000000
local STAGE1 = 103000800
local STAGE5 = 103000804
local BONUS = 103000805
local LIGATOR = 9300001
local COUPON = 4001007
local PASS = 4001008

local stage1_answers = {
	{ "전사 1차 전직을 하기 위한 최소 레벨", 10 },
	{ "전사 1차 전직을 하기 위한 최소 힘", 35 },
	{ "마법사 1차 전직을 하기 위한 최소 지력", 20 },
	{ "궁수 1차 전직을 하기 위한 최소 민첩성", 25 },
	{ "도적 1차 전직을 하기 위한 최소 민첩성", 25 },
	{ "2차 전직을 하기 위한 최소 레벨", 30 },
	{ "궁수 1차 전직을 하기 위한 최소 레벨", 10 },
	{ "레벨1에서 레벨2가 되기 위해 필요한 경험치량", 15 },
	{ "도적 1차 전직을 하기 위한 최소 레벨", 10 },
	{ "마법사 1차 전직을 하기 위한 최소 레벨", 8 },
}

local function enter_stage1(ctx, i)
	local bot = ctx:bot(i)
	local pkt = nil
	if i == 0 then
		local oid = bot:npc(LAKELIS)
		if oid == nil then
			return ctx:fail("라케리스 NPC가 맵에 없음")
		end
		pkt = req.npc_click { oid = oid }
	end

	local warp = bot:request(resp.warp, pkt, function(p)
		return p.character.map == STAGE1
	end, 15000)
	if warp == false then
		return ctx:fail(bot:name() .. " 1스테이지 미입장")
	end
	return true
end

local function talk(ctx, bot, expected)
	local oid = pq.npc(ctx, bot, CLOTO)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg == false then
		return ctx:fail(bot:name() .. " 클로토 대화 응답 없음")
	end
	if dlg.text:find(expected, 1, true) == nil then
		return ctx:fail(bot:name() .. " 예상과 다른 대화: " .. dlg.text)
	end
	bot:dialog(false)
	return true
end

local function next_stage(ctx)
	local from = ctx:bot(0):map()
	for i = 0, ctx:bot_count() - 1 do
		local bot = ctx:bot(i)
		local warp = bot:warp("next00")
		if warp == false or warp.character.map ~= from + 1 then
			return ctx:fail(bot:name() .. " 다음 스테이지 이동 실패: " .. from)
		end
	end
	return true
end

local function stage1_answer(text)
	for _, entry in ipairs(stage1_answers) do
		if text:find(entry[1], 1, true) ~= nil then
			return entry[2]
		end
	end
	return nil
end

local function solve_coupon_quiz(ctx, member, leader)
	local oid = pq.npc(ctx, member, CLOTO)
	if oid == false then
		return false
	end
	local dlg = member:npc_click(oid)
	if dlg == false then
		return ctx:fail(member:name() .. " 클로토 대화 응답 없음")
	end
	dlg = member:dialog(true)
	if dlg == nil then
		return ctx:fail(member:name() .. " 두번째 안내 없음")
	end
	dlg = member:dialog(true)
	if dlg == nil then
		return ctx:fail(member:name() .. " 문제 대화 없음")
	end
	member:dialog(false)
	local answer = stage1_answer(dlg.text)
	if answer == nil then
		return ctx:fail("모르는 문제: " .. dlg.text)
	end

	if pq.collect(ctx, member, LIGATOR, COUPON, answer) == false then
		return false
	end
	if talk(ctx, member, "정답을 맞추셨습니다") == false then
		return false
	end
	return pq.give(ctx, member, leader, PASS, 1)
end

local function solve_area_stage(ctx)
	local leader = ctx:bot(0)
	if talk(ctx, leader, "스테이지에 대해 설명해 드리겠습니다") == false then
		return false
	end

	local oid = leader:npc(CLOTO)
	if oid == nil then
		return ctx:fail("클로토 NPC가 보이지 않음")
	end
	local areas = leader:areas()
	for _, combo in ipairs(pq.combinations(#areas, 3)) do
		for j, index in ipairs(combo) do
			local area = areas[index]
			if pq.move(ctx:bot(j), area.x, area.y) == false then
				return ctx:fail(ctx:bot(j):name() .. " 영역 이동 실패: " .. index)
			end
		end

		local p, name = leader:request({ resp.dialog, resp.environment_change }, req.npc_click { oid = oid }, function(p, name)
			return name == resp.dialog or p.env:find("wrong", 1, true) ~= nil
		end)
		if p == false then
			return ctx:fail("발판 판정 응답 없음")
		end
		if name == resp.dialog then
			leader:dialog(false)
			if p.text:find("포탈이 열렸습니다", 1, true) == nil then
				return ctx:fail("예상과 다른 판정 대화: " .. p.text)
			end
			return next_stage(ctx)
		end
	end
	return ctx:fail("정답 조합을 찾지 못함: " .. leader:map())
end

test_suite {
	name = "Party Quest: 커닝 PQ",
	bot_count = 4,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/레벨바꾸기 30", "레벨 설정") == false then
				return ctx:fail(bot:name() .. " 레벨 설정 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(KERNING) == false then
				return false
			end
		end
		return true
	end,

	on_finished = function(ctx)
		ctx:bot(0):request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
	end,

	scenarios = {
		function(ctx)
			local leader = ctx:bot(0)
			if leader:request(resp.party_created, req.party_operation { operation = PARTY.Create }) == false then
				return ctx:fail("파티 생성 실패")
			end

			for i = 1, ctx:bot_count() - 1 do
				local member = ctx:bot(i)
				local invite = leader:request_on(member, resp.party_invite,
					req.party_operation { operation = PARTY.Invite, target_name = member:name() })
				if invite == false then
					return ctx:fail(member:name() .. " 초대 패킷 없음")
				end

				local joined = member:request(resp.party_update_join,
					req.party_operation { operation = PARTY.AcceptInvite, party_id = invite.party_id })
				if joined == false then
					return ctx:fail(member:name() .. " 파티 가입 실패")
				end
				if #joined.members ~= i + 1 then
					return ctx:fail("파티원 수: " .. #joined.members)
				end
			end
			return true
		end,
		{
			parallel = {
				function(ctx)
					return enter_stage1(ctx, 0)
				end,
				function(ctx)
					return enter_stage1(ctx, 1)
				end,
				function(ctx)
					return enter_stage1(ctx, 2)
				end,
				function(ctx)
					return enter_stage1(ctx, 3)
				end,
			},
		},
		function(ctx)
			local leader = ctx:bot(0)
			if talk(ctx, leader, "첫번째 스테이지에 오신 것을 환영합니다") == false then
				return false
			end
			for i = 1, ctx:bot_count() - 1 do
				if solve_coupon_quiz(ctx, ctx:bot(i), leader) == false then
					return false
				end
			end
			if talk(ctx, leader, "포탈이 열렸습니다") == false then
				return false
			end
			return next_stage(ctx)
		end,
		solve_area_stage,
		solve_area_stage,
		solve_area_stage,
		function(ctx)
			local leader = ctx:bot(0)
			if leader:map() ~= STAGE5 then
				return ctx:fail("5스테이지가 아님: " .. leader:map())
			end
			if talk(ctx, leader, "마지막 스테이지에 오신것을 환영합니다") == false then
				return false
			end
			if pq.collect(ctx, leader, nil, PASS, 10) == false then
				return false
			end
			if talk(ctx, leader, "모든 문제를 훌륭히 해결하셨습니다") == false then
				return false
			end

			for i = 0, ctx:bot_count() - 1 do
				local bot = ctx:bot(i)
				local oid = pq.npc(ctx, bot, CLOTO)
				if oid == false then
					return false
				end
				local warp = bot:request(resp.warp, req.npc_click { oid = oid }, function(p)
					return p.character.map == BONUS
				end)
				if warp == false then
					return ctx:fail(bot:name() .. " 보너스 맵 이동 실패")
				end
			end
			return true
		end,
	},
}
