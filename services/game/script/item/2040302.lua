-- Item name (String.wz/Consume.img.xml): 귀 장식 지력 주문서 10%

function on_scroll(me, scroll_slot, target_slot)
    return me:enhance(scroll_slot, target_slot, EquipmentPart.Eye, nil, function(target, scroll)
        target:add_bonus_stats(scroll:wz():bonus_stats())
    end)
end
