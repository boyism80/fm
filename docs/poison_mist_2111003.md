# 포이즌 미스트 (2111003) — 자바 동작 정리 및 Go 구현 시 추가 과제

## 1. 스킬 성격 (요약)

`String.wz/Skill.img.xml` 설명(주변에 독 안개, 안개 안 몬스터 중독, 대상 수 제한 없음)과 같이, **직접 타격 스킬이 아니라 맵에 설치되는 필드(미스트) 오브젝트**로 처리된다.

| 구분 | 동작 |
|------|------|
| 클라이언트 분류 | `Skill.java`: `isBuff = false` (2111003) — 일반 자기 버프가 아님 |
| 서버 분류 | `MapleStatEffect.isMist()`: 2111003 포함 — 도어/버프와 별도 **미스트** 경로 |
| 설치 | `applyMist` → `Rectangle` 바운드 + `MapleMist` 생성 → `map.spawnMist(mist, duration, false)` |
| 지속 | 스킬 이펙트 `duration` 만큼 맵에 유지 후 제거 |
| 효과 | 미스트 **사각형 안 몬스터**에 대해 **주기적으로** 중독 적용 시도 |

자바 참조 파일: `old/src/client/Skill.java`, `old/src/server/MapleStatEffect.java`, `old/src/server/maps/MapleMist.java`, `old/src/server/maps/MapleMap.java`.

---

## 2. 자바 상세 동작

### 2.1 스킬 적용 분기 (`MapleStatEffect.applyTo`)

- MP/기타 처리 후 `isMist()` 이면 **`applyMist(caster, pos)`** 만 수행.
- `applyMist`: 시전자 위치(또는 `pos`)와 좌우를 반영한 **`calculateBoundingBox`** 로 영역을 만들고,  
  `new MapleMist(bounds, caster, this)` (`this` = `MapleStatEffect`) 후 **`spawnMist(mist, getDuration(), false)`**.

### 2.2 미스트 객체 (`MapleMist`)

- 플레이어 미스트: `MapleMist(bounds, owner, MapleStatEffect source)`.
- 2111003 / 14111006 등은 생성자에서 **`isPoisonMist = 0`** (맵 타이머 분기용 타입).
- `makeChanceResult()` → WZ **prop** 기반 확률.

### 2.3 맵 스폰 및 주기 (`MapleMap.spawnMist`)

`mist.isPoisonMist() == 0` 인 경우(포이즌 미스트 등):

- 맵에 미스트 오브젝트 등록 + 스폰 패킷 브로드캐스트.
- **`MapTimer.register(runnable, 2000, 2500)`**: 초기 지연 **2000ms**, 이후 **2500ms** 주기 반복.
- 매 틱: `mist.getBox()` 와 교차하는 **몬스터**만 순회.
  - `mist.makeChanceResult()` (prop) **그리고** 해당 몬스터가 이미 `POISON` 이 아니면  
    `MapleMonster.applyStatus(..., POISON, ..., mist.getSource())` 로 중독 부여.
- 별도 스케줄: `duration` 경과 시 미스트 제거 패킷 + 맵에서 제거 + 위 주기 작업 **취소**.

즉, **매직 컴포지션(2111006)처럼 `on_attack` 한 번에 걸리는 타입이 아니라**, 미스트가 깔린 뒤 **틱마다 영역 내 몬스터에게 독을 시도**하는 구조다.

### 2.4 `isPoisonMist == 1` (참고)

같은 `spawnMist` 스위치에서 `case 1` 은 **플레이어** 대상 독 디버프 등 다른 처리. 2111003 구현 시에는 **case 0** 경로가 대응 대상.

---

## 3. 현재 Go 코드베이스와의 관계

이미 있는 것:

- 몬스터 `MobStatus.POISON`, `ApplyMobStatus` / `GetMobStatusValue`, 맵 단위 **독 도트 타이머**(`MobPoisonTickTimer`)로 틱 데미지.
- Lua: `compute_poison_tick_damage`, `element_amp_from_class`, `element_weak_multiplier` 등 **pDam 계산** (직접 스킬 적중 시나리오에 맞춤).

아직 없는 것:

