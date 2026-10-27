local pq = require("script/integration/lib/party_quest")

local LOBBY = 200080101
local ENTRY_NPC = 2013000
local GUIDE = 2013001
local MINERVA = 2013002
local START = 920010000
local HUB = 920010100
local WALKWAY = 920010200
local STORAGE = 920010300
local LOUNGE = 920010400
local SEALED = 920010500
local LODGING = 920010600
local UP_PATH = 920010700
local GARDEN = 920010800
local BONUS = 920011100
local REWARD = 920011300
local EXIT = 920011200
local DIALOG_NEXT = 0
local DIALOG_YES_NO = 1
local DIALOG_LIST = 4
local MENU_ENTER = 0
local MENU_FEED = 2
local RED_BEAN_SNOW = 2001002
local RED_BEAN_SNOW_INDEX = 6
local FIGHTER = 110
local FIGHTER_66_70 = 4

local CLOUD = 4001063
local WALK_PIECE = 4001050
local LODGING_PIECE = 4001052
local SEED = 4001053
local PIXIE_SEED = 4001054
local LIFE_GRASS = 4001055
local DISC_BASE = 4001055
local STATUE_PIECES = { 4001044, 4001045, 4001046, 4001047, 4001048, 4001049 }

local EAK_REACTOR = 2006000
local MINERVA_REACTOR = 2006001
local MUSIC_REACTOR = 2008006
local DISC_REACTOR_BASE = 2002003
local LOUNGE_BOX = 2002011
local SEALED_BOX = 2002012
local UP_PATH_BOX = 2002013
local NEPENTHES_POT = 2001001
local GRASS_POT = 2002003
local STATUE_REACTORS = {
    [4001044] = 2008003,
    [4001045] = 2008005,
    [4001046] = 2008000,
    [4001047] = 2008002,
    [4001048] = 2008004,
    [4001049] = 2008001,
}

local CELION = 9300040
local GARDEN_MOB = 9300048
local DARK_NEPENTHES = 9300049
local PAPA_PIXIE = 9300039
local CELION_X = { 200, -300, -300, -300, 200, 200, 200, -300, -300, 200, 200, -300, -300, 200 }
local CELION_Y = { -2321, -2114, -2910, -2510, -1526, -2716, -717, -1310, -3357, -1912, -1122, -1736, -915, -3116 }

local function leave_room(ctx)
    if pq.portal(ctx, ctx:bot(0), "st00", HUB) == false then
        return false
    end
    return pq.wait_all(ctx, HUB, 10000)
end

local function enter_room(ctx, portal, map_id)
    local leader = ctx:bot(0)
    if pq.portal(ctx, leader, portal, map_id) == false then
        return false
    end
    return leader
end

local function wonky_menu(ctx, bot, selected, reply)
    local oid = pq.npc(ctx, bot, ENTRY_NPC)
    if oid == false then
        return false
    end
    local dlg = bot:npc_click(oid)
    if dlg == false or dlg.selections == nil then
        return ctx:fail("웡키 선택지가 오지 않음")
    end
    return bot:request(reply, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = selected }, nil,
        5000)
end

local function feed_wonky(ctx, bot)
    if wonky_menu(ctx, bot, MENU_FEED, resp.dialog_list) == false then
        return ctx:fail("웡키 음식 목록이 오지 않음")
    end
    return bot:request(resp.dialog, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = RED_BEAN_SNOW_INDEX },
        nil, 5000)
end

local function stage4_feedback(text)
    if text:find("축하해요", 1, true) ~= nil then
        return 3
    end
    local ok = text:match("파티원들이 (%d+) 개의 발판")
    if ok ~= nil then
        return tonumber(ok)
    end
    return nil
end

