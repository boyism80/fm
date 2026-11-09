local pq = require("script/integration/lib/party_quest")

local MAP = 100000000
local TEMPLATE = 9901000
local HALL = 102000003
local SEAT_X, SEAT_Y = 69, 196
local WARRIOR = 100

local function chat(text)
    return req.normal_chat { message = text }
end

test_suite {
    name = "Player NPC: 생성·채널 동기화·삭제",
    bot_count = 2,

    on_initialize = function(ctx)
        local gm, other = ctx:bot(0), ctx:bot(1)
        if gm:map_move(MAP) == false then
            return ctx:fail("GM 맵 이동 실패")
        end
        if other:transfer(1) == false then
            return ctx:fail("다른 봇 채널 이동 실패")
        end
        if other:map_move(MAP) == false then
            return ctx:fail("다른 봇 맵 이동 실패")
        end
        pq.command(gm, "/플레이어NPC삭제 " .. TEMPLATE, "플레이어 NPC 삭제")
        return true
    end,

    on_finished = function(ctx)
        pq.command(ctx:bot(0), "/플레이어NPC삭제 " .. TEMPLATE, "플레이어 NPC 삭제")
        ctx:bot(1):transfer(0)
    end,

    scenarios = {
        function(ctx)
            local gm, other = ctx:bot(0), ctx:bot(1)
            local look = gm:request_on(other, resp.player_npc, chat("/플레이어NPC생성 " .. TEMPLATE), function(p)
                return p.npc_id == TEMPLATE
            end, 5000)
            if look == false then
                return ctx:fail("다른 채널에 플레이어 NPC 외형이 오지 않음")
            end
            if look.look.name ~= gm:name() then
                return ctx:fail("플레이어 NPC 이름이 다름: " .. tostring(look.look.name))
            end
            if other:npc(TEMPLATE) == nil then
                return ctx:fail("다른 채널에 플레이어 NPC가 보이지 않음")
            end
            for _ = 1, 10 do
                if gm:npc(TEMPLATE) ~= nil then
                    break
                end
                ctx:sleep(200)
            end
            local oid = gm:npc(TEMPLATE)
            if oid == nil then
                return ctx:fail("만든 채널에 플레이어 NPC가 보이지 않음")
            end
            local dlg = gm:npc_click(oid)
            if dlg == false or dlg.text:find("레벨 200", 1, true) == nil then
                return ctx:fail("플레이어 NPC 클릭 대사가 WZ 기본 대사가 아님")
            end
            gm:dialog(false)
            return true
        end,
        function(ctx)
            local gm = ctx:bot(0)
            if pq.command(gm, "/플레이어NPC생성 " .. TEMPLATE, "이미 플레이어 NPC로 쓰이는") == false then
                return ctx:fail("같은 NPC로 두 번 만들어짐")
            end
            if pq.command(gm, "/플레이어NPC생성 1", "존재하지 않는 NPC") == false then
                return ctx:fail("없는 NPC로 만들어짐")
            end
            return true
        end,
        function(ctx)
            local gm, other = ctx:bot(0), ctx:bot(1)
            local oid = other:npc(TEMPLATE)
            local removed = gm:request_on(other, resp.remove_npc, chat("/플레이어NPC삭제 " .. TEMPLATE), function(p)
                return p.oid == oid
            end, 5000)
            if removed == false then
                return ctx:fail("다른 채널에서 플레이어 NPC가 사라지지 않음")
            end
            for _ = 1, 10 do
                if gm:npc(TEMPLATE) == nil then
                    return true
                end
                ctx:sleep(200)
            end
            return ctx:fail("만든 채널에서 플레이어 NPC가 사라지지 않음")
        end,
        -- A hall of fame seat stays hidden and gives a player npc of its template its place.
        function(ctx)
            local gm = ctx:bot(0)
            if gm:map_move(HALL) == false then
                return ctx:fail("전직관 이동 실패")
            end
            if gm:npc(TEMPLATE, 1000) ~= nil then
                return ctx:fail("명예의 전당 빈 자리가 NPC로 보임")
            end
            if gm:request(resp.player_npc, chat("/플레이어NPC생성 " .. TEMPLATE), function(p)
                    return p.npc_id == TEMPLATE
                end, 5000) == false then
                return ctx:fail("빈 자리에 플레이어 NPC를 세우지 못함")
            end
            local npc = gm:npcs(TEMPLATE)[1]
            pq.command(gm, "/플레이어NPC삭제 " .. TEMPLATE, "플레이어 NPC 삭제")
            if npc == nil or npc.x ~= SEAT_X or npc.y ~= SEAT_Y then
                return ctx:fail(string.format("플레이어 NPC가 WZ 자리(%d, %d)에 서지 않음: %s", SEAT_X, SEAT_Y,
                    npc and (npc.x .. ", " .. npc.y) or "없음"))
            end
            return true
        end,
        -- Reaching level 200 tells the world and stands the character's player npc on its job instructor's seat.
        function(ctx)
            local gm, other = ctx:bot(0), ctx:bot(1)
            if other:map_move(HALL) == false then
                return ctx:fail("다른 봇 전직관 이동 실패")
            end
            if pq.command(gm, "/직업바꾸기 " .. WARRIOR, "직업") == false or pq.command(gm, "/레벨바꾸기 199", "레벨 설정") == false then
                return ctx:fail("레벨 199 전사 설정 실패")
            end
            local noticed, seat = false, nil
            ctx:hook("level_200", function(_, b, p, name)
                if b:name() ~= other:name() then
                    return
                end
                if name == resp.notice and p.message:find("레벨 200을 달성", 1, true) ~= nil and p.message:find(gm:name(), 1, true) ~= nil then
                    noticed = true
                elseif name == resp.player_npc and p.look.name == gm:name() then
                    seat = p.npc_id
                end
            end)
            local told = gm:request(resp.notice, chat("/경험치얻기 4000000000"), function(p)
                return p.message:find("직업 전직관에 생겨났습니다", 1, true) ~= nil
            end, 10000)
            for _ = 1, 25 do
                if noticed and seat ~= nil then
                    break
                end
                ctx:sleep(200)
            end
            ctx:unhook("level_200")
            if seat ~= nil then
                pq.command(gm, "/플레이어NPC삭제 " .. seat, "플레이어 NPC 삭제")
            end
            if told == false then
                return ctx:fail("레벨 200에 전직관 플레이어 NPC 안내가 오지 않음 (레벨 " .. gm:level() .. ")")
            end
            if noticed == false then
                return ctx:fail("다른 채널에 레벨 200 축하 공지가 오지 않음")
            end
            if seat == nil or seat < TEMPLATE or seat >= TEMPLATE + 20 then
                return ctx:fail("전사 전직관 자리에 플레이어 NPC가 서지 않음: " .. tostring(seat))
            end
            return true
        end,
    },
}
