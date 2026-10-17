local pq = require("script/integration/lib/party_quest")

local MAP = 100000000
local TEMPLATE = 9901000

local function chat(text)
	return req.normal_chat { message = text }
end

test_suite {
	name = "Player NPC: 생성·채널 동기화·삭제",
	bot_count = 2,

	on_initialize = function(ctx)
		local gm, other = ctx:bot(0), ctx:bot(1)
		if gm:map_move(MAP) == false then
			return ctx:fail("GM 맵 이동 실패")
		end
		if other:transfer(1) == false then
			return ctx:fail("다른 봇 채널 이동 실패")
		end
		if other:map_move(MAP) == false then
			return ctx:fail("다른 봇 맵 이동 실패")
		end
		pq.command(gm, "/플레이어NPC삭제 " .. TEMPLATE, "플레이어 NPC 삭제")
		return true
	end,

	on_finished = function(ctx)
		pq.command(ctx:bot(0), "/플레이어NPC삭제 " .. TEMPLATE, "플레이어 NPC 삭제")
		ctx:bot(1):transfer(0)
	end,

	scenarios = {
		function(ctx)
			local gm, other = ctx:bot(0), ctx:bot(1)
			local look = gm:request_on(other, resp.player_npc, chat("/플레이어NPC생성 " .. TEMPLATE), function(p)
				return p.npc_id == TEMPLATE
			end, 5000)
			if look == false then
				return ctx:fail("다른 채널에 플레이어 NPC 외형이 오지 않음")
			end
			if look.look.name ~= gm:name() then
				return ctx:fail("플레이어 NPC 이름이 다름: " .. tostring(look.look.name))
			end
			if other:npc(TEMPLATE) == nil then
				return ctx:fail("다른 채널에 플레이어 NPC가 보이지 않음")
			end
			for _ = 1, 10 do
				if gm:npc(TEMPLATE) ~= nil then
					return true
				end
				ctx:sleep(200)
			end
			return ctx:fail("만든 채널에 플레이어 NPC가 보이지 않음")
		end,
		function(ctx)
			local gm = ctx:bot(0)
			if pq.command(gm, "/플레이어NPC생성 " .. TEMPLATE, "이미 플레이어 NPC로 쓰이는") == false then
				return ctx:fail("같은 NPC로 두 번 만들어짐")
			end
			if pq.command(gm, "/플레이어NPC생성 1", "존재하지 않는 NPC") == false then
				return ctx:fail("없는 NPC로 만들어짐")
			end
			return true
		end,
		function(ctx)
			local gm, other = ctx:bot(0), ctx:bot(1)
			local oid = other:npc(TEMPLATE)
			local removed = gm:request_on(other, resp.remove_npc, chat("/플레이어NPC삭제 " .. TEMPLATE), function(p)
				return p.oid == oid
			end, 5000)
			if removed == false then
				return ctx:fail("다른 채널에서 플레이어 NPC가 사라지지 않음")
			end
			for _ = 1, 10 do
				if gm:npc(TEMPLATE) == nil then
					return true
				end
				ctx:sleep(200)
			end
			return ctx:fail("만든 채널에서 플레이어 NPC가 사라지지 않음")
		end,
	},
}
