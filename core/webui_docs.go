package core

const adminAPIDocsHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>cx-connect API 文档</title>
  <style>
    :root {
      --bg: #f5efe5;
      --panel: rgba(255, 252, 247, 0.92);
      --ink: #1b1d21;
      --muted: #5e6472;
      --line: rgba(27, 29, 33, 0.12);
      --brand: #c8553d;
      --brand-deep: #7f2f1d;
      --code: #182026;
      --shadow: 0 20px 60px rgba(69, 51, 31, 0.14);
      --radius: 22px;
    }

    * { box-sizing: border-box; }

    body {
      margin: 0;
      color: var(--ink);
      font-family: "IBM Plex Sans", "PingFang SC", "Noto Sans SC", sans-serif;
      background:
        radial-gradient(circle at top left, rgba(200, 85, 61, 0.16), transparent 28%),
        linear-gradient(180deg, #fbf6ee 0%, #f2ebdf 100%);
    }

    .shell {
      width: min(1120px, calc(100vw - 32px));
      margin: 24px auto 40px;
      display: grid;
      gap: 18px;
    }

    .hero, .section {
      background: var(--panel);
      border: 1px solid rgba(255,255,255,0.7);
      border-radius: var(--radius);
      box-shadow: var(--shadow);
      backdrop-filter: blur(16px);
    }

    .hero {
      padding: 24px;
      display: flex;
      justify-content: space-between;
      gap: 16px;
      align-items: end;
    }

    h1 {
      margin: 0;
      font-size: clamp(28px, 4vw, 42px);
      line-height: 0.95;
      letter-spacing: -0.05em;
    }

    p {
      margin: 10px 0 0;
      color: var(--muted);
      line-height: 1.6;
    }

    .nav {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
    }

    .button-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      text-decoration: none;
      border-radius: 999px;
      padding: 11px 16px;
      background: rgba(27, 29, 33, 0.08);
      color: var(--ink);
      transition: transform 120ms ease;
    }

    .button-link.accent {
      background: linear-gradient(135deg, var(--brand) 0%, var(--brand-deep) 100%);
      color: #fff;
    }

    .button-link:hover { transform: translateY(-1px); }

    .section {
      padding: 22px 24px 24px;
      display: grid;
      gap: 14px;
    }

    h2 {
      margin: 0;
      font-size: 18px;
    }

    h3 {
      margin: 6px 0 0;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: 0.12em;
      color: var(--muted);
    }

    .meta {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
    }

    .badge {
      display: inline-flex;
      align-items: center;
      padding: 5px 10px;
      border-radius: 999px;
      background: rgba(27, 29, 33, 0.06);
      color: var(--muted);
      font-size: 12px;
    }

    code.inline {
      padding: 2px 8px;
      border-radius: 999px;
      background: rgba(27, 29, 33, 0.06);
      font-family: "IBM Plex Mono", "SFMono-Regular", monospace;
      font-size: 13px;
    }

    pre {
      margin: 0;
      padding: 16px;
      overflow: auto;
      border-radius: 18px;
      background: var(--code);
      color: #ecf2f8;
      font-family: "IBM Plex Mono", "SFMono-Regular", monospace;
      font-size: 13px;
      line-height: 1.6;
    }

    ul {
      margin: 0;
      padding-left: 20px;
      color: var(--muted);
      line-height: 1.7;
    }

    .grid {
      display: grid;
      gap: 18px;
    }

    @media (max-width: 840px) {
      .hero { flex-direction: column; align-items: start; }
    }
  </style>
