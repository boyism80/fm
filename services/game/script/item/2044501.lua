-- Item name (String.wz/Consume.img.xml): 활 공격력 주문서 60%

function on_scroll(me, scroll_slot, target_slot)
    return me:enhance(scroll_slot, target_slot, EquipmentPart.Weapon, WeaponType.Bow, function(target, scroll)
        target:add_bonus_stats(scroll:wz():bonus_stats())
    end)
end
