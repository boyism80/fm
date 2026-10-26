-- Item name (String.wz/Cash.img.xml): AP 되돌리기 주문서

return {
    on_cash = function(me, item_id, text, ear, pet_sn, input)
        return me:reset_ap(input.to, input.from)
    end,
}
