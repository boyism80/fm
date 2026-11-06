local pq = require("script/lib/party_quest")

local PQ_ITEMS = { 4001130, 4001131, 4001132, 4001133, 4001134, 4001135 }
local RANKING_QUEST = 1205
local DURATION_MS = 45 * 60000
local SCALE_LEVEL = 200
local MIN_PARTY = 4
local MIN_LEVEL = 71

local function strip_items(me)
    for _, id in ipairs(PQ_ITEMS) do
        pq.remove_all(id, me)
    end
end

local function shuffle_lab(map)
    if map == nil then
        return
    end
    local reactors = {}
    local positions = {}
    for _, reactor in pairs(map:reactors()) do
        local name = reactor:name()
        if name == nil or not string.find(name, "out", 1, true) then
            local x, y = reactor:position()
            reactors[#reactors + 1] = reactor
            positions[#positions + 1] = { x, y }
        end
    end
    for i = #positions, 2, -1 do
        local j = math.random(i)
        positions[i], positions[j] = positions[j], positions[i]
    end
    for i, reactor in ipairs(reactors) do
        reactor:position(positions[i][1], positions[i][2])
    end
end

local function hit_all_reactors(map)
    if map == nil then
        return
    end
    for _, reactor in pairs(map:reactors()) do
        if reactor ~= nil and reactor:state() == 0 then
            reactor:hit(1)
        end
    end
end

local function clear_fx(map)
    if map ~= nil then
        map:clear_effect()
    end
end

local function init_stage6(sm)
    for y = 0, 3 do
        sm:set_property("stage6_" .. tostring(y), "0")
        for b = 0, 9 do
            for c = 0, 3 do
                sm:set_property(string.format("stage6_%d_%d_%d", y, b, c), "0")
            end
        end
    end
    for y = 0, 3 do
        for i = 0, 9 do
            local found = false
            while not found do
                for x = 0, 3 do
                    if found then
                        break
                    end
                    local taken = false
                    for z = 0, 3 do
                        if sm:get_property(string.format("stage6_%d_%d_%d", z, i, x)) == "1" then
                            taken = true
                            break
                        end
                    end
                    if not taken and math.random() < 0.25 then
                        sm:set_property(string.format("stage6_%d_%d_%d", y, i, x), "1")
                        found = true
                    end
                end
            end
        end
    end
end

local function pick_investigate(sm, map, npc_id)
    local oids = {}
    for _, n in pairs(map:npcs()) do
        if n ~= nil and n:id() == npc_id then
            oids[#oids + 1] = n:oid()
        end
    end
    if #oids < 2 then
        return
    end
    local a = oids[math.random(#oids)]
    local b = a
    while b == a do
        b = oids[math.random(#oids)]
    end
    sm:set_property("stage1_" .. tostring(a), "2")
    sm:set_property("stage1_" .. tostring(b), "3")
end

local function create(cfg)
    local maps = cfg.maps
    local start_map = cfg.start_map
    local exit_map = cfg.exit_map
    local town_map = cfg.town_map
    local investigate_npc = cfg.investigate_npc
    local letter_item = cfg.letter_item
    local protect_mob = cfg.protect_mob
    local protect_fail_msg = cfg.protect_fail_msg or "보호에 실패했습니다."
    local protect_ok_msg = cfg.protect_ok_msg or "보호에 성공했습니다."
    local marble = cfg.marble
    local door_prefix = cfg.door_prefix
    local stage1_exp = cfg.stage1_exp
    local stage4_exp = cfg.stage4_exp
    local urete_office = cfg.urete_office
    local boss_map = cfg.boss_map
    local urete_map = cfg.urete_map
    local reward_map = cfg.reward_map
    local hallway = cfg.hallway
    local beaker_map = cfg.beaker_map
    local hub_map = cfg.hub_map
    local lab1 = cfg.lab1
    local lab2 = cfg.lab2
    local ending_npcs = cfg.ending_npcs

    local function end_run(sm)
        sm:finish(exit_map)
    end

    local function clear_props(sm)
        sm:set_property("stage", "0")
        sm:set_property("stage1", "0")
        sm:set_property("stage1_way_clear", "")
        sm:set_property("stage3", "0")
        sm:set_property("stage4", "0")
        sm:set_property("stage5", "0")
        sm:set_property("stage7", "0")
        sm:set_property("summoned", "0")
        sm:set_property("clear_protect", "")
        sm:set_property("urete_boss", "")
        sm:set_property("persuade_urete", "")
    end

    return {
        config = cfg,
        strip_items = strip_items,
        letter_item = letter_item,
        marble = marble,
        door_prefix = door_prefix,
        stage1_exp = stage1_exp,
        stage4_exp = stage4_exp,
        town_map = town_map,
        exit_map = exit_map,
        boss_map = boss_map,
        urete_map = urete_map,
        reward_map = reward_map,
        hub_map = hub_map,
        beaker_map = beaker_map,
        hallway = hallway,
        hit_all_reactors = hit_all_reactors,
        clear_fx = clear_fx,

        on_init = function(group)
            group:min_players(1)
            group:exit_map(exit_map)
        end,

        on_create = function(sm)
            local group = sm:group()
            clear_props(sm)
            init_stage6(sm)
            return maps
        end,

        on_prepare = function(sm)
            for _, map_id in ipairs(maps) do
                sm:map(map_id):respawn({ include_one_time = true })
            end
            pick_investigate(sm, sm:map(start_map), investigate_npc)
            shuffle_lab(sm:map(lab1))
            shuffle_lab(sm:map(lab2))
            sm:map(urete_office):spawn_npc(2112000, 200, 188)
            local boss = sm:map(boss_map)
            boss:set_respawn(false)
            boss:kill_all_mobs()
        end,

        on_start = function(sm)
            sm:start_timer(DURATION_MS)
        end,

        on_player_enter = function(sm, player)
            local function on_arrive(player)
                strip_items(player)
                player:try_party_quest(RANKING_QUEST)
            end
            player:map(sm:map(start_map), { callback = on_arrive })
        end,

        on_changed_map = function(sm, player, map_id)
            if map_id == boss_map and sm:get_property("urete_boss") == "" then
                local boss = sm:map(boss_map)
                if boss ~= nil then
                    boss:spawn_npc(2112010, 242, 150)
                    sm:set_property("urete_boss", "1")
                end
            end
        end,

        on_mob_die = function(sm, mob)
            if mob == nil then
                return
            end
            local id = mob:id()
            local hall = sm:map(hallway)
            if hall ~= nil and pq.mob_count(hall) == 0 and sm:get_property("stage1_way_clear") == "" then
                clear_fx(hall)
                sm:set_property("stage1_way_clear", "1")
            end
            if id >= 9300142 and id <= 9300146 then
                local office = sm:map(urete_office)
                if office ~= nil and pq.mob_count(office) == 0 and sm:get_property("stage5") == "1" then
                    sm:set_property("stage5", "2")
                    clear_fx(office)
                    hit_all_reactors(office)
                end
                return
            end
            if id == protect_mob then
                if sm:get_property("clear_protect") == "" then
                    sm:set_property("stage7", "1")
                    sm:message(protect_fail_msg, Msg.PinkText)
                end
                return
            end
            if id == 9300139 or id == 9300140 then
                sm:set_property("clear_protect", "1")
                local boss = sm:map(boss_map)
                if boss ~= nil then
                    clear_fx(boss)
                    boss:set_respawn(false)
                    boss:kill_all_mobs()
                end
                local ok = sm:get_property("stage7") ~= "1"
                if ok then
                    sm:message(protect_ok_msg)
                    pq.party_exp(sm, 45000, boss)
                else
                    pq.party_exp(sm, 30000, boss)
                end
                local spawns = ok and ending_npcs.success or ending_npcs.fail
                for _, map_id in ipairs({ urete_map, reward_map, boss_map }) do
                    local map = sm:map(map_id)
                    if map ~= nil then
                        for _, npc in ipairs(spawns[map_id]) do
                            map:spawn_npc(npc[1], npc[2], npc[3])
                        end
                    end
                end
            end
        end,

        on_left_party = function(sm, player)
            end_run(sm)
        end,

        on_disband_party = function(sm)
            end_run(sm)
        end,

        on_scheduled_timeout = function(sm)
            end_run(sm)
        end,

        on_finish = function(sm)
            clear_props(sm)
        end
    }
end

local function exchange_marbles(me, npc)
    local sel = me:dialog_list(npc,
        "제뉴미스트의 구슬과 알카드노의 구슬을 각각 25개씩 가져오면 호루스의 눈을, 둘 중 어느쪽의 구슬이라도 10개 가져오면 목걸이에 새로운 힘을 부여할 수 있는 지혜의 돌을 만들어 주겠네. 자, 어떤 아이템을 만들겠는가?\r\n#b",
        {
            "호루스의 눈을 만들어 주세요.",
            "제뉴미스트 구슬로 지혜의 돌을 만들어 주세요.",
            "알카드노 구슬로 지혜의 돌을 만들어 주세요.",
        })
    if sel == nil then
        return false
    end
    local material
    local reward_item
    if sel == 1 then
        material = { [4001159] = 25, [4001160] = 25 }
        reward_item = 1122010
    elseif sel == 2 then
        material = { [4001159] = 10 }
        reward_item = 2041212
    elseif sel == 3 then
        material = { [4001160] = 10 }
        reward_item = 2041212
    else
        return false
    end
    local str = "교환할 아이템을 확인하게.\r\n\r\n"
    if sel == 1 then
        str = str .. "#i4001159# #t4001159# 25개\r\n"
        str = str .. "#i4001160# #t4001160# 25개\r\n"
    elseif sel == 2 then
        str = str .. "#i4001159# #t4001159# 10개\r\n"
    else
        str = str .. "#i4001160# #t4001160# 10개\r\n"
    end
    str = str .. "\r\n위 아이템으로 \r\n"
    str = str .. string.format("#i%d# #t%d# 1개로 교환하겠네. 계속하겠는가?", reward_item, reward_item)
    if not me:dialog_yes_no(npc, str) then
        me:dialog(npc, "흐음, 잘못 선택한건가? 마음이 바뀌면 다시 오게. 아마 오래 만나지 못할게야.", false, false)
        return false
    end
    local code = me:exchange({ item = material }, { item = { [reward_item] = 1 } })
    if code == ExchangeResult.LackCapacity then
        me:dialog(npc, "인벤토리 공간이 부족한건 아닌지, 혹은 재료를 분명 제대로 갖고 계신건지 확인해 주시게.", false, false)
        return false
    end
    if code ~= ExchangeResult.OK then
        me:dialog(npc, "인벤토리 공간이 부족한건 아닌지, 혹은 재료를 분명 제대로 갖고 계신건지 확인해 주시게.", false, false)
        return false
    end
    if reward_item == 1122010 then
        local q1205 = me:quest(1205)
        if q1205 ~= nil and q1205:record_ex("have") == nil then
            q1205:record_ex("have", "1")
        end
    end
    return true
end

local M = {
    create = create,
    strip_items = strip_items,
    pq_items = PQ_ITEMS,
    ranking_quest = RANKING_QUEST,
    min_party = MIN_PARTY,
    min_level = MIN_LEVEL,
    scale_level = SCALE_LEVEL,
    hit_all_reactors = hit_all_reactors,
    clear_fx = clear_fx,
    exchange_marbles = exchange_marbles,
}

function M.path_for_map(map_id)
    if map_id >= 926100000 and map_id <= 926100700 then
        return "romeo"
    end
    if map_id >= 926110000 and map_id <= 926110700 then
        return "juliet"
    end
    return nil
end

function M.group_name(map_id)
    local path = M.path_for_map(map_id)
    if path == "romeo" then
        return "romeo_party_quest"
    end
    if path == "juliet" then
        return "juliet_party_quest"
    end
    return nil
end

local function try_start(me, npc, cfg)
    strip_items(me)
    local party = me:party()
    if party == nil or party:leader_id() ~= me:id() then
        me:dialog(npc, "파티가 없으시거나, 혹은 파티장이 아니신 건 아닌가요? 또한 4명의 파티만 입장할 수 있습니다.")
        return
    end
    local map = me:map()
    if map == nil or map:wz() == nil then
        return
    end
    local ok = true
    local size = 0
    local has_gm = false
    for _, mem in ipairs(party:members()) do
        if mem ~= nil then
            local ch = map:characters()[mem:id()]
            if ch == nil or mem:level() < MIN_LEVEL then
                ok = false
                break
            end
            if pq.is_gm(ch) then
                has_gm = true
            end
            size = size + 1
        end
    end
    if not ok or (size ~= MIN_PARTY and not has_gm) then
        me:dialog(npc, "파티원 중 레벨 제한 71이상이 아닌 파티원이 있거나 파티원이 4명이 아닌 것 같군요.")
        return
    end
    local group = state_machine(cfg.group)
    if group == nil then
        me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
        return
    end
    local sm, err = group:start_party(me, party, SCALE_LEVEL)
    if sm == nil then
        if err == StateMachineError.Full then
            me:dialog(npc, "이미 이 안에서 다른 파티가 퀘스트에 도전하는 중입니다. 잠시 후 다시 시도해 주세요.")
        else
            me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
            log(cfg.group .. " start_party:", err)
        end
    end
end

local function read_letter(me, npc, sm, cfg)
    if not pq.has_item(me, cfg.letter, 1) then
        return false
    end
    pq.remove_all(cfg.letter, me)
    sm:set_property("stage", "1")
    sm:message(cfg.letter_notice)
    me:dialog(npc, cfg.letter_reply)
    return true
end

local function handle_hub(me, npc, sm, map, cfg)
    if cfg.letter_before_leader and read_letter(me, npc, sm, cfg) then
        return
    end
    if not pq.is_leader(me) then
        me:dialog(npc,
            "알카드노와 제뉴미스트 사이의 분쟁을 멈추기 위해서는 알카드노와 제뉴미스트의 실험 자료가 필요해요. 알카드노의 실험 자료를 먼저 제게 가져다 주세요. 또한, 파티장이 제게 말을 걸어주셔야 합니다.")
        return
    end
    if not cfg.letter_before_leader and read_letter(me, npc, sm, cfg) then
        return
    end
    if pq.has_item(me, 4001134, 1) then
        pq.remove_all(4001134, me)
        sm:set_property("stage4", "1")
        me:dialog(npc, "알카드노 실험 자료를 찾아오셨군요! 이제 제뉴미스트 실험 자료를 찾으면 되겠군요.")
        return
    end
    if pq.has_item(me, 4001135, 1) and sm:get_property("stage4") == "1" then
        pq.remove_all(4001135, me)
        clear_fx(map)
        map:set_respawn(false)
        map:kill_all_mobs()
        pq.party_exp(sm, cfg.stage4_exp)
        sm:set_property("stage4", "2")
        local door = map:find_reactor_name(cfg.door)
        if door ~= nil then
            door:hit(1)
        end
        me:dialog(npc, "제뉴미스트 실험 자료도 찾아오셨군요. 다음 스테이지로 진행해 주세요.")
        return
    end
    me:dialog(npc, "알카드노와 제뉴미스트 사이의 분쟁을 멈추기 위해서는 알카드노와 제뉴미스트의 실험 자료가 필요해요. 알카드노의 실험 자료를 먼저 제게 가져다 주세요.")
end

local function handle_reward(me, npc, cfg)
    me:dialog(npc,
        "이렇게 유레테도 마음을 고쳐먹었고, 제뉴미스트와 알카드노의 전쟁도 조금은 수그러 들 수 있을까요? 비록 저희 사랑이 이뤄지기까지 아직도 험난한 고난이 기다리고 있지만, 결코 포기하지 않고 " ..
        cfg.beloved .. " 지키겠습니다.")
    local code = me:exchange({}, { item = { [cfg.marble] = 1 } })
    if code ~= ExchangeResult.OK then
        me:dialog(npc, "인벤토리 공간이 부족하신 건 아닌지 확인해 주시겠어요?")
        return
    end
    strip_items(me)
    me:end_party_quest(RANKING_QUEST)
    me:records():add("romeo_juliet_pq.clears", 1)
    me:map(cfg.base + 700)
end

-- cfg.base is the first map of this NPC's own run, cfg.other the first map of the other run.
function M.npc(cfg)
    return {
        on_click = function(me, npc)
            local map = me:map()
            if map == nil or map:wz() == nil then
                return
            end
            local map_id = map:wz():id()
            if map_id == cfg.entry_map then
                local sel = me:dialog_list(npc,
                    "#e<파티퀘스트 : 로미오와 줄리엣>#n\r\n마가티아는 지금 크나큰 위기를 맞이하고 있습니다. 용감한 메이플의 모험가들께서 저희를 도와주시지 않으시겠어요?\r\n\r\n#b",
                    {
                        cfg.story,
                        "함께 할 파티를 찾고 싶어요.",
                        "퀘스트를 시작한다.",
                    })
                if sel == 1 then
                    me:dialog(npc,
                        "#e<파티퀘스트 : 로미오와 줄리엣>#n\r\n로미오와 줄리엣을 도와주세요.\r\n - #e레벨#n : 71이상 #r(추천레벨 : 71 ~ 85)#k\r\n - #e참가인원#n : 4명\r\n - #e최종 보상#n : #v1122010# #z1122010##k #b(알카드노, 제뉴미스트의 구슬 각각 25개로 교환)#k")
                elseif sel == 2 then
                    me:dialog(npc,
                        "현재 미구현된 컨텐츠입니다.\r\n\r\n컨텐츠 오픈일 : #b하루속히 빨리 오픈하겠습니다.#k\r\n\r\n파티를 구하고 싶으시면 #b파티원 찾기#k 버튼을 누르시면 됩니다.")
                elseif sel == 3 then
                    try_start(me, npc, cfg)
                end
                return
            end
            local sm = me:state_machine()
            if map_id == cfg.base then
                me:dialog(npc, cfg.lab_hint)
            elseif map_id == cfg.base + 1 then
                me:dialog(npc, cfg.mobs_hint)
            elseif map_id == cfg.base + 100 then
                me:dialog(npc, cfg.beaker_hint)
            elseif map_id == cfg.base + 200 then
                if sm == nil then
                    return
                end
                handle_hub(me, npc, sm, map, cfg)
            elseif map_id == cfg.base + 300 then
                me:dialog(npc, "각 파티원이 각각 연구실을 통해 꼭대기에 있는 중앙 연구실 까지 올라가야 해요.")
            elseif map_id == cfg.base + 400 then
                me:dialog(npc, "준비가 되면, 제 사랑을 구하러 가요.")
            elseif map_id == cfg.base + 401 then
                if sm ~= nil then
                    pq.party_warp(sm, cfg.base + 500)
                end
            elseif map_id == cfg.other + 401 then
                me:dialog(npc, "구해주셔서 진심으로 감사드립니다.")
            elseif map_id == cfg.base + 600 then
                handle_reward(me, npc, cfg)
            elseif map_id == cfg.other + 600 then
                me:dialog(npc, "저를 구해주신 모험가님. 진심으로 감사드립니다.")
            end
        end
    }
end

return M
