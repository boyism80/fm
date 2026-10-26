-- Item name (String.wz/Cash.img.xml): 봉인의 자물쇠

return {
    on_cash = function(me, item_id, text, ear, pet_sn, input)
        return me:lock_item(input.target_type, input.target_slot)
    end,
}
