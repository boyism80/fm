local pq = require("script/integration/lib/party_quest")

local TRAINING_ROOM = 120000104
local ROAD_1 = 925010000
local ROAD_2 = 925010100
local ROAD_3 = 925010200
local PROTECT = 925010300
local SHELTER = 925010400
local SCHRINTZ = 1092008
local DELLI_NPC = 2095000
local DELLI_MOB = 9300162

local DIALOG_DEFAULT = 0
local DIALOG_YES_NO = 1
local PROTECT_SECONDS = 10

test_suite {
    name = "Quest: 델리를 지켜라·파티 탈퇴 뒤에도 남음",
    bot_count = 2,

    on_finished = function(ctx)
        pq.command(ctx:bot(0), "/타이머 0", "타이머 제한: 0")
    end,

    on_initialize = function(ctx)
        local bot = ctx:bot(0)
        if pq.command(bot, "/봇초기화 70 500 0 - 6410:1", "봇초기화 완료") == false then
            return ctx:fail("봇 초기화 실패")
        end
        if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
            return ctx:fail("플레이어 모드 설정 실패")
        end
        if pq.command(ctx:bot(1), "/봇초기화 70 0 0 -", "봇초기화 완료") == false then
            return ctx:fail("파티장 봇 초기화 실패")
        end
        return bot:instance_move(TRAINING_ROOM)
    end,

    scenarios = {
        -- The bot enters as a party member, and leaving the party leaves it on the road: the original kept it.
        function(ctx)
            local bot, leader = ctx:bot(0), ctx:bot(1)
            if leader:request(resp.party_created, req.party_operation { operation = PARTY.Create }) == false then
                return ctx:fail("파티 생성 실패")
            end
            local invite = leader:request_on(bot, resp.party_invite,
                req.party_operation { operation = PARTY.Invite, target_name = bot:name() })
            if invite == false or bot:request(resp.party_update_join,
                    req.party_operation { operation = PARTY.AcceptInvite, party_id = invite.party_id }) == false then
                return ctx:fail("파티 가입 실패")
            end

            local oid = pq.npc(ctx, bot, SCHRINTZ)
            if oid == false then
                return false
            end
            if bot:npc_click(oid) == false then
                return ctx:fail("슈린츠 대화가 오지 않음")
            end
            if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
                    return p.character.map == ROAD_1
                end, 10000) == false then
                return ctx:fail("델리를 찾으러 가는 길 입장 실패")
            end

            if bot:request(resp.party_update_leave, req.party_operation { operation = PARTY.Leave }, nil, 5000) == false then
                return ctx:fail("파티 탈퇴 실패")
            end
            leader:request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
            ctx:sleep(1500)
            if bot:map() ~= ROAD_1 then
                return ctx:fail("파티를 나갔더니 델리를 찾으러 가는 길에서 쫓겨남: " .. bot:map())
            end
            return true
        end,
        function(ctx)
            local bot = ctx:bot(0)
            if pq.portal(ctx, bot, "out00", ROAD_2) == false then
                return false
            end
            return pq.portal(ctx, bot, "out01", ROAD_3)
        end,
        function(ctx)
            local bot = ctx:bot(0)
            local oid = pq.npc(ctx, bot, DELLI_NPC)
            if oid == false then
                return false
            end
            if bot:npc_click(oid) == false then
                return ctx:fail("델리 대화가 오지 않음")
            end
            if bot:dialog(true) == false then
                return ctx:fail("델리 부탁 대화가 오지 않음")
            end
            if bot:dialog(true) == false then
                return ctx:fail("델리 보호 안내가 오지 않음")
            end
            if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
                    return p.character.map == PROTECT
                end, 10000) == false then
                return ctx:fail("델리를 지켜라 맵으로 이동하지 않음")
            end
            if bot:mobs(DELLI_MOB)[1] == nil and bot:request(resp.spawn_mob, nil, function(p)
                    return p.mob.mob_id == DELLI_MOB
                end, 5000) == false then
                return ctx:fail("보호할 델리가 없음")
            end
            if pq.command(bot, "/타이머 " .. PROTECT_SECONDS, "타이머 제한: " .. PROTECT_SECONDS) == false then
                return ctx:fail("타이머 제한 설정 실패")
            end
            return true
        end,
        function(ctx)
            local bot = ctx:bot(0)
            if pq.wait_map(ctx, bot, SHELTER, PROTECT_SECONDS * 1000 + 15000) == false then
                return ctx:fail("보호 시간이 끝난 뒤 은신처로 이동하지 않음")
            end
            if bot:mobs(DELLI_MOB)[1] ~= nil then
                return ctx:fail("은신처에 몬스터 델리가 남아 있음")
            end
            local oid = pq.npc(ctx, bot, DELLI_NPC)
            if oid == false then
                return false
            end
            if bot:npc_click(oid) == false then
                return ctx:fail("은신처 델리 감사 대화가 오지 않음")
            end
            return true
        end,
        function(ctx)
            local bot = ctx:bot(0)
            if bot:instance_move(TRAINING_ROOM) == false then
                return ctx:fail("훈련장으로 돌아가기 실패")
            end
            return true
        end,
    },
}