- **맵에 존재하는 미스트 엔티티**(영역, 소유자, 소스 스킬 이펙트, 만료).
- 스킬 사용 시 **미스트 스폰**으로 분기하는 **핸들러/액터 흐름**.
- 미스트 전용 **주기 작업**(자바: 2s 지연 + 2.5s 주기; Go에서는 `MapActor` 스케줄러 또는 전용 메시지로 정렬 필요).
- 미스트 영역과 몹 교차 판정 + **prop** 확률 + 이미 중독이면 스킵.
- (클라이언트) **미스트 스폰/제거 패킷** 및 기존 `tools/packet` / `MobPacket` 과의 정합.

---

## 4. Go에서 구현하기 위해 추가로 필요한 작업 (체크리스트)

### 4.1 도메인 / 엔티티

- [ ] **`Mist` (또는 `MapMist`) 타입**: 맵 ID, OID, 소유 캐릭터 ID, `Rectangle`/LT·RB, 지속 시간, 소스 스킬 ID·레벨(또는 `SkillEntry` 참조), 미스트 종류(`PoisonMist` 등).
- [ ] **`Map`에 미스트 컬렉션**: 스폰/제거/만료, (선택) 쿼리용 공간 인덱스는 초기에는 단순 순회로도 가능.

### 4.2 스킬 사용 → 미스트 설치

- [ ] 스킬 2111003 (및 `isMist()` 에 해당하는 다른 ID) 사용 시 **일반 마법 데미지 패킷이 아닌** `spawn mist` 경로.
- [ ] WZ에서 **duration**, **prop**, **lt/rb** (또는 `calculateBoundingBox` 에 쓰는 동일 데이터) 읽기 — 이미 스킬 `LevelData` 일부는 로드됨; 미스트 전용 필드가 부족하면 **로더·`wz.Skill` 보강**.
- [ ] 시전 위치·방향에 따른 **바운드 계산** (자바 `calculateBoundingBox` 와 동일 규칙 정의).

### 4.3 주기 로직 (자바 `spawnMist` case 0)

- [ ] **첫 틱 지연 2초**, 이후 **2.5초마다** 반복 (또는 WZ/밸런스에 맞게 상수화). `MapActor` + `TimerScheduler` / 전용 타이머 메시지로 **맵 액터 단일 스레드**에서 실행 권장.
- [ ] 각 틱: 미스트 박스와 교차하는 **몹** 나열 → `prop` 판정 → 이미 독이 없을 때만 `ApplyMobStatus` + Lua에서 **pDam** 을 채우려면 **시전자·스킬**을 넘기는 경로 통일 (`apply_prob_status` 와 유사하거나 공용 헬퍼).

### 4.4 독 데미지 수치와의 연결

- [ ] 미스트로 걸린 독도 **자바처럼** `MapleStatEffect` 기반 pDam 이면, Go에서는 **같은 Lua `compute_poison_tick_damage`** 를 쓰되 인자로 **시전자 `me`**, **스킬**, **대상 몹** 전달.
- [ ] 또는 미스트 전용 분기가 필요하면 `script/skill/2111003.lua` 등에서만 분기.

### 4.5 네트워크

- [ ] 미스트 **생성/제거** 클라이언트 패킷 (자바 `MaplePacketCreator.spawnMist` / `removeMist` 참고).
- [ ] 맵 입장 시 기존 미스트 동기화 여부 결정.

### 4.6 상수·스킬 ID

- [ ] `constants_skill.go` 에 이미 `PoisonMist: 2111003` 존재 — 핸들러에서 분기 시 사용.

### 4.7 테스트·검증

- [ ] 단일 맵에 미스트 1개: 만료 시 오브젝트·타이머 정리.
- [ ] prop·중복 POISON 스킵·영역 밖 몹 미적용.

---

## 5. 구현 순서 제안

1. **맵에 미스트 엔티티 + 수명 만료**만 먼저 넣고 스폰/제거 패킷까지 (시각·동기화).
2. **주기 틱** + 영역 내 몹 + `ApplyMobStatus(POISON)` + 기존 **독 틱 데미지 타이머**와 연동.
3. **Lua pDam** 및 원소/앰프 공식을 미스트 적용 경로에서 재사용.
4. 나머지 미스트 스킬 ID(12111005, 4221006, 14111006 등)는 같은 프레임에 확장.

---

## 6. 참고 (자바 위치)

| 내용 | 파일 |
|------|------|
| 2111003 비버프 | `old/src/client/Skill.java` |
| `isMist`, `applyMist` | `old/src/server/MapleStatEffect.java` |
| 미스트 생성자, `isPoisonMist` | `old/src/server/maps/MapleMist.java` |
| 스폰·주기·만료 | `old/src/server/maps/MapleMap.java` (`spawnMist`) |
