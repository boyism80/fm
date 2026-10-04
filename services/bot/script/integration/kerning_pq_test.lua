local LAKELIS = 9020000
local KERNING = 103000000
local STAGE1 = 103000800

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

test_suite {
	name = "Party Quest: 커닝 PQ 입장",
	bot_count = 4,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			local level = bot:request(resp.notice, req.normal_chat { message = "/레벨바꾸기 30" }, function(p)
				return p.message:find("레벨 설정", 1, true) ~= nil
			end)
			if level == false then
				return ctx:fail(bot:name() .. " 레벨 설정 실패")
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
	},
}
