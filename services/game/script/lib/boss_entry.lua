local pq = require("script/lib/party_quest")
local bosses = require("script/lib/bosses")

local M = {}

local SCRIPT = "script/lib/boss_entry.lua"
local EXIT_MAP = 123256780
local MIN_LEVEL = 120
local MIN_MEMBERS = 1
local MAX_MEMBERS = 6
local BUSY = "이미 누군가가 이용 중에 있습니다."
local NEXT_PHASE_DELAY_MS = 5000

function M.spawn_boss(map, x, y, ...)
    for _, mob_id in ipairs({ ... }) do
        map:spawn_mob(mob_id, x, y)
    end
end

-- Sends characters into a phase: its mobs come out, and its clock sends everyone back to the waiting room when it
-- runs out.
local function enter_phase(characters, phase)
    run_on_map(phase.map, SCRIPT, "spawn_boss", phase.x, phase.y, phase.mobs[1], phase.mobs[2])
    local seconds = phase.time * 60
    for _, ch in pairs(characters) do
        local function on_arrive(ch)
            ch:clock(seconds, function(ch)
                ch:map(EXIT_MAP)
            end)
        end
        ch:map(phase.map, 0, { callback = on_arrive })
    end
end

local function enter(me, boss)
    local party = me:party()
    if party == nil then
        return "파티를 구성하고 있지 않으십니다."
    end
    if party:leader_id() ~= me:id() then
        return "파티장만이 저에게 말을 걸 수 있습니다."
    end

    local characters = me:map():characters()
    local size = 0
    local leveled = 0
    local in_map = 0
    local entries = 0
    for _, mem in pairs(party:members()) do
        size = size + 1
        if mem:level() >= MIN_LEVEL then
            leveled = leveled + 1
        end
        local ch = characters[mem:id()]
        if ch ~= nil then
            in_map = in_map + 1
            if ch:records():get(boss.record) < boss.limit then
                entries = entries + 1
            end
        end
    end
    if pq.is_gm(me) == false then
        if size > MAX_MEMBERS or size < MIN_MEMBERS or in_map < size or leveled < size then
            return "파티원 중 레벨이 부족한 사람이 있거나, 파티원 수가 안맞는거 같습니다. 입장 조건을 다시 확인하여주세요."
        end
        if entries < size then
            return "파티원 중 오늘 이미 " .. boss.limit .. "번 입장하신 분이 있습니다."
        end
    end

    local pid = party:id()
    local members = {}
    for _, ch in pairs(characters) do
        local p = ch:party()
        if p ~= nil and p:id() == pid then
            ch:records():add(boss.record, 1, { daily = true })
            members[#members + 1] = ch
        end
    end
    enter_phase(members, boss.phases[1])
    return nil
end

-- A phase's boss died on map. The world hears of it, and 5 seconds later everyone still there goes on to the next
-- phase, as the original did. The move runs on the boss map's own actor, which mob scripts call into.
function M.on_phase_cleared(map, attacker, key)
    local leader = ""
    local party = attacker ~= nil and attacker:party() or nil
    if party ~= nil then
        for _, mem in pairs(party:members()) do
            if mem:id() == party:leader_id() then
                leader = mem:name()
            end
        end
    end
    run_on_map(map:wz():id(), SCRIPT, "clear_phase", key, leader)
end

function M.clear_phase(map, key, leader)
    local boss = bosses[key]
    local index = nil
    for i, phase in ipairs(boss.phases) do
        if phase.map == map:wz():id() then
            index = i
        end
    end
    if index == nil then
        return
    end
    local next_phase = boss.phases[index + 1]
    if leader ~= "" then
        local cleared = boss.name .. " " .. index .. "페이즈를 클리어하였습니다."
        if next_phase == nil then
            cleared = boss.victory
        end
        world_notice("[헤이븐] : " .. leader .. "님의 파티가 " .. cleared, Msg.Megaphone)
    end
    if next_phase == nil then
        return
    end

    sleep(NEXT_PHASE_DELAY_MS)
    local characters = {}
    for _, ch in pairs(map:characters()) do
        characters[#characters + 1] = ch
    end
    if #characters > 0 then
        enter_phase(characters, next_phase)
    end
end

function M.on_click(me, npc, boss)
    local party = me:party()
    if party == nil then
        me:dialog(npc, "파티를 구성하여 주세요.")
        return
    end
    if me:map():wz():id() ~= EXIT_MAP then
        me:dialog(npc, "#m" .. EXIT_MAP .. "#맵에서만 저에게 말을 걸어 입장하실 수 있습니다.")
        return
    end
    if party:leader_id() ~= me:id() then
        me:dialog(npc, "파티장만이 제게 말을 거실 수 있습니다.")
        return
    end
    local text = "#e< #b" .. boss.name .. " 입장 조건#k >#n#b\r\n"
        .. "- #e#r난이도 #k: " .. boss.difficulty .. "\r\n"
        .. "- 120레벨 이상의 파티원\r\n"
        .. "- 파티원 모두 같은맵에 존재\r\n"
        .. "- 최소 1인부터 최대 6인의 파티\r\n\r\n"
        .. "#k"
        .. "위의 조건을 모두 충족하셨다면, 수락하기를 눌러주세요.\r\n"
    if me:dialog_accept(npc, text) == false then
        return
    end

    local maps = {}
    for _, phase in ipairs(boss.phases) do
        maps[#maps + 1] = phase.map
    end
    local refused = pq.enter_rooms(me:map(), "boss_entry:" .. maps[1], { maps }, BUSY, function()
        return enter(me, boss)
    end)
    if refused ~= nil then
        me:dialog(npc, refused)
    end
end

return M
