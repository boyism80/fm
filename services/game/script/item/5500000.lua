-- Item name (String.wz/Cash.img.xml): 마법의 모래시계

local DAYS = 1

return {
    on_cash = function(me, item_id, text, ear, pet_sn, input)
        return me:extend_equipment(input.equip_slot, DAYS)
    end,
}
