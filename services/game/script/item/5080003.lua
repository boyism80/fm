-- Item name (String.wz/Cash.img.xml): 입학축하 플래카드

return {
    on_cash = function(me, item_id, text)
        return me:place_message_box(item_id, text)
    end,
}
