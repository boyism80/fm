-- Item name (String.wz/Cash.img.xml): 이름 새기기

return {
    on_cash = function(me, item_id, text, ear, pet_sn, input)
        return me:tag_equipment(input.equip_slot)
    end,
}
