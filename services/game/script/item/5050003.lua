-- Item name (String.wz/Cash.img.xml): 3차 스킬 SP 되돌리기 주문서

-- The scroll moves SP only between skills of its own job advancement.
return {
    on_cash = function(me, item_id, text, ear, pet_sn, input)
        return me:reset_sp(input.to, input.from, item_id - 5050000)
    end,
}
