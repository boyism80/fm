-- VIP존 입장하기 (Quest.wz/QuestData/2291.img): VIP존 입장하기

local pq = require("script/lib/party_quest")

local quest_id = 2291
local vip_base = 103040440
local ENTRY_KEY = "vip_entry"
local BUSY = "모든 VIP존에 누군가 있군요. 다른 채널을 이용해주세요."

return {
    on_end = function(me, npc)
        local q = me:quest(quest_id)
        if q == nil then
            return
        end

        if not me:dialog_yes_no(npc, "소지하신 VIP티켓은 모두 모아 온건가요? 어서 제게 주세요. #bVIP존#k으로 입장시켜 드릴게요. 참고로 VIP존은 #b30분에 한 번씩만 이용이 가능#k합니다.") then
            me:dialog(npc, "지금은 별로 가고 싶지 않으신가 보군요. 언제든 저를 통해서 입장할 수 있으니 말을 걸어주세요.", false, false)
            return
        end

        local rooms = {}
        for i = 0, 9 do
            rooms[#rooms + 1] = { vip_base + i, vip_base + 10 + i, vip_base + 20 + i }
        end
        local refused = pq.enter_rooms(me:map(), ENTRY_KEY, rooms, BUSY, function(index)
            if me:exchange({ item = { [4032521] = 10 } }, {}) ~= ExchangeResult.OK then
                return nil
            end
            q:force_complete(npc)
            me:show_effect(EffectType.QuestCompletion)
            me:play_portal_sound()
            me:map(rooms[index][1], 1)
        end)
        if refused ~= nil then
            me:dialog(npc, refused, false, false)
        end
    end
}
