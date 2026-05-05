-- Item name (String.wz/Consume.img.xml): 투구 지력 주문서 100%

function on_scroll(me, scroll_slot, target_slot)
    return me:enhance(scroll_slot, target_slot, EquipmentPart.Cap, nil, function(target, scroll)
        target:add_bonus_stats(scroll:wz():bonus_stats())
    end)
end
