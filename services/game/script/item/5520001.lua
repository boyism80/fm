-- Item name (String.wz/Cash.img.xml): 플래티넘 카르마의 가위

return {
    on_cash = function(me, item_id, text, ear, pet_sn, input)
        return me:use_karma(input.target_type, input.target_slot, true)
    end,
}
