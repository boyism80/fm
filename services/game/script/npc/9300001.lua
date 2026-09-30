-- NPC name (String.wz/Npc.img.xml): 사장법사

local SCRIPT = "script/npc/9300001.lua"
local DEST = 950101100
local FLAME = 4001433
local FLAME_COST = 30
local MIN_LEVEL = 50
local MAX_LEVEL = 190

local function busy_maps()
	local maps = { DEST }
	for id = 970040200, 970042000, 100 do
		maps[#maps + 1] = id
	end
	return maps
end

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	prepare_stage = function(map)
		map:kill_all_mobs()
		map:spawn_mob(9500390, 69, 513)
	end,

	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "나무관세음보살.. 황금사원에 큰 위험이 닥치고 있다네 라바나를 죽이지 않으면 130년 전통의 황금사원이 무너지고 말께야.\r\n#b", {
			"이지 모드 (레벨 50~190",
		})
		if sel ~= 1 then
			return
		end
		local party = me:party()
		if party == nil then
			me:dialog(npc, "파티를 만들고 말을 걸어주세요.")
			return
		end
		if party:leader_id() ~= me:id() then
			me:dialog(npc, "파티장이 저에게 말을 걸어주세요.")
			return
		end
		local characters = me:map():characters()
		local members = party:members()
		for _, mem in pairs(members) do
			if characters[mem:id()] == nil then
				me:dialog(npc, "파티를 구성한게 맞는지 혹은 태양의 불꽃 30개를 소지하고 계신게 맞는지 확인 해주세요.")
				return
			end
		end
		for _, map_id in ipairs(busy_maps()) do
			local ok, count = run_on_map(map_id, SCRIPT, "character_count")
			if ok and count ~= nil and count > 0 then
				me:dialog(npc, "이 안에 이미 다른 파티가 입장하여 클리어에 도전중입니다. 잠시 후에 다시 시도해보세요.")
				return
			end
		end
		for _, mem in pairs(members) do
			if mem:level() < MIN_LEVEL or mem:level() > MAX_LEVEL then
				me:dialog(npc, "파티원 중에서 레벨 #r50~90#k 에 해당하지 않는 캐릭터가 있습니다.")
				return
			end
		end

		if me:exchange({ item = { [FLAME] = FLAME_COST } }, nil) ~= ExchangeResult.OK then
			me:dialog(npc, "파티를 구성한게 맞는지 혹은 태양의 불꽃 30개를 소지하고 계신게 맞는지 확인 해주세요.")
			return
		end
		run_on_map(DEST, SCRIPT, "prepare_stage")
		local pid = party:id()
		for _, ch in pairs(characters) do
			local p = ch:party()
			if p ~= nil and p:id() == pid then
				ch:map(DEST, 0)
			end
		end
	end
}
