local pq = require("script/integration/lib/party_quest")

local MEGAPHONE = 5071000
local SUPER_MEGAPHONE = 5072000
local MSG_MEGAPHONE = 2
local MSG_SUPER_MEGAPHONE = 3

local function shout(ctx, sender, receiver, item_id, text, ear, msg_type)
	local count = sender:items()[item_id] or 0
	local slot = sender:slot(item_id)
	if slot == nil then
		return ctx:fail("확성기가 인벤토리에 없음: " .. item_id)
	end
	local use = req.use_cash_item { slot = slot, item_id = item_id, text = text, ear = ear }
	local heard = sender:request_on(receiver, resp.notice, use, function(p)
		return p.type == msg_type and p.message:find(text, 1, true) ~= nil
	end, 5000)
	if heard == false then
		return ctx:fail("확성기 메시지가 전달되지 않음: " .. item_id)
	end
	if heard.message:find(sender:name(), 1, true) == nil then
		return ctx:fail("확성기 메시지에 보낸 사람이 없음: " .. heard.message)
	end
	if (sender:items()[item_id] or 0) ~= count - 1 then
		sender:request(resp.inventory_operation, nil, nil, 3000)
	end
	if (sender:items()[item_id] or 0) ~= count - 1 then
		local state = pq.command(sender, "/봇상태", "봇상태")
		return ctx:fail("확성기가 소모되지 않음: " .. item_id .. " " .. (state and state.message or ""))
	end
	return heard
end

local function refused(ctx, sender, receiver, text)
	local count = sender:items()[MEGAPHONE] or 0
	local use = req.use_cash_item { slot = sender:slot(MEGAPHONE), item_id = MEGAPHONE, text = text }
	local heard = sender:request_on(receiver, resp.notice, use, function(p)
		return p.type == MSG_MEGAPHONE
	end, 1500)
	if heard ~= false then
		return ctx:fail("거절되어야 할 확성기가 전달됨: " .. heard.message)
	end
	if (sender:items()[MEGAPHONE] or 0) ~= count then
		return ctx:fail("거절된 확성기가 소모됨")
	end
	return true
end

test_suite {
	name = "Megaphone: 확성기·고성능 확성기",
	bot_count = 2,

	on_initialize = function(ctx)
		local maps = { 100000000, 102000000 }
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			local items = string.format("%d:5,%d:1", MEGAPHONE, SUPER_MEGAPHONE)
			if pq.command(bot, "/봇초기화 30 0 0 " .. items .. " -", "봇초기화 완료") == false then
				return ctx:fail("봇 초기화 실패")
			end
			if bot:map_move(maps[i + 1]) == false then
				return ctx:fail("맵 이동 실패: " .. maps[i + 1])
			end
		end
		return true
	end,

	scenarios = {
		function(ctx)
			return shout(ctx, ctx:bot(0), ctx:bot(1), MEGAPHONE, "확성기 테스트", false, MSG_MEGAPHONE) ~= false
		end,
		function(ctx)
			local heard = shout(ctx, ctx:bot(0), ctx:bot(1), SUPER_MEGAPHONE, "고성능 확성기 테스트", true, MSG_SUPER_MEGAPHONE)
			if heard == false then
				return false
			end
			if heard.mega_ear ~= true then
				return ctx:fail("귀 표시가 전달되지 않음")
			end
			return true
		end,
		function(ctx)
			local sender = ctx:bot(0)
			if pq.command(sender, "/확성기금지", "확성기 사용 금지: 켜짐") == false then
				return ctx:fail("확성기 금지 설정 실패")
			end
			local ok = refused(ctx, sender, ctx:bot(1), "금지 중 확성기")
			if pq.command(sender, "/확성기금지", "확성기 사용 금지: 꺼짐") == false then
				return ctx:fail("확성기 금지 해제 실패")
			end
			return ok
		end,
		function(ctx)
			return refused(ctx, ctx:bot(0), ctx:bot(1), string.rep("가", 31))
		end,
	},
}
