local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
local SNOWMAN = 8220001
local POISON_MIST = 131
-- Poison mist level 7 reaches lt (-50, -160) to rb (300, 40).
local MIST_AREA = { left = -50, top = -160, right = 300, bottom = 40 }
local CAST_ACTION = 21
local FACE_RIGHT = 4

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

test_suite {
    name = "MobSkill: 독 안개",
    bot_count = 1,

    scenarios = {
        poison_mist,
    },
}
