local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
local SNOWMAN = 8220001
local POISON_MIST = 131
-- Poison mist level 7 reaches lt (-50, -160) to rb (300, 40).
local MIST_AREA = { left = -50, top = -160, right = 300, bottom = 40 }
local CAST_ACTION = 21
local FACE_RIGHT = 4
-- Weapon attack up for every mob within (-300, -200) to (300, 200).
local POWER_UP_MOB = 5120506
local POWER_UP = 110
-- Darkness within (-300, -120) to (300, 120).
local DARKNESS_MOB = 9501001
local DARKNESS = 121

local function setup(ctx, bot)
    if pq.command(bot, "/봇초기화 70 0 0 -", "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    if pq.command(bot, "/체력바꾸기 30000", "체력 설정") == false then
        return ctx:fail(bot:name() .. " 체력 설정 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail(bot:name() .. " 사냥터 이동 실패")
    end
    return true
end

local function spawn(ctx, bot, mob_id)
    local p = bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. mob_id }, function(p)
        return p.mob.mob_id == mob_id
    end, 5000)
    if p == false then
        return ctx:fail(mob_id .. " 생성 응답 없음")
    end
    return p.mob
end

-- The controller asks for a skill until the server gives skill_id, then casts it as the client does.
local function cast(ctx, controller, oid, skill_id, response, filter)
    local level = nil
    for _ = 1, 5 do
        local chosen = controller:request(resp.control_move_mob,
            req.move_mob { oid = oid, active_skill = true, center_split = -1 },
            function(p)
                return p.oid == oid
            end, 3000)
        if chosen == false then
            return ctx:fail("몹 스킬 선택 응답이 오지 않음")
        end
        if chosen.skill_id == skill_id then
            level = chosen.skill_level
            break
        end
    end
    if level == nil then
        return ctx:fail("몹이 스킬 " .. skill_id .. "을 고르지 않음")
    end
    local p = controller:request(response,
        req.move_mob { oid = oid, center_split = CAST_ACTION * 2, skill_id = skill_id, skill_level = level }, filter,
        3000)
    if p == false then
        return ctx:fail("몹 스킬 " .. skill_id .. " 시전 결과가 오지 않음")
    end
    return p
end

-- The poison mist is laid as WZ describes it, facing left, whichever way the mob faces, as in the original.
local function poison_mist(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot) == false then
        return false
    end
    local mob = spawn(ctx, bot, SNOWMAN)
    if mob == false then
        return false
    end
    local x, y = mob.position.x, mob.position.y
    if bot:move_mob(mob.oid, x, y, FACE_RIGHT) == false then
        return ctx:fail("눈사람을 오른쪽으로 돌리지 못함")
    end

    local mist = cast(ctx, bot, mob.oid, POISON_MIST, resp.spawn_mist, function(p)
        return p.skill_id == POISON_MIST
    end)
    pq.command(bot, "/몬스터죽이기", "마리 제거")
    if mist == false then
        return false
    end
    local b = mist.bounds
    if b.left ~= x + MIST_AREA.left or b.top ~= y + MIST_AREA.top
        or b.right ~= x + MIST_AREA.right or b.bottom ~= y + MIST_AREA.bottom then
        return ctx:fail(string.format("독 안개 범위가 다름: (%d, %d) (%d, %d) 몹 (%d, %d)",
            b.left, b.top, b.right, b.bottom, x, y))
    end
    return true
end

-- An area buff reaches the other mobs in the area besides the caster.
local function area_buff(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot) == false then
        return false
    end
    local caster = spawn(ctx, bot, POWER_UP_MOB)
    local other = caster ~= false and spawn(ctx, bot, POWER_UP_MOB)
    if other == false then
        return false
    end

    local buffed = {}
    ctx:hook("area_buff", function(_, _, p, name)
        if name == resp.apply_mob_buff then
            buffed[p.oid] = true
        end
    end)
    local cast_ok = cast(ctx, bot, caster.oid, POWER_UP, resp.apply_mob_buff, function(p)
        return p.oid == caster.oid
    end)
    ctx:sleep(500)
    ctx:unhook("area_buff")
    pq.command(bot, "/몬스터죽이기", "마리 제거")
    if cast_ok == false then
        return false
    end
    if buffed[other.oid] ~= true then
        return ctx:fail("범위 안 다른 몹에 버프가 걸리지 않음")
    end
    return true
end

-- An area disease reaches only the controller standing in the area, as in the original.
local function area_disease(ctx)
    local first, second = ctx:bot(0), ctx:bot(1)
    if setup(ctx, first) == false or setup(ctx, second) == false then
        return false
    end
    local controllers = {}
    ctx:hook("disease_control", function(_, b, p, name)
        if name == resp.start_control_mob then
            controllers[p.mob.oid] = b
        end
    end)
    local mob = spawn(ctx, first, DARKNESS_MOB)
    if mob == false then
        return false
    end
    ctx:sleep(500)
    ctx:unhook("disease_control")
    local controller = controllers[mob.oid]
    if controller == nil then
        return ctx:fail("생성한 몹의 컨트롤이 오지 않음")
    end
    -- The mob lands on the ground below the spawn point.
    for _, bot in ipairs({ first, second }) do
        if pq.move(bot, mob.position.x, mob.position.y) == false then
            return ctx:fail(bot:name() .. " 몹 옆으로 이동 실패")
        end
    end

    local diseased = {}
    ctx:hook("area_disease", function(_, _, p, name)
        if name == resp.give_debuff and p.skill_id == DARKNESS then
            diseased[p.character_id] = true
        end
    end)
    local cast_ok = cast(ctx, controller, mob.oid, DARKNESS, resp.control_move_mob, function(p)
        return p.oid == mob.oid
    end)
    ctx:sleep(1000)
    ctx:unhook("area_disease")
    pq.command(first, "/몬스터죽이기", "마리 제거")
    if cast_ok == false then
        return false
    end
    if diseased[controller:id()] ~= true then
        return ctx:fail("범위 안 컨트롤러에게 암흑이 걸리지 않음")
    end
    local other = first
    if controller:id() == first:id() then
        other = second
    end
    if diseased[other:id()] == true then
        return ctx:fail("컨트롤러가 아닌 캐릭터에게도 암흑이 걸림")
    end
    return true
end

test_suite {
    name = "MobSkill: 독 안개·범위 버프·범위 상태 이상",
    bot_count = 2,

    scenarios = {
        poison_mist,
        area_buff,
        area_disease,
    },
}
