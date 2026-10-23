-- Item name (String.wz/Cash.img.xml): 하트 풍선

return {
    on_cash = function(me, item_id, text)
        return me:place_message_box(item_id, text)
    end,
}
