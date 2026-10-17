local pq = require("script/integration/lib/party_quest")

local MODE = { add = 1, accept = 2 }
local SYNC = { update = 10 }
local GROUP = "봇친구"

local lists = {}
local requests = {}

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function listen(ctx)
	ctx:hook("buddy", function(_, bot, p, name)
		if name == "buddy_list_update" and p.action == SYNC.update then
			lists[bot:id()] = p.entries
		elseif name == "buddy_add_request" then
			requests[bot:id()] = p
		end
	end)
end

local function find(entries, id)
	for _, entry in ipairs(entries or {}) do
		if entry.character_id == id then
			return entry
		end
	end
	return nil
end

local function wait_list(ctx, bot, ids, what)
	for _ = 1, 50 do
		local entries = lists[bot:id()]
		if entries ~= nil and #entries == #ids then
			local all = true
			for _, id in ipairs(ids) do
				if find(entries, id) == nil then
					all = false
				end
			end
			if all then
				return entries
			end
		end
		ctx:sleep(100)
	end
	return ctx:fail(string.format("%s %s: 친구 목록 %d명 (기대 %d명)", bot:name(), what, #(lists[bot:id()] or {}), #ids))
end

local function wait_request(ctx, bot, from, what)
	for _ = 1, 50 do
		local p = requests[bot:id()]
		if p ~= nil and p.from_character_id == from:id() then
			return p
		end
		ctx:sleep(100)
	end
	return ctx:fail(string.format("%s %s: 친구 신청이 오지 않음", bot:name(), what))
end

local function add(bot, target)
	lists[bot:id()] = nil
	requests[target:id()] = nil
	return bot:send(req.buddy { mode = MODE.add, name = target:name(), group = GROUP })
end

-- Adding or accepting a friend sends the whole list, so earlier friends stay on the client.
local function list_flow(ctx)
	local a, b, c = ctx:bot(0), ctx:bot(1), ctx:bot(2)
	for _, bot in ipairs({ a, b, c }) do
		if pq.command(bot, "/봇초기화 30 0 0 -", "봇초기화 완료") == false then
			return ctx:fail(bot:name() .. " 봇 초기화 실패")
		end
	end
	listen(ctx)

	add(a, b)
	if wait_list(ctx, a, { b:id() }, "첫 친구 추가") == false or wait_request(ctx, b, a, "첫 친구 추가") == false then
		return false
	end

	add(a, c)
	if wait_list(ctx, a, { b:id(), c:id() }, "두 번째 친구 추가") == false or wait_request(ctx, c, a, "두 번째 친구 추가") == false then
		return false
	end

	lists[a:id()] = nil
	lists[b:id()] = nil
	b:send(req.buddy { mode = MODE.accept, character_id = a:id() })
	local mine = wait_list(ctx, b, { a:id() }, "친구 수락")
	if mine == false then
		return false
	end
	if check(ctx, find(mine, a:id()).pending == false, "수락한 친구가 대기 중으로 남음") == false then
		return false
	end
	local theirs = wait_list(ctx, a, { b:id(), c:id() }, "상대가 수락한 뒤")
	if theirs == false then
		return false
	end
	return check(ctx, find(theirs, b:id()).pending == false, "상대가 수락했는데 대기 중으로 남음")
end

test_suite {
	name = "Buddy: 친구 추가·수락 시 전체 목록",
	bot_count = 3,

	scenarios = {
		list_flow,
	},
}