</head>
<body>
  <div class="shell">
    <section class="hero">
      <div>
        <h1>Web 管理 API</h1>
        <p>这是一套本地 HTTP API，用于创建 Web Session、列出会话、查看详情，以及向指定 session 追加新的 prompt。默认配合内置 Web 管理页面一起使用。</p>
      </div>
      <div class="nav">
        <a class="button-link" href="/">返回首页</a>
      </div>
    </section>

    <section class="section">
      <h2>基本说明</h2>
      <ul>
        <li>默认监听你配置的本地 Web 地址，例如 <code class="inline">http://127.0.0.1:6380</code>。</li>
        <li>创建 Session 时不需要自己传 <code class="inline">session_key</code>，后端会自动生成 <code class="inline">web:xxxxxx</code>。</li>
        <li>当前 Web 页面使用的就是下面这组接口。</li>
      </ul>
    </section>

    <section class="section">
      <h3>GET</h3>
      <h2>/api/admin/sessions</h2>
      <div class="meta">
        <span class="badge">用途：列出所有 session 分组</span>
      </div>
      <pre>{
  "project": "default",
  "groups": [
    {
      "session_key": "web:a1b2c3d4e5f6",
      "platform": "web",
      "active_session_id": "s3",
      "interactive": true,
      "sessions": [
        {
          "id": "s3",
          "name": "修复登录问题",
          "work_dir": "/Users/me/project",
          "agent_session_id": "thread_xxx",
          "history_count": 12,
          "created_at": "2026-03-11T15:00:00+08:00",
          "updated_at": "2026-03-11T15:08:00+08:00",
          "busy": false,
          "active": true
        }
      ]
    }
  ]
}</pre>
    </section>

    <section class="section">
      <h3>GET</h3>
      <h2>/api/admin/session?session_key=web:...&session_id=s3</h2>
      <div class="meta">
        <span class="badge">用途：查看单个 session 的完整详情</span>
      </div>
      <pre>{
  "project": "default",
  "session_key": "web:a1b2c3d4e5f6",
  "platform": "web",
  "active_session_id": "s3",
  "interactive": true,
  "session": {
    "id": "s3",
    "name": "修复登录问题",
    "work_dir": "/Users/me/project",
    "agent_session_id": "thread_xxx",
    "history_count": 12,
    "created_at": "2026-03-11T15:00:00+08:00",
    "updated_at": "2026-03-11T15:08:00+08:00",
    "busy": false,
    "history": [
      {
        "role": "user",
        "content": "先看一下登录失败的原因",
        "timestamp": "2026-03-11T15:01:00+08:00"
      }
    ]
  }
}</pre>
    </section>

    <section class="section">
      <h3>POST</h3>
      <h2>/api/admin/session/create</h2>
      <div class="meta">
        <span class="badge">用途：创建新的 Web Session</span>
      </div>
      <pre>{
  "name": "修复登录问题",
  "work_dir": "/Users/me/project"
}</pre>
      <p>字段说明：</p>
      <ul>
        <li><code class="inline">name</code> 可选，留空会自动生成。</li>
        <li><code class="inline">work_dir</code> 可选，留空时会使用默认目录 <code class="inline">~/.cx-connect/workspace/&lt;session-name&gt;</code>。</li>
      </ul>
      <pre>{
  "project": "default",
  "session_key": "web:a1b2c3d4e5f6",
  "platform": "web",
  "active_session_id": "s3",
  "interactive": false,
  "display_work_dir": "~/.cx-connect/workspace/修复登录问题",
  "session": {
    "id": "s3",
    "name": "修复登录问题",
    "work_dir": "/Users/me/.cx-connect/workspace/修复登录问题",
    "agent_session_id": "",
    "history_count": 0,
    "created_at": "2026-03-11T15:00:00+08:00",
    "updated_at": "2026-03-11T15:00:00+08:00",
    "busy": false,
    "history": []
  }
}</pre>
    </section>

    <section class="section">
      <h3>POST</h3>
      <h2>/api/admin/prompt</h2>
      <div class="meta">
        <span class="badge">用途：向指定 session 追加新的 prompt</span>
      </div>
      <pre>{
  "session_key": "web:a1b2c3d4e5f6",
  "session_id": "s3",
  "prompt": "继续分析登录流程，并给出修复方案。"
}</pre>
      <pre>{
  "status": "accepted",
  "project": "default"
}</pre>
      <ul>
        <li>调用成功后会异步把 prompt 送进对应 session。</li>
        <li>如果你传入的不是当前 active session，后端会先切换到该 session，再发送 prompt。</li>
      </ul>
    </section>

    <section class="section">
      <h2>curl 示例</h2>
      <div class="grid">
        <pre>curl http://127.0.0.1:6380/api/admin/sessions</pre>
        <pre>curl -X POST http://127.0.0.1:6380/api/admin/session/create \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "新会话",
    "work_dir": "/Users/me/project"
  }'</pre>
        <pre>curl "http://127.0.0.1:6380/api/admin/session?session_key=web:abc123&session_id=s1"</pre>
        <pre>curl -X POST http://127.0.0.1:6380/api/admin/prompt \
  -H 'Content-Type: application/json' \
  -d '{
    "session_key": "web:abc123",
    "session_id": "s1",
    "prompt": "继续上一次的任务"
  }'</pre>
      </div>
    </section>
  </div>
</body>
</html>
`