test_suite {
    name = "Party Quest: 오르비스 PQ",
    bot_count = 6,

    on_initialize = function(ctx)
        for i = 0, ctx:bot_count() - 1 do
            local bot = ctx:bot(i)
            if pq.command(bot, "/레벨바꾸기 70", "레벨 설정") == false then
                return ctx:fail(bot:name() .. " 레벨 설정 실패")
            end
            if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
                return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
            end
            if bot:instance_move(LOBBY) == false then
                return false
            end
        end
        return true
    end,

    on_finished = function(ctx)
        pq.command(ctx:bot(0), "/SM상한 orbis_party_quest 기본", "SM상한 orbis_party_quest: 기본")
        ctx:bot(0):request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
    end,

    scenarios = {
        function(ctx)
            local leader = ctx:bot(0)
            local init = string.format("/봇초기화 70 %d 0 %d:1 - - wonky.mission=%d", FIGHTER, RED_BEAN_SNOW, FIGHTER_66_70)
            if pq.command(leader, init, "봇초기화 완료") == false then
                return ctx:fail("파티장 초기화 실패")
            end
            local taste = feed_wonky(ctx, leader)
            if taste == false or taste.text:find("정말 고마워", 1, true) == nil then
                return ctx:fail("처음 먹는 음식 반응이 오지 않음")
            end
            local mood = leader:request(resp.notice, req.dialog { dialog_type = DIALOG_NEXT, next = true }, function(p)
                return p.message:find("팥빙수을(를) 먹고 웡키의 기분이", 1, true) ~= nil
            end, 5000)
            if mood == false then
                return ctx:fail("웡키의 기분 변화 안내가 없음")
            end
            if (leader:items()[RED_BEAN_SNOW] or 0) ~= 0 then
                return ctx:fail("먹인 팥빙수가 남아 있음")
            end

            local scold = feed_wonky(ctx, leader)
            if scold == false or scold.text:find("장난치지 말라구", 1, true) == nil then
                return ctx:fail("없는 음식을 줬는데 꾸중이 없음")
            end
            leader:dialog(false)
            return true
        end,
        pq.form_party,
        -- Another party in the tower: Wonky refuses before blessing, so the mission stays for the next entry.
        function(ctx)
            local leader = ctx:bot(0)
            if pq.command(leader, "/SM상한 orbis_party_quest 0", "SM상한 orbis_party_quest: 0") == false then
                return ctx:fail("SM상한 변경 실패")
            end
            local blessed = false
            ctx:hook("notice", function(_, _, p)
                if p.message ~= nil and p.message:find("요정 웡키의 미션을 수행하여", 1, true) ~= nil then
                    blessed = true
                end
            end)
            local mission = wonky_menu(ctx, leader, MENU_ENTER, resp.dialog_yes_no)
            if mission == false then
                ctx:unhook("notice")
                return ctx:fail("웡키의 축복 미션 안내가 없음")
            end
            local matched = leader:request(resp.dialog, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, nil,
                5000)
            if matched == false or matched.text:find("딱 맞췄어", 1, true) == nil then
                ctx:unhook("notice")
                return ctx:fail("미션을 맞췄는데 축복 안내가 없음")
            end
            local refused = leader:request(resp.dialog, req.dialog { dialog_type = DIALOG_NEXT, next = true }, nil, 5000)
            ctx:sleep(1000)
            ctx:unhook("notice")
            leader:dialog(false)
            if refused == false or refused.text:find("이미 다른 파티가", 1, true) == nil then
                return ctx:fail("다른 파티 진행 중 거절 대사가 없음")
            end
            if blessed then
                return ctx:fail("입장이 거절됐는데 웡키의 축복이 걸림")
            end
            if pq.command(leader, "/SM상한 orbis_party_quest 기본", "SM상한 orbis_party_quest: 기본") == false then
                return ctx:fail("SM상한 복원 실패")
            end
            return true
        end,
        function(ctx)
            local leader = ctx:bot(0)
            local mission = wonky_menu(ctx, leader, MENU_ENTER, resp.dialog_yes_no)
            if mission == false or mission.text:find("파이터", 1, true) == nil or mission.text:find("66~70", 1, true) == nil then
                return ctx:fail("웡키의 축복 미션 안내가 없음")
            end
            local blessed = leader:request(resp.dialog, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, nil,
                5000)
            if blessed == false or blessed.text:find("딱 맞췄어", 1, true) == nil then
                return ctx:fail("미션을 맞췄는데 축복 안내가 없음")
            end
            local notice = leader:request(resp.notice, req.dialog { dialog_type = DIALOG_NEXT, next = true }, function(p)
                return p.message:find("요정 웡키의 미션을 수행하여", 1, true) ~= nil
            end, 5000)
            if notice == false then
                return ctx:fail("웡키의 축복이 적용되지 않음")
            end
            if pq.wait_all(ctx, START, 15000) == false then
                return ctx:fail("여신의 탑 입장 실패")
            end
            for i = 0, ctx:bot_count() - 1 do
                ctx:bot(i):dialog(false)
            end
            return true
        end,
        function(ctx)
            local leader = ctx:bot(0)
            if pq.sweep(ctx, leader, CLOUD, 20, { [EAK_REACTOR] = true }) == false then
                return false
            end
            if (leader:items()[CLOUD] or 0) < 20 then
                return ctx:fail("구름 조각 부족: " .. (leader:items()[CLOUD] or 0))
            end
            if pq.feed(ctx, leader, EAK_REACTOR, CLOUD, 20) == false then
                return false
            end
            local oid = leader:npc(GUIDE)
            if oid == nil then
                return ctx:fail("시종 이크가 나타나지 않음")
            end
            local dlg = leader:npc_click(oid)
            if dlg == false or dlg.text:find("입구로 데려다", 1, true) == nil then
                return ctx:fail("이크 복구 대화가 오지 않음")
            end
            leader:dialog(false)
            if leader:request(resp.warp, nil, nil, 10000) == false then
                return ctx:fail("이크가 입구로 보내지 않음")
            end
            for i = 0, ctx:bot_count() - 1 do
                if pq.portal(ctx, ctx:bot(i), "in00", HUB) == false then
                    return false
                end
            end
            return true
        end,
        function(ctx)
            local leader = enter_room(ctx, "in00", WALKWAY)
            if leader == false then
                return false
            end
            if pq.talk(ctx, leader, GUIDE, "산책로") == false then
                return false
            end
            if pq.collect(ctx, leader, nil, WALK_PIECE, 30) == false then
                return false
            end
            if pq.talk(ctx, leader, GUIDE, "훌륭해요") == false then
                return false
            end
            return leave_room(ctx)
        end,
        function(ctx)
            local leader = enter_room(ctx, "in01", STORAGE)
            if leader == false then
                return false
            end
            if pq.talk(ctx, leader, GUIDE, "창고") == false then
                return false
            end
            local spot = nil
            for _, s in ipairs(leader:mob_spots()) do
                if s.id == CELION then
                    spot = s
                end
            end
            if spot == nil then
                return ctx:fail("샐리온 배치가 없음")
            end
            if pq.kill_at(ctx, leader, CELION, spot.x, spot.y) == false then
                return false
            end
            for i = 1, #CELION_X do
                if pq.kill_at(ctx, leader, CELION, CELION_X[i], CELION_Y[i]) == false then
                    return false
                end
            end
            if pq.pick_up(ctx, leader, STATUE_PIECES[2]) == false then
                return false
            end
            if pq.talk(ctx, leader, GUIDE, "훌륭해요") == false then
                return false
            end
            return leave_room(ctx)
        end,
        function(ctx)
            local leader = enter_room(ctx, "in02", LOUNGE)
            if leader == false then
                return false
            end
            local today = os.date("*t").wday
            local disc = DISC_BASE + today
            if pq.open_box(ctx, leader, DISC_REACTOR_BASE + today, disc) == false then
                return false
            end
            if pq.feed(ctx, leader, MUSIC_REACTOR, disc, 1) == false then
                return false
            end
            if pq.talk(ctx, leader, GUIDE, "바로 이 음악") == false then
                return false
            end
            if pq.open_box(ctx, leader, LOUNGE_BOX, STATUE_PIECES[3]) == false then
                return false
            end
            return leave_room(ctx)
        end,
        function(ctx)
            local leader = enter_room(ctx, "in03", SEALED)
            if leader == false then
                return false
            end
            for i = 1, ctx:bot_count() - 1 do
                if pq.portal(ctx, ctx:bot(i), "in03", SEALED) == false then
                    return false
                end
            end
            local areas = leader:areas()
            if #areas < 3 then
                return ctx:fail("봉인된 방 발판 정보가 없음")
            end
            local candidates = {}
            for a = 0, 5 do
                for b = 0, 5 - a do
                    candidates[#candidates + 1] = { a, b, 5 - a - b }
                end
            end
            local oid = pq.npc(ctx, leader, GUIDE)
            if oid == false then
                return false
            end
            for _ = 1, 7 do
                local guess = candidates[1]
                if guess == nil then
                    return ctx:fail("봉인된 방 정답 후보가 없음")
                end
                local bot = 1
                for area, count in ipairs(guess) do
                    for _ = 1, count do
                        if pq.move(ctx:bot(bot), areas[area].x, areas[area].y) == false then
                            return ctx:fail(ctx:bot(bot):name() .. " 발판 이동 실패: " .. area)
                        end
                        bot = bot + 1
                    end
                end
                local dlg = leader:npc_click(oid)
                if dlg == false then
                    return ctx:fail("봉인된 방 판정 대화가 오지 않음")
                end
                leader:dialog(false)
                local ok = stage4_feedback(dlg.text)
                if ok == 3 then
                    return pq.open_box(ctx, leader, SEALED_BOX, STATUE_PIECES[4]) ~= false and leave_room(ctx)
                end
                local left = {}
                for _, c in ipairs(candidates) do
                    local same = 0
                    for k = 1, 3 do
                        if c[k] == guess[k] then
                            same = same + 1
                        end
                    end
                    if c ~= guess and (ok == nil or same == ok) then
                        left[#left + 1] = c
                    end
                end
                candidates = left
            end
            return ctx:fail("봉인된 방 정답을 7번 안에 찾지 못함")
        end,
        function(ctx)
            local leader = enter_room(ctx, "in04", LODGING)
            if leader == false then
                return false
            end
            if pq.talk(ctx, leader, GUIDE, "숙박실") == false then
                return false
            end
            for i = 0, 3 do
                if (leader:items()[LODGING_PIECE] or 0) >= 40 then
                    break
                end
                if pq.portal(ctx, leader, string.format("in%02d", i), LODGING + 1 + i) == false then
                    return false
                end
                if pq.sweep(ctx, leader, LODGING_PIECE, 40) == false then
                    return false
                end
                if pq.portal(ctx, leader, "out00", LODGING) == false then
                    return false
                end
            end
            if pq.talk(ctx, leader, GUIDE, "훌륭해요") == false then
                return false
            end
            return leave_room(ctx)
        end,
        function(ctx)
            local leader = enter_room(ctx, "in05", UP_PATH)
            if leader == false then
                return false
            end
            if pq.talk(ctx, leader, GUIDE, "꼭대기로") == false then
                return false
            end
            for _, combo in ipairs(pq.combinations(5, 2)) do
                local want = { 0, 0, 0, 0, 0 }
                for _, i in ipairs(combo) do
                    want[i] = 1
                end
                for i = 1, 5 do
                    local lever = pq.seek_reactor(ctx, leader, function(r)
                        return r.name == tostring(i)
                    end)
                    if lever == false then
                        return false
                    end
                    if lever.state ~= want[i] and leader:hit_reactor(lever.oid) == false then
                        return ctx:fail("레버 응답 없음: " .. i)
                    end
                end
                local oid = pq.npc(ctx, leader, GUIDE)
                if oid == false then
                    return false
                end
                local dlg = leader:npc_click(oid)
                if dlg == false then
                    return ctx:fail("레버 판정 대화가 오지 않음")
                end
                leader:dialog(false)
                if dlg.text:find("정답 레버를 찾으셨군요", 1, true) ~= nil then
                    if pq.open_box(ctx, leader, UP_PATH_BOX, STATUE_PIECES[6]) == false then
                        return false
                    end
                    return leave_room(ctx)
                end
                if dlg.text:find("정답 레버가 아닙니다", 1, true) == nil then
                    return ctx:fail("예상과 다른 레버 판정: " .. dlg.text)
                end
            end
            return ctx:fail("정답 레버 조합을 찾지 못함")
        end,
        function(ctx)
            local leader = ctx:bot(0)
            for _, piece in ipairs(STATUE_PIECES) do
                if (leader:items()[piece] or 0) < 1 then
                    return ctx:fail("여신상 조각이 없음: " .. piece)
                end
                if pq.feed(ctx, leader, STATUE_REACTORS[piece], piece, 1) == false then
                    return false
                end
            end
            if pq.talk(ctx, leader, GUIDE, "석상이 복구되었군요") == false then
                return false
            end
            return pq.wait_all(ctx, GARDEN, 10000)
        end,
        function(ctx)
            local leader = ctx:bot(0)
            if pq.collect(ctx, leader, GARDEN_MOB, SEED, 1) == false then
                return false
            end
            local pot = pq.feed(ctx, leader, NEPENTHES_POT, SEED, 1)
            if pot == false then
                return false
            end
            if pq.kill_at(ctx, leader, DARK_NEPENTHES, pot.x, pot.y) == false then
                return false
            end
            if pq.kill_at(ctx, leader, PAPA_PIXIE, -830, 563) == false then
                return false
            end
            if pq.pick_up(ctx, leader, PIXIE_SEED) == false then
                return false
            end
            local grass_pot = pq.feed(ctx, leader, GRASS_POT, PIXIE_SEED, 1)
            if grass_pot == false then
                return false
            end
            if pq.open_box(ctx, leader, GRASS_POT, LIFE_GRASS) == false then
                return false
            end
            if pq.portal(ctx, leader, "in00", HUB) == false then
                return false
            end
            return pq.wait_all(ctx, HUB, 10000)
        end,
        function(ctx)
            local leader = ctx:bot(0)
            if pq.feed(ctx, leader, MINERVA_REACTOR, LIFE_GRASS, 1) == false then
                return false
            end
            local oid = leader:npc(MINERVA)
            if oid == nil then
                return ctx:fail("여신 미네르바가 나타나지 않음")
            end
            if leader:npc_click(oid) == false then
                return ctx:fail("미네르바 대화가 오지 않음")
            end
            leader:dialog(false)
            if pq.wait_all(ctx, BONUS, 10000) == false then
                return false
            end
            if pq.command(leader, "/타이머 1", "타이머 제한: 1") == false then
                return ctx:fail("보너스 시간 단축 실패")
            end
            if pq.command(leader, "/타이머 0", "타이머 제한: 0") == false then
                return ctx:fail("타이머 제한 해제 실패")
            end
            return pq.wait_all(ctx, REWARD, 15000)
        end,
        function(ctx)
            local leader = ctx:bot(0)
            local oid = pq.npc(ctx, leader, MINERVA)
            if oid == false then
                return false
            end
            local dlg = leader:npc_click(oid)
            if dlg == false then
                return ctx:fail("보상 대화가 오지 않음")
            end
            leader:dialog(false)
            return pq.wait_map(ctx, leader, EXIT, 10000)
        end,
    },
}
