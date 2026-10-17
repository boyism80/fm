local pq = require("script/integration/lib/party_quest")

local LITH_HARBOR = 104000000

local POTION = 2000000
local SWORD = 1302000
local SHELL = 4000000

local EQUIP = 1
local CONSUME = 2
local ETC = 4
local TRADE = 3

local MODE = {
	create = 0x00,
	invite = 0x02,
	decline = 0x03,
	visit = 0x04,
	chat = 0x06,
	exit = 0x0A,
	put_item = 0x0D,
	put_meso = 0x0E,
	confirm = 0x0F,
}
local INVITE = { declined = 3 }
local LEAVE = { exit = 0, cancel = 2, done = 6 }

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, meso, items)
	local text = string.format("/봇초기화 30 0 %d %s -", meso, items)
	if pq.command(bot, text, "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map() ~= LITH_HARBOR and bot:map_move(LITH_HARBOR) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

local function has(ctx, bot, item_id, count)
	local have = bot:items()[item_id] or 0
	return check(ctx, have == count, string.format("%s 인벤토리 %d 개수 %d (기대 %d)", bot:name(), item_id, have, count))
end

local function meso_is(ctx, bot, meso)
	return check(ctx, bot:meso() == meso, string.format("%s 메소 %d (기대 %d)", bot:name(), bot:meso(), meso))
end

local function invite(ctx, inviter, invitee, what)
	local enter = inviter:request(resp.trade_enter, req.mini_room { mode = MODE.create, type = TRADE }, nil, 5000)
	if check(ctx, enter ~= false and enter.my_slot == 0 and #enter.members == 1, "교환창 생성 응답이 다름") == false then
		return false
	end
	local p = inviter:request_on(invitee, resp.trade_invite, req.mini_room { mode = MODE.invite, target_id = invitee:id() }, nil, 5000)
	if check(ctx, p ~= false and p.inviter == inviter:name(), (what or "교환") .. " 초대가 상대에게 오지 않음") == false then
		return false
	end
	return p.sn
end

local function open(ctx, inviter, invitee)
	local sn = invite(ctx, inviter, invitee)
	if sn == false then
		return false
	end
	local enter = invitee:request(resp.trade_enter, req.mini_room { mode = MODE.visit, sn = sn }, nil, 5000)
	if check(ctx, enter ~= false and enter.my_slot == 1 and #enter.members == 2, "초대 수락 입장 응답이 다름") == false then
		return false
	end
	ctx:sleep(300)
	return true
end

local function put_item(ctx, bot, other, inv_type, item_id, count, trade_slot)
	local request = req.mini_room {
		mode = MODE.put_item,
		inventory_type = inv_type,
		slot = bot:slot(item_id) or 0,
		count = count,
		trade_slot = trade_slot,
	}
	local p = bot:request_on(other, resp.trade_item, request, nil, 5000)
	return check(ctx, p ~= false and p.who == 1 and p.slot == trade_slot, string.format("%s 아이템 %d 올리기가 상대에게 보이지 않음", bot:name(), item_id))
end

local function put_meso(ctx, bot, other, meso, total)
	local p = bot:request_on(other, resp.trade_meso, req.mini_room { mode = MODE.put_meso, meso = meso }, nil, 5000)
	return check(ctx, p ~= false and p.who == 1 and p.meso == total, string.format("%s 메소 올리기가 상대에게 보이지 않음", bot:name()))
end

local function blocked(ctx, bot, request, what)
	local p = bot:request(resp.update_stats, request, nil, 3000)
	return check(ctx, p ~= false, bot:name() .. " 교환 중 " .. what .. " 거부 응답 없음")
end

local function trade_flow(ctx)
	local a = ctx:bot(0)
	local b = ctx:bot(1)
	if setup(ctx, a, 2000000, POTION .. ":100") == false or setup(ctx, b, 0, SWORD .. ":1") == false then
		return false
	end
	local dropped = a:drop(POTION, 10)
	if check(ctx, dropped ~= nil, "줍기 확인용 물약을 떨어뜨리지 못함") == false then
		return false
	end

	if open(ctx, a, b) == false then
		return false
	end
	if put_item(ctx, a, b, CONSUME, POTION, 30, 1) == false or put_meso(ctx, a, b, 1000000, 1000000) == false then
		return false
	end
	if put_item(ctx, b, a, EQUIP, SWORD, 1, 1) == false then
		return false
	end
	ctx:sleep(300)
	if has(ctx, a, POTION, 60) == false or meso_is(ctx, a, 1000000) == false or has(ctx, b, SWORD, 0) == false then
		return false
	end

	if check(ctx, b:loot(dropped), "교환 중 물약 줍기") == false then
		return false
	end
	if blocked(ctx, a, req.move_item { inventory_type = CONSUME, source = a:slot(POTION), dest = 0, count = 1 }, "버리기") == false then
		return false
	end
	ctx:sleep(300)
	if has(ctx, b, POTION, 10) == false or has(ctx, a, POTION, 60) == false then
		return false
	end

	local chat = a:request_on(b, resp.mini_room_chat, req.mini_room { mode = MODE.chat, message = "교환 테스트" }, nil, 5000)
	if check(ctx, chat ~= false and chat.slot == 0 and chat.message:find("교환 테스트", 1, true) ~= nil, "교환창 채팅이 상대에게 오지 않음") == false then
		return false
	end

	if a:request_on(b, resp.trade_confirm, req.mini_room { mode = MODE.confirm }, nil, 5000) == false then
		return ctx:fail("교환 확정이 상대에게 보이지 않음")
	end
	local done = b:request_on(a, resp.mini_room_leave, req.mini_room { mode = MODE.confirm }, function(p)
		return p.reason == LEAVE.done
	end, 5000)
	if check(ctx, done ~= false and done.slot == 0, "양쪽 확정 후 교환이 성사되지 않음") == false then
		return false
	end
	ctx:sleep(500)
	if has(ctx, a, POTION, 60) == false or has(ctx, a, SWORD, 1) == false or meso_is(ctx, a, 1000000) == false then
		return false
	end
	return has(ctx, b, POTION, 40) and has(ctx, b, SWORD, 0) and meso_is(ctx, b, 982000)
end

local function cancel_flow(ctx)
	local a = ctx:bot(0)
	local b = ctx:bot(1)
	if setup(ctx, a, 50000, POTION .. ":100," .. SWORD .. ":1") == false or setup(ctx, b, 0, "-") == false then
		return false
	end
	local potion_slot = a:slot(POTION)
	local sword_slot = a:slot(SWORD)

	if open(ctx, a, b) == false then
		return false
	end
	if put_item(ctx, a, b, CONSUME, POTION, 40, 1) == false or put_item(ctx, a, b, EQUIP, SWORD, 1, 2) == false then
		return false
	end
	if put_meso(ctx, a, b, 5000, 5000) == false or put_meso(ctx, a, b, 1000, 6000) == false then
		return false
	end
	ctx:sleep(300)
	if has(ctx, a, POTION, 60) == false or has(ctx, a, SWORD, 0) == false or meso_is(ctx, a, 44000) == false then
		return false
	end

	local left = b:request_on(a, resp.mini_room_leave, req.mini_room { mode = MODE.exit }, nil, 5000)
	if check(ctx, left ~= false and left.slot == 0 and left.reason == LEAVE.cancel, "상대 취소 알림이 오지 않음") == false then
		return false
	end
	ctx:sleep(500)
	if has(ctx, a, POTION, 100) == false or has(ctx, a, SWORD, 1) == false or meso_is(ctx, a, 50000) == false then
		return false
	end
	return check(ctx, a:slot(POTION) == potion_slot and a:slot(SWORD) == sword_slot, "취소 후 원래 칸으로 돌아오지 않음")
end

-- An item put into a trade keeps its slot, so a loot cannot take the slot the item must return to.
local function restore_slot(ctx)
	local a = ctx:bot(0)
	local b = ctx:bot(1)
	local items = {}
	for i = 1, 32 do
		items[#items + 1] = (SHELL + i) .. ":1"
	end
	if setup(ctx, a, 0, SHELL .. ":1") == false or setup(ctx, b, 0, table.concat(items, ",")) == false then
		return false
	end
	local traded = SHELL + 1
	local slot = b:slot(traded)
	local dropped = a:drop(SHELL, 1)
	if check(ctx, dropped ~= nil, "줍기 확인용 아이템을 떨어뜨리지 못함") == false then
		return false
	end

	if open(ctx, a, b) == false or put_item(ctx, b, a, ETC, traded, 1, 1) == false then
		return false
	end
	if check(ctx, b:loot(dropped) == false, "교환에 올린 아이템이 돌아갈 칸으로 아이템을 주움") == false then
		return false
	end

	local left = b:request_on(a, resp.mini_room_leave, req.mini_room { mode = MODE.exit }, nil, 5000)
	if check(ctx, left ~= false and left.reason == LEAVE.cancel, "취소 알림이 오지 않음") == false then
		return false
	end
	ctx:sleep(500)
	if has(ctx, b, traded, 1) == false or has(ctx, b, SHELL, 0) == false then
		return false
	end
	return check(ctx, b:slot(traded) == slot, "취소 후 원래 칸으로 돌아오지 않음")
end

local function move_flow(ctx)
	local a = ctx:bot(0)
	local b = ctx:bot(1)
	if setup(ctx, a, 0, POTION .. ":100") == false or setup(ctx, b, 0, "-") == false then
		return false
	end

	if open(ctx, a, b) == false or put_item(ctx, a, b, CONSUME, POTION, 40, 1) == false then
		return false
	end
	local refused = a:request(resp.server_blocked, req.enter_cash_shop(), nil, 5000)
	if check(ctx, refused ~= false, "교환 중 캐시샵 이동이 거절되지 않음") == false then
		return false
	end

	local left = a:request_on(b, resp.mini_room_leave, req.mini_room { mode = MODE.exit }, nil, 5000)
	if check(ctx, left ~= false and left.reason == LEAVE.cancel, "이동 거절 뒤 교환 취소 알림이 오지 않음") == false then
		return false
	end
	ctx:sleep(500)
	return has(ctx, a, POTION, 100) and has(ctx, b, POTION, 0)
end

local function decline_flow(ctx)
	local a = ctx:bot(0)
	local b = ctx:bot(1)
	local sn = invite(ctx, a, b)
	if sn == false then
		return false
	end
	local result = b:request_on(a, resp.trade_invite_result, req.mini_room { mode = MODE.decline, sn = sn, reason = INVITE.declined }, nil, 5000)
	if check(ctx, result ~= false and result.result == INVITE.declined and result.name == b:name(), "거절 결과가 초대자에게 오지 않음") == false then
		return false
	end
	ctx:sleep(300)

	if invite(ctx, a, b) == false then
		return false
	end
	local left = a:request(resp.mini_room_leave, req.mini_room { mode = MODE.exit }, nil, 5000)
	return check(ctx, left ~= false and left.reason == LEAVE.exit, "초대 대기 중 나가기 응답이 다름")
end

local function invitee_moves(ctx)
	local a, b = ctx:bot(0), ctx:bot(1)
	if setup(ctx, a, 0, "-") == false or setup(ctx, b, 0, "-") == false then
		return false
	end
	local old_sn = invite(ctx, a, b, "이동 전 교환")
	if old_sn == false then
		return false
	end
	local left = b:request_on(a, resp.mini_room_leave, req.normal_chat { message = "/맵이동 100000000" }, nil, 5000)
	if check(ctx, left ~= false and left.reason == LEAVE.cancel, "초대받은 캐릭터 이동 뒤 교환이 취소되지 않음") == false then
		return false
	end
	ctx:sleep(500)
	if check(ctx, b:map() == 100000000, "초대받은 캐릭터 이동이 끝나지 않음") == false then
		return false
	end
	if b:map_move(LITH_HARBOR) == false then
		return ctx:fail("초대받은 캐릭터 복귀 실패")
	end
	ctx:sleep(300)
	local new_sn = invite(ctx, a, b, "복귀 후 교환")
	if new_sn == false then
		return false
	end
	if check(ctx, new_sn ~= old_sn, "새 초대의 SN이 이전과 같음") == false then
		return false
	end
	local rejected = b:request(resp.mini_room_enter_failed, req.mini_room { mode = MODE.visit, sn = old_sn }, nil, 5000)
	if check(ctx, rejected ~= false, "지난 초대 SN으로 입장됨") == false then
		return false
	end
	local closed = a:request(resp.mini_room_leave, req.mini_room { mode = MODE.exit }, nil, 5000)
	return check(ctx, closed ~= false and closed.reason == LEAVE.exit, "이동 후 다시 초대할 수 없음")
end

test_suite {
	name = "Trade: 개인 거래 성사·취소·복구 칸·이동 거절·거절·초대 중 이동",
	bot_count = 2,

	scenarios = {
		trade_flow,
		cancel_flow,
		restore_slot,
		move_flow,
		decline_flow,
		invitee_moves,
	},
}
