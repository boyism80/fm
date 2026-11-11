local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/dragonNest.lua"
local NEST = 240040611
local EXIT = 240040610
local EGG = 4001094
local GUARDIAN = 2081008

local BUSY = "이미 이 안에 들어가 누군가가 퀘스트를 완수하는 중입니다. 잠시 후 다시 시도해 주세요."

local function enter(me)
    run_on_map(NEST, SCRIPT, "remove_guardian")
    me:play_portal_sound()
    local function on_arrive(me)
        me:clock(300, function(me)
            me:map(EXIT)
        end)
    end
    me:map(NEST, 0, { callback = on_arrive })
end

return {

    remove_guardian = function(map)
        for _, n in pairs(map:npcs()) do
            if n:id() == GUARDIAN then
                map:remove_npc(n)
            end
        end
    end,

    on_enter = function(me)
        if pq.has_item(me, EGG) == false then
            me:message("이 안에 들어가기 위해서는 <나인 스피릿의 알> 이 필요하다.", Msg.PinkText)
            return
        end
        local refused = pq.enter_rooms(me:map(), SCRIPT, { { NEST } }, BUSY, function()
            enter(me)
        end)
        if refused ~= nil then
            me:message(refused, Msg.PinkText)
        end
    end
}
