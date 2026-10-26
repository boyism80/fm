-- Item name (String.wz/Cash.img.xml): 보따리상인 묘묘

local MERCHANT = 9090000
local MIN_LEVEL = 10
local STOCK = {
    { item = 2070011, price = 400000 },
    { item = 2000003, price = 100 },
    { item = 2000004, price = 50000 },
    { item = 2002000, price = 500 },
    { item = 2002001, price = 500 },
    { item = 2002002, price = 500 },
    { item = 2002003, price = 550 },
    { item = 2002004, price = 500 },
    { item = 2002005, price = 500 },
    { item = 2001000, price = 1600 },
    { item = 2001001, price = 1150 },
    { item = 2001002, price = 2000 },
    { item = 2010000, price = 50 },
    { item = 2010001, price = 100 },
    { item = 2010002, price = 50 },
    { item = 2010003, price = 100 },
    { item = 2010004, price = 300 },
    { item = 2020012, price = 2250 },
    { item = 2020013, price = 5600 },
    { item = 2020014, price = 8100 },
    { item = 2020015, price = 10200 },
    { item = 2022003, price = 1100 },
    { item = 2022000, price = 1650 },
    { item = 2060000, price = 1 },
    { item = 2061000, price = 1 },
    { item = 2060001, price = 5 },
    { item = 2061001, price = 5 },
    { item = 2030000, price = 250 },
    { item = 2330000, price = 500 },
    { item = 2070000, price = 500 },
}

return {
    on_cash = function(me, item_id, text, ear, pet_sn, input)
        local map = me:map()
        if map:wz():limits(FieldLimit.Migrate) then
            me:message("이곳에서는 사용할 수 없습니다.", Msg.Popup)
            return false
        end
        if me:level() < MIN_LEVEL then
            me:message("레벨 10 이상만 사용 가능합니다.", Msg.PinkText)
            return false
        end
        if me:state_machine() ~= nil then
            me:message("이곳에서는 사용할 수 없습니다.", Msg.PinkText)
            return false
        end
        return me:open_shop(MERCHANT, STOCK)
    end,
}
