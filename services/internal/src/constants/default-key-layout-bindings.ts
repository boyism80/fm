const DEFAULT_SLOTS: Record<number, { type: number; action: number }> = {
    2: { type: 4, action: 10 }, 3: { type: 4, action: 12 }, 4: { type: 4, action: 13 }, 5: { type: 4, action: 18 },
    6: { type: 4, action: 0x17 }, 16: { type: 4, action: 8 }, 17: { type: 4, action: 5 }, 18: { type: 4, action: 0 },
    19: { type: 4, action: 4 }, 23: { type: 4, action: 1 }, 24: { type: 4, action: 0x18 }, 25: { type: 4, action: 19 },
    26: { type: 4, action: 14 }, 27: { type: 4, action: 15 }, 29: { type: 5, action: 52 }, 31: { type: 4, action: 2 },
    33: { type: 4, action: 0x19 }, 34: { type: 4, action: 17 }, 35: { type: 4, action: 11 }, 37: { type: 4, action: 3 },
    38: { type: 4, action: 20 }, 40: { type: 4, action: 16 }, 41: { type: 4, action: 0x16 }, 42: { type: 4, action: 0x15 },
    43: { type: 4, action: 9 }, 44: { type: 5, action: 50 }, 45: { type: 5, action: 51 }, 46: { type: 4, action: 6 },
    50: { type: 4, action: 7 }, 56: { type: 5, action: 53 }, 57: { type: 5, action: 0x36 }, 59: { type: 6, action: 100 },
    60: { type: 6, action: 101 }, 61: { type: 6, action: 102 }, 62: { type: 6, action: 103 }, 63: { type: 6, action: 104 },
    64: { type: 6, action: 105 }, 65: { type: 6, action: 106 },
};

export function getDefaultKeyLayoutBindings() {
    return Object.entries(DEFAULT_SLOTS)
        .map(([k, v]) => ({ slot: Number(k), type: v.type >>> 0, action: v.action | 0 }))
        .sort((a, b) => a.slot - b.slot);
}
