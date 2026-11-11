local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/2ndTowerTop.lua"
local TOWER = 211060401
local TEMP_KEY = 4032833

local REFUSED = "두번째 탑의 임시열쇠가 없거나 이미 누군가가 안에서 퀘스트를 진행중입니다."

local function enter(me)
    if pq.has_item(me, TEMP_KEY) then
        me:rmitem(TEMP_KEY, 1)
    end
    me:map(TOWER, 0)
end

return {

    on_enter = function(me)
        local allowed = me:quest(3140):started()
            and (pq.has_item(me, TEMP_KEY) or me:quest(3166):completed())
        if allowed == false then
            me:message(REFUSED, Msg.PinkText)
            return
        end
        local refused = pq.enter_rooms(me:map(), SCRIPT, { { TOWER } }, REFUSED, function()
            enter(me)
        end)
        if refused ~= nil then
            me:message(refused, Msg.PinkText)
        end
    end
}
