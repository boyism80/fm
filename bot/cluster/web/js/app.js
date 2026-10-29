'use strict';

const WORD = {
  pending: '대기',
  running: '실행 중',
  passed: '통과',
  failed: '실패',
  skipped: '건너뜀',
  stopped: '중지',
  lost: '끊김',
};
const DONE = ['passed', 'failed', 'skipped', 'stopped', 'lost'];
const AGENT_WORD = { run: '실행 중', idle: '대기', slow: '보고 지연', lost: '끊김' };
const AGENT_CLS = { run: 'running', idle: 'passed', slow: 'running', lost: 'failed' };

const $ = id => document.getElementById(id);
const esc = value =>
  String(value ?? '').replace(/[&<>"']/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[c]);
const clock = ms => {
  const total = Math.max(0, Math.floor((ms || 0) / 1000));
  const h = Math.floor(total / 3600),
    m = Math.floor(total / 60) % 60,
    s = total % 60;
  const pad = v => String(v).padStart(2, '0');
  return h ? `${h}:${pad(m)}:${pad(s)}` : `${pad(m)}:${pad(s)}`;
};
const seconds = value => (value >= 60 ? clock(value * 1000) : `${value.toFixed(1)}s`);
const local = iso => (iso ? new Date(iso).toLocaleString('ko-KR', { hour12: false }) : '-');
const minutes = ms => (ms >= 60000 ? `${Math.round(ms / 60000)}분` : `${Math.round(ms / 1000)}초`);

const store = {
  get(key, fallback) {
    try {
      const v = localStorage.getItem(key);
      return v === null ? fallback : JSON.parse(v);
    } catch {
      return fallback;
    }
  },
  set(key, value) {
    try {
      localStorage.setItem(key, JSON.stringify(value));
    } catch {}
  },
};

const S = {
  catalog: null,
  run: null,
  agents: [],
  offset: 0,
  online: false,
  selected: store.get('fm.dashboard.selected', null) && new Set(store.get('fm.dashboard.selected', [])),
  filter: 'all',
  search: '',
  view: 'run',
  lastState: null,
};

const serverNow = () => Date.now() + S.offset;
const since = iso => serverNow() - Date.parse(iso);
const active = () => !!S.run && S.run.state !== 'done';
const live = () => S.agents.filter(a => a.state !== 'lost');

async function api(path, body) {
  const options =
    body === undefined
      ? {}
      : { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) };
  const res = await fetch(path, options);
  const json = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(json.error || res.statusText);
  return json;
}

function notice(text, bad) {
  const box = $('notice');
  box.hidden = !text;
  box.textContent = text || '';
  box.classList.toggle('bad', !!bad);
}

function counts(suites) {
  const c = { pending: 0, running: 0, passed: 0, failed: 0, skipped: 0, stopped: 0, lost: 0 };
  suites.forEach(s => c[s.state]++);
  return c;
}

function pill(suite) {
  const time =
    suite.state === 'running' && suite.started_at
      ? ` <span data-since="${esc(suite.started_at)}"></span>`
      : DONE.includes(suite.state) && suite.seconds
        ? ` ${seconds(suite.seconds)}`
        : '';
  const word = suite.infra ? '인프라 실패' : WORD[suite.state];
  const cls = suite.state === 'lost' ? 'failed' : suite.state;
  return `<span class="pill dot ${cls}">${word}${time}</span>`;
}

const titled = s => (s.title && s.title !== s.name ? `<span class="sub">${esc(s.title)}</span>` : '');

function selectedNames() {
  return S.catalog ? S.catalog.suites.map(s => s.name).filter(n => S.selected.has(n)) : [];
}

function visible() {
  const query = S.search.trim().toLowerCase();
  return S.catalog.suites
    .map((s, i) => [s, i])
    .filter(([s]) => {
      if (S.filter === 'parallel' && s.serial) return false;
      if (S.filter === 'serial' && !s.serial) return false;
      if (S.filter === 'explicit' && !s.explicit) return false;
      return !query || s.name.toLowerCase().includes(query) || s.title.toLowerCase().includes(query);
    });
}

function renderCatalog() {
  const catalog = S.catalog;
  if (!catalog) return;

  const all = catalog.suites;
  const filters = [
    ['all', '전체', all.length],
    ['parallel', 'parallel', all.filter(s => !s.serial).length],
    ['serial', 'serial', all.filter(s => s.serial).length],
    ['explicit', 'explicit', all.filter(s => s.explicit).length],
  ];
  $('filters').innerHTML = filters
    .map(
      ([key, word, n]) =>
        `<button class="${S.filter === key ? 'on' : ''}" data-filter="${key}">${word} (${n})</button>`,
    )
    .join('');

  const results = new Map((S.run ? S.run.suites : []).map(s => [s.name, s]));
  const rows = visible();
  $('suites').innerHTML = rows.length
    ? rows
        .map(([s, i]) => {
          const result = results.get(s.name);
          return `<div class="suite${result ? ` state-${result.state}` : ''}" data-suite="${esc(s.name)}">
<span class="box${S.selected.has(s.name) ? ' on' : ''}"></span>
<span class="order">#${String(i + 1).padStart(2, '0')}</span>
<div><div class="name"><b>${esc(s.name)}</b><span class="pill ${s.serial ? 'serial' : 'parallel'}">${s.serial ? 'serial' : 'parallel'}</span>${s.explicit ? '<span class="pill explicit">explicit</span>' : ''}<span class="meta">봇 ${s.bots} · ${minutes(s.timeout_ms)}</span></div>
<div class="title${s.error ? ' bad' : ''}">${esc(s.error || (s.title !== s.name ? s.title : ''))}</div></div>
<div class="end">${result ? pill(result) : ''}</div></div>`;
        })
        .join('')
    : '<div class="empty">조건에 맞는 스위트가 없습니다</div>';

  $('tab-pick').textContent = `${S.selected.size}/${all.length}`;
  const broken = all.filter(s => s.error).length;
  const time = new Date(catalog.loaded_at).toLocaleTimeString('ko-KR', { hour12: false });
  $('loaded').textContent = `${time} 읽음 · ${all.length}개${broken ? ` · 오류 ${broken}` : ''}`;
  $('foot').textContent = `agents ${live().length} · suites ${all.length}`;
  updateButtons();
  tick();
}

function updateButtons() {
  const n = selectedNames().length;
  $('start').textContent = live().length ? `선택한 ${n}개 실행` : '연결된 Agent가 없습니다';
  $('start').disabled = !S.online || active() || n === 0 || live().length === 0;
  $('stop').disabled = !S.run || S.run.state !== 'running';
  $('reload').disabled = active();
  $('pick-failed').disabled = !S.run || !S.run.suites.some(s => s.state === 'failed' || s.state === 'lost');
}

function renderStatus() {
  const run = S.run;
  const badge = $('status');
  let state = 'idle',
    word = 'IDLE';
  if (run && run.state === 'running') [state, word] = ['running', 'RUNNING'];
  else if (run && run.state === 'stopping') [state, word] = ['stopping', 'STOPPING'];
  else if (run && run.suites.some(s => s.state === 'failed' || s.state === 'lost')) [state, word] = ['failed', 'FAILED'];
  else if (run && run.suites.some(s => s.state === 'stopped')) [state, word] = ['stopped', 'STOPPED'];
  else if (run) [state, word] = ['passed', 'PASSED'];
  badge.className = `badge ${state}`;
  badge.querySelector('span').textContent = word;
}

function renderRun() {
  const run = S.run;
  const suites = run ? run.suites : [];
  const c = counts(suites);
  const total = suites.length;
  const done = DONE.reduce((t, k) => t + c[k], 0);

  renderStatus();

  const pct = k => (total ? (c[k] / total) * 100 : 0);
  $('ribbon')
    .querySelectorAll('i')
    .forEach(bar => (bar.style.width = `${bar.className === 'failed' ? pct('failed') + pct('lost') : pct(bar.className)}%`));
  $('ribbon-text').textContent = run
    ? `${done} / ${total} suites (${total ? ((done / total) * 100).toFixed(1) : 0}%)`
    : '실행 전';
  $('run-id').textContent = run ? `run ${run.id}` : '';
  $('legend').innerHTML = run
    ? ['passed', 'failed', 'running', 'pending', 'lost', 'skipped', 'stopped']
        .filter(k => c[k] || ['passed', 'failed', 'running', 'pending'].includes(k))
        .map(k => `<span class="pill dot ${k === 'lost' ? 'failed' : k}">${WORD[k]} ${c[k]} (${pct(k).toFixed(1)}%)</span>`)
        .join('')
    : '<span class="faint">스위트를 골라 실행하면 진행률이 여기에 표시됩니다</span>';

  $('k-progress').innerHTML = `${done}<span class="sep">/</span>${total}`;
  $('k-progress-sub').textContent = run ? `실행 중 ${c.running} · 대기 ${c.pending}` : '-';
  $('k-result').innerHTML = `<span class="good">${c.passed}</span><span class="sep">·</span><span class="bad">${c.failed + c.lost}</span>`;
  const judged = c.passed + c.failed + c.lost;
  $('k-result-sub').textContent = run
    ? `성공률 ${judged ? ((c.passed / judged) * 100).toFixed(1) : '0.0'}% · 끊김 ${c.lost} · 중지 ${c.stopped}`
    : '-';
  const agents = live();
  const seats = agents.reduce((t, a) => t + a.seats, 0);
  const busy = agents.reduce((t, a) => t + a.suites.length, 0);
  $('k-seats').innerHTML = `${agents.length}<span class="sep">·</span>${busy}<span class="sep">/</span>${seats}`;
  $('k-seats-sub').textContent = `Agent ${agents.length}대 · 빈 seat ${seats - busy}${S.agents.length > agents.length ? ` · 끊김 ${S.agents.length - agents.length}` : ''}`;
  $('k-elapsed-sub').textContent = run ? `${local(run.started_at)} 시작` : '-';

  $('tab-progress').textContent = run ? `${done}/${total}` : '';
  $('tab-progress').parentElement.classList.toggle('live', active());
  $('tab-agents').textContent = S.agents.length ? `${agents.length}/${S.agents.length}` : '';

  renderBoard(run);
  renderPlan(run, done, total, c);
  renderAgents();
  renderFailures(run);
  renderCatalog();
}

function renderBoard(run) {
  const suites = run ? run.suites : [];
  const running = suites.filter(s => s.state === 'running');
  const tile = (label, suite, extra) => {
    if (!suite) return `<div class="seat idle${extra}"><div class="no">${label}</div><div class="who">빈 자리</div></div>`;
    return `<div class="seat${extra}" title="${esc(suite.title)}"><div class="no">${label}</div><div class="who">${esc(suite.name)}</div>
<div class="time" data-since="${esc(suite.started_at)}"></div><i class="bar" data-since="${esc(suite.started_at)}" data-timeout="${suite.timeout_ms || 0}"></i></div>`;
  };
  const groups = live().map(a => {
    const tiles = [];
    for (let i = 1; i <= a.seats; i++) {
      const seat = a.seat_base + i;
      const suite = running.find(s => s.agent === a.name && s.seat === seat);
      tiles.push(tile(`SEAT ${String(seat).padStart(2, '0')}`, suite, suite && suite.serial ? ' serial-run' : ''));
    }
    return `<div class="agent-group"><div class="agent-head"><b>${esc(a.name)}</b><span class="faint mono">${esc(a.host)} · seat ${a.seat_base + 1}~${a.seat_base + a.seats}</span>
<span class="pill dot ${AGENT_CLS[a.state]}">${AGENT_WORD[a.state]}</span></div><div class="board">${tiles.join('')}</div></div>`;
  });
  $('board').innerHTML = groups.length
    ? groups.join('')
    : '<div class="empty">연결된 Agent가 없습니다. <span class="mono">-mode agent</span>로 Agent를 띄우세요</div>';
  const serialLeft = suites.filter(s => s.serial && s.state === 'pending').length;
  const serialRun = running.find(s => s.serial);
  $('board-note').textContent = serialRun
    ? `serial 실행 중: ${serialRun.name} · 대기 ${serialLeft}개`
    : run
      ? `serial 대기 ${serialLeft}개`
      : '실행 중인 테스트가 없습니다';
}

function renderPlan(run, done, total, c) {
  $('plan-note').textContent = run
    ? `${done} / ${total}개 종료 · ${run.state === 'done' ? '끝남' : c.running ? '진행 중' : 'Agent 대기 중'} · 재시도 ${run.retry}회`
    : '-';
  if (!run || !run.suites.length) {
    $('plan').innerHTML = '<tr><td colspan="7" class="empty">실행 계획이 없습니다</td></tr>';
    return;
  }
  $('plan').innerHTML = run.suites
    .map(
      s => `<tr class="${s.state === 'running' ? 'active' : ''}">
<td><span class="pill ${s.serial ? 'serial' : 'parallel'}">${s.serial ? 'SERIAL' : 'PARALLEL'}</span></td>
<td><span class="mono">${esc(s.name)}</span>${titled(s)}</td>
<td>${pill(s)}</td><td class="mono">${esc(s.agent || '—')}</td><td class="num">${s.seat || '—'}</td><td class="num">${s.attempts}</td>
<td class="num">${s.state === 'running' ? `<span data-since="${esc(s.started_at)}"></span>` : s.seconds ? seconds(s.seconds) : '-'}</td></tr>`,
    )
    .join('');
}

function renderAgents() {
  const agents = S.agents;
  const word = ['run', 'idle', 'slow', 'lost']
    .map(k => [k, agents.filter(a => a.state === k).length])
    .filter(([, n]) => n)
    .map(([k, n]) => `${AGENT_WORD[k]} ${n}`)
    .join(' · ');
  $('agents-note').textContent = agents.length ? word : '-';
  $('agents').innerHTML = agents.length
    ? agents
        .map(
          a => `<tr><td class="mono">${esc(a.name)}<span class="sub">#${a.id} · pid ${a.pid}</span></td><td class="mono faint">${esc(a.host)}</td>
<td><span class="pill dot ${AGENT_CLS[a.state]}">${AGENT_WORD[a.state]}</span></td>
<td class="mono">${a.seat_base + 1}~${a.seat_base + a.seats} <span class="faint">(${a.seats})</span></td>
<td class="mono" title="${esc(a.suites.join(', '))}">${a.suites.length ? esc(a.suites.length === 1 ? a.suites[0] : `${a.suites[0]} 외 ${a.suites.length - 1}`) : '—'}</td>
<td class="num">${a.state === 'lost' ? '-' : `${Math.floor(a.seen)}초 전`}</td>
<td>${a.state === 'lost' ? '' : `<button class="link" data-quit="${a.id}">종료</button>`}</td></tr>`,
        )
        .join('')
    : '<tr><td colspan="7" class="empty">연결된 Agent가 없습니다</td></tr>';
}

function failureCards(suites) {
  return suites
    .filter(s => s.state === 'failed' || s.state === 'lost')
    .map(
      s => `<details class="fail"><summary>${pill(s)}<b>${esc(s.name)}</b><span class="faint">${esc(s.title && s.title !== s.name ? s.title : '')}</span>
<span class="faint mono">${esc(s.agent || '-')} · seat ${s.seat || '-'} · ${s.seconds ? seconds(s.seconds) : '-'}</span></summary>
<div class="body"><ol class="msgs">${(s.failures || []).map(f => `<li>${esc(f)}</li>`).join('')}</ol>
${s.log && s.log.length ? `<pre class="log">${esc(s.log.join('\n'))}</pre>` : ''}</div></details>`,
    )
    .join('');
}

function renderFailures(run) {
  const failed = run ? run.suites.filter(s => s.state === 'failed' || s.state === 'lost') : [];
  $('tab-failures').textContent = failed.length || '';
  const open = new Set([...$('failures').querySelectorAll('details[open] b')].map(b => b.textContent));
  $('failures').innerHTML = failed.length
    ? failureCards(run.suites)
    : `<div class="card empty-card">${run ? '실패한 스위트가 없습니다' : '아직 실행 기록이 없습니다'}</div>`;
  $('failures')
    .querySelectorAll('details')
    .forEach(d => (d.open = open.has(d.querySelector('b').textContent)));
}

function tick() {
  const run = S.run;
  if (run) {
    const elapsed =
      run.state === 'done' ? Date.parse(run.finished_at) - Date.parse(run.started_at) : since(run.started_at);
    $('status-clock').textContent = clock(elapsed);
    $('k-elapsed').textContent = clock(elapsed);
  } else {
    $('status-clock').textContent = '00:00';
    $('k-elapsed').textContent = '-';
  }
  document.querySelectorAll('[data-since]').forEach(el => {
    const ms = since(el.dataset.since);
    if (el.dataset.timeout !== undefined) {
      const timeout = Number(el.dataset.timeout);
      el.style.width = timeout ? `${Math.min(100, (ms / timeout) * 100)}%` : '0%';
    } else {
      el.textContent = clock(ms);
    }
  });
}

async function loadCatalog() {
  try {
    S.catalog = await api('/api/suites');
    const known = new Set(S.catalog.suites.map(s => s.name));
    S.selected =
      S.selected === null
        ? new Set(S.catalog.suites.filter(s => !s.explicit).map(s => s.name))
        : new Set([...S.selected].filter(n => known.has(n)));
    store.set('fm.dashboard.selected', [...S.selected]);
    renderCatalog();
  } catch (error) {
    S.selected = S.selected instanceof Set ? S.selected : new Set();
    notice(`스위트 목록을 읽지 못했습니다: ${error.message}`, true);
  }
}

function toggle(names, on) {
  names.forEach(n => (on ? S.selected.add(n) : S.selected.delete(n)));
  store.set('fm.dashboard.selected', [...S.selected]);
  renderCatalog();
}

function setRetry(value) {
  const retry = Math.max(0, Math.round(Number(value) || 0));
  $('retry').value = retry;
  store.set('fm.dashboard.retry', retry);
}

async function loadHistory() {
  try {
    const runs = await api('/api/runs');
    $('runs').innerHTML = runs.length
      ? runs
          .map(
            r => `<tr data-run="${esc(r.id)}"><td class="mono">${esc(r.id)}</td><td>${local(r.started_at)}</td>
<td class="num">${clock(Date.parse(r.finished_at) - Date.parse(r.started_at))}</td><td class="num">${r.agents}</td>
<td><span class="pill passed">통과 ${r.passed}</span> <span class="pill ${r.failed ? 'failed' : 'stopped'}">실패 ${r.failed}</span>${r.lost ? ` <span class="pill failed">끊김 ${r.lost}</span>` : ''}${r.skipped ? ` <span class="pill skipped">건너뜀 ${r.skipped}</span>` : ''}${r.stopped ? ` <span class="pill stopped">중지 ${r.stopped}</span>` : ''}</td></tr>`,
          )
          .join('')
      : '<tr><td colspan="5" class="empty">끝난 실행이 없습니다</td></tr>';
  } catch (error) {
    $('runs').innerHTML = `<tr><td colspan="5" class="empty bad">${esc(error.message)}</td></tr>`;
  }
}

async function showRun(id) {
  $('runs')
    .querySelectorAll('tr')
    .forEach(r => r.classList.toggle('on', r.dataset.run === id));
  const run = await api(`/api/runs/${encodeURIComponent(id)}`);
  const c = counts(run.suites);
  const agents = new Set(run.suites.map(s => s.agent).filter(Boolean));
  $('detail').hidden = false;
  $('detail-title').textContent = run.id;
  $('detail-sub').textContent = `${local(run.started_at)} ~ ${local(run.finished_at)} · ${clock(Date.parse(run.finished_at) - Date.parse(run.started_at))} · Agent ${agents.size}대 · 통과 ${c.passed} · 실패 ${c.failed} · 끊김 ${c.lost}`;
  $('detail-file').href = `/api/runs/${encodeURIComponent(run.id)}`;
  $('detail-file').download = `summary-${run.id}.json`;
  $('detail-suites').innerHTML = run.suites
    .map(
      s => `<tr><td><span class="mono">${esc(s.name)}</span>${titled(s)}</td><td>${pill(s)}</td><td class="mono">${esc(s.agent || '—')}</td>
<td class="num">${s.seat || '—'}</td><td class="num">${s.attempts}</td><td class="num">${s.seconds ? seconds(s.seconds) : '-'}</td></tr>`,
    )
    .join('');
  $('detail-failures').innerHTML = failureCards(run.suites);
  $('detail').scrollIntoView({ behavior: 'smooth', block: 'start' });
}

function show(view) {
  S.view = view;
  document.querySelectorAll('.view').forEach(v => (v.hidden = v.id !== `view-${view}`));
  document.querySelectorAll('.views a').forEach(a => a.classList.toggle('on', a.dataset.view === view));
  history.replaceState(null, '', `#${view}`);
  if (view === 'history') loadHistory();
}

function connect() {
  const events = new EventSource('/api/events');
  const online = on => {
    S.online = on;
    $('online').textContent = on ? 'Master 연결됨' : 'Master 연결 끊김';
    $('online').classList.toggle('off', !on);
    updateButtons();
  };
  events.onopen = () => online(true);
  events.onerror = () => online(false);
  events.onmessage = event => {
    const message = JSON.parse(event.data);
    S.offset = Date.parse(message.now) - Date.now();
    S.run = message.run;
    S.agents = message.agents || [];
    const state = S.run ? S.run.state : null;
    if (state === 'done' && S.lastState && S.lastState !== 'done' && S.view === 'history') loadHistory();
    S.lastState = state;
    if (!S.online) online(true);
    renderRun();
  };
}

function bind() {
  document.querySelectorAll('.views a').forEach(a => a.addEventListener('click', () => show(a.dataset.view)));
  document.querySelector('.tabs').addEventListener('click', e => {
    const button = e.target.closest('[data-tab]');
    if (!button) return;
    document.querySelectorAll('.tabs [data-tab]').forEach(x => x.setAttribute('aria-selected', String(x === button)));
    document.querySelectorAll('[data-panel]').forEach(p => (p.hidden = p.dataset.panel !== button.dataset.tab));
  });

  $('suites').addEventListener('click', e => {
    const row = e.target.closest('[data-suite]');
    if (row) toggle([row.dataset.suite], !S.selected.has(row.dataset.suite));
  });
  $('filters').addEventListener('click', e => {
    const button = e.target.closest('[data-filter]');
    if (!button) return;
    S.filter = button.dataset.filter;
    renderCatalog();
  });
  $('pick-all').addEventListener('click', () => {
    if (!S.catalog) return;
    toggle(
      visible()
        .map(([s]) => s)
        .filter(s => S.filter === 'explicit' || !s.explicit)
        .map(s => s.name),
      true,
    );
  });
  $('pick-none').addEventListener('click', () => S.catalog && toggle([...S.selected], false));
  $('pick-failed').addEventListener('click', () => {
    if (!S.run) return;
    S.selected = new Set(S.run.suites.filter(s => s.state === 'failed' || s.state === 'lost').map(s => s.name));
    store.set('fm.dashboard.selected', [...S.selected]);
    renderCatalog();
  });
  $('search').addEventListener('input', e => {
    S.search = e.target.value;
    renderCatalog();
  });

  $('retry-down').addEventListener('click', () => setRetry(Number($('retry').value) - 1));
  $('retry-up').addEventListener('click', () => setRetry(Number($('retry').value) + 1));
  $('retry').addEventListener('change', e => setRetry(e.target.value));

  $('start').addEventListener('click', async () => {
    $('start').disabled = true;
    try {
      await api('/api/runs', { suites: selectedNames(), retry: Number($('retry').value) });
      notice('');
      document.querySelector('.tabs [data-tab="progress"]').click();
    } catch (error) {
      notice(error.message, true);
    }
    updateButtons();
  });
  $('stop').addEventListener('click', () =>
    api('/api/runs/current/stop', {}).catch(error => notice(error.message, true)),
  );
  $('reload').addEventListener('click', async () => {
    await loadCatalog();
    notice(`registry.lua와 스위트 파일을 다시 읽었습니다 · 스위트 ${S.catalog ? S.catalog.suites.length : 0}개`);
  });
  $('agents').addEventListener('click', async e => {
    const button = e.target.closest('[data-quit]');
    if (!button) return;
    button.disabled = true;
    try {
      await api(`/api/agents/${button.dataset.quit}/quit`, {});
    } catch (error) {
      notice(error.message, true);
      button.disabled = false;
    }
  });
  $('runs').addEventListener('click', e => {
    const row = e.target.closest('[data-run]');
    if (row) showRun(row.dataset.run).catch(error => notice(error.message, true));
  });
}

bind();
setRetry(store.get('fm.dashboard.retry', 1));
show(location.hash === '#history' ? 'history' : 'run');
setInterval(tick, 250);
loadCatalog();
connect();
