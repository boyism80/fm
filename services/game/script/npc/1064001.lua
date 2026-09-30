-- NPC name (String.wz/Npc.img.xml): 의문의 소녀

local SCRIPT = "script/npc/1064001.lua"
local TICKET = 3980069
local TICKET_PRICE = 100000000
local MAX_MEMBERS = 4
local FIRST_MOB = 8910000
local WAVES = { 123456772, 123456773, 123456774, 123456775, 123456776 }

local function challenge(me, npc)
	local party = me:party()
	if party == nil then
		me:dialog(npc, "파티를 구성하여주세요.")
		return
	end
	if party:leader_id() ~= me:id() then
		me:dialog(npc, "파티장에게 말을 걸어달라고 해주세요.")
		return
	end
	local characters = me:map():characters()
	local size = 0
	local in_map = 0
	for _, mem in pairs(party:members()) do
		size = size + 1
		if characters[mem:id()] ~= nil then
			if mem:class_id() == Class.GM then
				in_map = in_map + 2
			else
				in_map = in_map + 1
			end
		end
	end
	if size > MAX_MEMBERS or in_map < MAX_MEMBERS then
		local m = "입장 조건을 충족하지 못하셨습니다.\r\n\r\n"
			.. "#e[ #r입장조건#k ]#n\r\n\r\n"
			.. "- 파티원 수 : 4명\r\n"
			.. "- 입장권 수 : 4개\r\n"
			.. "- 티어 등급 : 플래티넘\r\n"
			.. "- 파티원 맵 : 파티장과 같은 맵\r\n"
		me:dialog(npc, m)
		return
	end

	for _, map_id in ipairs(WAVES) do
		local ok, count = run_on_map(map_id, SCRIPT, "character_count")
		if ok and count ~= nil and count > 0 then
			me:dialog(npc, "이미 다른 파티가 입장 중이거나, 파티원의 수만큼 입장권을 가지고 계시지 않습니다.")
			return
		end
	end
	if me:exchange({ item = { [TICKET] = size } }, nil) ~= ExchangeResult.OK then
		me:dialog(npc, "파티원 중 입장권이 부족한 사람이 있습니다.")
		return
	end
	for _, map_id in ipairs(WAVES) do
		run_on_map(map_id, SCRIPT, "reset_map")
	end
	run_on_map(WAVES[1], SCRIPT, "spawn_first", FIRST_MOB, -63, 112)
	local pid = party:id()
	for _, ch in pairs(characters) do
		local p = ch:party()
		if p ~= nil and p:id() == pid then
			ch:map(WAVES[1], 0)
		end
	end
end

local function buy_ticket(me, npc)
	local max = math.floor(me:meso() / TICKET_PRICE)
	local input = me:dialog_input(npc, "#z" .. TICKET .. "#은 개당 #b" .. TICKET_PRICE .. " 메소#k입니다.\r\n입장권을 총 몇개 구매하시겠습니까?\r\n\r\n")
	local n = tonumber(input)
	if n == nil or n < 1 or n > max or n > 25000 then
		return
	end
	n = math.floor(n)
	if me:exchange({ meso = TICKET_PRICE * n }, { item = { [TICKET] = n } }) ~= ExchangeResult.OK then
		me:dialog(npc, "메소가 부족하시거나, 값을 조작하였습니다.")
		return
	end
	me:dialog(npc, "정상적으로 #z" .. TICKET .. "# " .. n .. "개를 구매하셨습니다.")
end

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	reset_map = function(map)
		map:reset()
	end,

	spawn_first = function(map, mob_id, x, y)
		map:spawn_mob(mob_id, x, y)
	end,

	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "#r#e더 시드#n#k에 오신 것을 환영합니다. 어떤 것을 도와드릴까요?#b", {
			"더 시드에 도전하겠습니다.",
			"더 시드의 입장권을 구매하겠습니다.",
		})
		if sel == nil then
			return
		end
		if sel == 1 then
			challenge(me, npc)
		elseif sel == 2 then
			buy_ticket(me, npc)
		end
	end
}
