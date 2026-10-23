local SUCCESS = 0
local INVALID = 0x14

local function relog(ctx, bot, opts, label)
    local result = bot:relog(opts)
    if result == false then
        ctx:fail(label .. ": 재접속 실패")
        return nil
    end
    return result
end

local function expect(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function set_select_remove(ctx)
    local bot = ctx:bot(0)

    local r = relog(ctx, bot, { set = "12", passwords = { "" } }, "짧은 비밀번호 설정")
    if r == nil then
        return false
    end
    if expect(ctx, r.has_second_password == false and r.set == INVALID and r.rejected == 0,
            string.format("짧은 2차 비밀번호가 설정됨: has=%s set=%s rejected=%d", tostring(r.has_second_password), tostring(r.set), r.rejected)) == false then
        return false
    end

    r = relog(ctx, bot, { set = "1234", delete_spare = true, passwords = { "", "9999", "1234" } }, "비밀번호 설정")
    if r == nil then
        return false
    end
    if expect(ctx, r.set == SUCCESS and r.rejected == 2,
            string.format("2차 비밀번호 설정 후 선택 결과가 다름: set=%s rejected=%d", tostring(r.set), r.rejected)) == false then
        return false
    end
    if expect(ctx, r.deleted == false and r.delete_notice ~= "",
            string.format("2차 비밀번호가 있는데 캐릭터가 삭제됨: deleted=%s notice=%q", tostring(r.deleted), r.delete_notice)) == false then
        return false
    end

    r = relog(ctx, bot, { set = "5678", passwords = { "1234" } }, "비밀번호 재설정")
    if r == nil then
        return false
    end
    if expect(ctx, r.has_second_password == true and r.set == INVALID and r.rejected == 0,
            string.format("이미 있는 2차 비밀번호를 덮어씀: has=%s set=%s rejected=%d", tostring(r.has_second_password), tostring(r.set), r.rejected)) == false then
        return false
    end

    r = relog(ctx, bot, { remove = "9999", passwords = { "1234" } }, "틀린 비밀번호로 해제")
    if r == nil then
        return false
    end
    if expect(ctx, r.remove == INVALID and r.rejected == 0,
            string.format("틀린 비밀번호로 2차 비밀번호가 해제됨: remove=%s rejected=%d", tostring(r.remove), r.rejected)) == false then
        return false
    end

    r = relog(ctx, bot, { remove = "1234", passwords = { "" } }, "비밀번호 해제")
    if r == nil then
        return false
    end
    if expect(ctx, r.remove == SUCCESS and r.rejected == 0,
            string.format("2차 비밀번호 해제 후 바로 선택되지 않음: remove=%s rejected=%d", tostring(r.remove), r.rejected)) == false then
        return false
    end

    r = relog(ctx, bot, { delete_spare = true, passwords = { "1234", "" } }, "해제 후 재접속")
    if r == nil then
        return false
    end
    return expect(ctx, r.has_second_password == false and r.rejected == 1 and r.deleted == true,
        string.format("해제 후 상태가 다름: has=%s rejected=%d deleted=%s", tostring(r.has_second_password), r.rejected,
            tostring(r.deleted)))
end

test_suite {
    name = "SecondPassword: 설정·선택·해제",
    bot_count = 1,

    scenarios = {
        set_select_remove,
    },
}
