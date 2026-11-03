local MIN_LEVEL = 120
local SKILL_BOOK = 2280003

local M = {}

function M.advance(me, npc, cfg)
    if me:level() < MIN_LEVEL then
        me:dialog(npc, cfg.text.intro)
        return
    end
    local advanced = cfg.advanced[me:class()]
    if advanced ~= nil then
        me:dialog(npc, advanced)
        return
    end
    local q = me:quest(cfg.quest)
    if q == nil or not q:completed() then
        me:dialog(npc, cfg.text.intro)
        return
    end
    local advance = cfg.advance[me:class()]
    if advance == nil then
        return
    end

    local sel = me:dialog_list(npc, cfg.text.choose, {
        advance.option,
        "#r#h # : #b잠시만 생각할 시간이 필요합니다.",
    })
    if sel == nil then
        return
    end
    if sel == 2 then
        me:dialog(npc, cfg.text.wait)
        return
    end
    if me:skill_point() > me:level() * 3 then
        me:dialog(npc, cfg.text.sp_left)
        return
    end

    local code = me:exchange(nil, { item = { [SKILL_BOOK] = 1 } })
    if code ~= ExchangeResult.OK then
        me:dialog(npc, cfg.text.no_space)
        return
    end
    me:class(advance.class)
    if not me:dialog(npc, advance.done, false, true) then
        return
    end
    me:dialog(npc, cfg.text.farewell, false, true)
end

return M
