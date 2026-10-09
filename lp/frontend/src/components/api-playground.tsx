"use client";

import { useState } from "react";

type TabAccent = {
  active: string;
  underline: string;
  post: string;
  string: string;
  number: string;
  responseText: string;
};

type ApiTab = {
  id: string;
  label: string;
  path: string;
  stats: { ms: string; tokens: string; cost: string };
  accent: TabAccent;
  request: React.ReactNode;
  response: React.ReactNode;
};

const accents: Record<string, TabAccent> = {
  chat: {
    active: "text-white",
    underline: "bg-success",
    post: "bg-success/15 text-success",
    string: "text-[#e8c27a]",
    number: "text-[#b8e986]",
    responseText: "text-success",
  },
  responses: {
    active: "text-[#e8b84a]",
    underline: "bg-[#e8b84a]",
    post: "bg-[#e8b84a]/15 text-[#e8b84a]",
    string: "text-[#e8c27a]",
    number: "text-[#e8b84a]",
    responseText: "text-[#e8c27a]",
  },
  claude: {
    active: "text-[#6aa8ff]",
    underline: "bg-[#6aa8ff]",
    post: "bg-[#6aa8ff]/15 text-[#6aa8ff]",
    string: "text-[#e8c27a]",
    number: "text-[#c39bff]",
    responseText: "text-[#6aa8ff]",
  },
  gemini: {
    active: "text-[#c39bff]",
    underline: "bg-[#c39bff]",
    post: "bg-[#c39bff]/15 text-[#c39bff]",
    string: "text-[#e8c27a]",
    number: "text-[#c39bff]",
    responseText: "text-[#c39bff]",
  },
};

const S = "text-[#e8c27a]";
const K = "text-[#7eb6ff]";
const P = "text-[#d0d8e2]";
const G = "text-success";

function ReqChat() {
  return (
    <pre className="mt-4 whitespace-pre-wrap break-words font-mono text-[12.5px] leading-[1.75] text-[#d0d8e2]">
      <code>
        <span className={G}>curl</span>
        <span className={S}> -X POST </span>
        <span className={S}>&quot;/v1/chat/completions&quot;</span>
        {" \\\n"}
        {"  "}
        <span className={K}>-H</span> <span className={S}>&quot;Authorization: Bearer sk-••••&quot;</span>
        {" \\\n"}
        {"  "}
        <span className={K}>-d</span> <span className={S}>{`'{`}</span>
        {"\n"}
        {"    "}
        <span className={K}>&quot;model&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;your-model&quot;</span>
        <span className={P}>,</span>
        {"\n"}
        {"    "}
        <span className={K}>&quot;messages&quot;</span>
        <span className={P}>: [</span>
        {"\n"}
        {"      "}
        <span className={P}>{`{ `}</span>
        <span className={K}>&quot;role&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;user&quot;</span>
        <span className={P}>, </span>
        <span className={K}>&quot;content&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;...&quot;</span>
        <span className={P}>{` }`}</span>
        {"\n"}
        {"    "}
        <span className={P}>]</span>
        {"\n"}
        <span className={S}>{`}'`}</span>
      </code>
    </pre>
  );
}

function ResChat({ a }: { a: TabAccent }) {
  return (
    <pre className="mt-4 whitespace-pre-wrap break-words font-mono text-[12.5px] leading-[1.75] text-[#c6d0dc]">
      <code>
        {"{\n"}
        {"  "}
        <span className={K}>&quot;choices&quot;</span>
        <span className={P}>: [{"{"}</span> <span className={K}>&quot;message&quot;</span>
        <span className={P}>: {"{"}</span> <span className={K}>&quot;content&quot;</span>
        <span className={P}>: </span>
        <span className={a.responseText}>&quot;Chat request routed.&quot;</span>
        <span className={P}>{` } }]`}</span>
        {"\n"}
        {"  "}
        <span className={K}>&quot;usage&quot;</span>
        <span className={P}>: {"{"}</span> <span className={K}>&quot;total_tokens&quot;</span>
        <span className={P}>: </span>
        <span className={a.number}>27</span>
        <span className={P}>{` }`}</span>
        {"\n"}
        {"}"}
      </code>
    </pre>
  );
}

function ReqResponses() {
  return (
    <pre className="mt-4 whitespace-pre-wrap break-words font-mono text-[12.5px] leading-[1.75] text-[#d0d8e2]">
      <code>
        <span className={G}>curl</span>
        <span className={S}> -X POST </span>
        <span className={S}>&quot;/v1/responses&quot;</span>
        {" \\\n"}
        {"  "}
        <span className={K}>-H</span> <span className={S}>&quot;Authorization: Bearer sk-••••&quot;</span>
        {" \\\n"}
        {"  "}
        <span className={K}>-d</span> <span className={S}>{`'{`}</span>
        {"\n"}
        {"    "}
        <span className={K}>&quot;model&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;your-model&quot;</span>
        <span className={P}>,</span>
        {"\n"}
        {"    "}
        <span className={K}>&quot;input&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;...&quot;</span>
        {"\n"}
        {"  "}
        <span className={S}>{`}'`}</span>
      </code>
    </pre>
  );
}

function ResResponses({ a }: { a: TabAccent }) {
  return (
    <pre className="mt-4 whitespace-pre-wrap break-words font-mono text-[12.5px] leading-[1.75] text-[#c6d0dc]">
      <code>
        {"{\n"}
        {"  "}
        <span className={K}>&quot;output&quot;</span>
        <span className={P}>: [{"{"}</span> <span className={K}>&quot;type&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;output_text&quot;</span>
        <span className={P}>, </span>
        <span className={K}>&quot;text&quot;</span>
        <span className={P}>: </span>
        <span className={a.responseText}>&quot;Response workflow ready.&quot;</span>
        <span className={P}>{` } }]`}</span>
        {"\n"}
        {"  "}
        <span className={K}>&quot;usage&quot;</span>
        <span className={P}>: {"{"}</span> <span className={K}>&quot;total_tokens&quot;</span>
        <span className={P}>: </span>
        <span className={a.number}>31</span>
        <span className={P}>{` }`}</span>
        {"\n"}
        {"}"}
      </code>
    </pre>
  );
}

function ReqClaude() {
  return (
    <pre className="mt-4 whitespace-pre-wrap break-words font-mono text-[12.5px] leading-[1.75] text-[#d0d8e2]">
      <code>
        <span className={G}>curl</span>
        <span className={S}> -X POST </span>
        <span className={S}>&quot;/v1/messages&quot;</span>
        {" \\\n"}
        {"  "}
        <span className={K}>-H</span> <span className={S}>&quot;x-api-key: sk-••••&quot;</span>
        {" \\\n"}
        {"  "}
        <span className={K}>-H</span> <span className={S}>&quot;anthropic-version: 2023-06-01&quot;</span>
        {" \\\n"}
        {"  "}
        <span className={K}>-d</span> <span className={S}>{`'{`}</span>
        {"\n"}
        {"    "}
        <span className={K}>&quot;model&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;your-model&quot;</span>
        <span className={P}>,</span>
        {"\n"}
        {"    "}
        <span className={K}>&quot;max_tokens&quot;</span>
        <span className={P}>: </span>
        <span className={S}>1024</span>
        <span className={P}>,</span>
        {"\n"}
        {"    "}
        <span className={K}>&quot;messages&quot;</span>
        <span className={P}>: [</span>
        {"\n"}
        {"      "}
        <span className={P}>{`{ `}</span>
        <span className={K}>&quot;role&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;user&quot;</span>
        <span className={P}>, </span>
        <span className={K}>&quot;content&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;...&quot;</span>
        <span className={P}>{` }`}</span>
        {"\n"}
        {"    "}
        <span className={P}>]</span>
        {"\n"}
        {"  "}
        <span className={S}>{`}'`}</span>
      </code>
    </pre>
  );
}

function ResClaude({ a }: { a: TabAccent }) {
  return (
    <pre className="mt-4 whitespace-pre-wrap break-words font-mono text-[12.5px] leading-[1.75] text-[#c6d0dc]">
      <code>
        {"{\n"}
        {"  "}
        <span className={K}>&quot;content&quot;</span>
        <span className={P}>: [{"{"}</span> <span className={K}>&quot;type&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;text&quot;</span>
        <span className={P}>, </span>
        <span className={K}>&quot;text&quot;</span>
        <span className={P}>: </span>
        <span className={a.responseText}>&quot;Claude message routed.&quot;</span>
        <span className={P}>{` } }]`}</span>
        {"\n"}
        {"  "}
        <span className={K}>&quot;usage&quot;</span>
        <span className={P}>: {"{"}</span> <span className={K}>&quot;input_tokens&quot;</span>
        <span className={P}>: </span>
        <span className={a.number}>11</span>
        <span className={P}>, </span>
        <span className={K}>&quot;output_tokens&quot;</span>
        <span className={P}>: </span>
        <span className={a.number}>18</span>
        <span className={P}>{` }`}</span>
        {"\n"}
        {"}"}
      </code>
    </pre>
  );
}

function ReqGemini() {
  return (
    <pre className="mt-4 whitespace-pre-wrap break-words font-mono text-[12.5px] leading-[1.75] text-[#d0d8e2]">
      <code>
        <span className={G}>curl</span>
        <span className={S}> -X POST </span>
        <span className={S}>&quot;/v1beta/models/{"{model}"}:generateContent&quot;</span>
        {" \\\n"}
        {"  "}
        <span className={K}>-H</span> <span className={S}>&quot;x-goog-api-key: sk-••••&quot;</span>
        {" \\\n"}
        {"  "}
        <span className={K}>-d</span> <span className={S}>{`'{`}</span>
        {"\n"}
        {"    "}
        <span className={K}>&quot;contents&quot;</span>
        <span className={P}>: [</span>
        {"\n"}
        {"      "}
        <span className={P}>{`{ `}</span>
        <span className={K}>&quot;role&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;user&quot;</span>
        <span className={P}>,</span>
        {"\n"}
        {"        "}
        <span className={K}>&quot;parts&quot;</span>
        <span className={P}>: [{"{"}</span> <span className={K}>&quot;text&quot;</span>
        <span className={P}>: </span>
        <span className={S}>&quot;...&quot;</span>
        <span className={P}>{` }] }`}</span>
        {"\n"}
        {"    "}
        <span className={P}>]</span>
        {"\n"}
        {"  "}
        <span className={S}>{`}'`}</span>
      </code>
    </pre>
  );
}

function ResGemini({ a }: { a: TabAccent }) {
  return (
    <pre className="mt-4 whitespace-pre-wrap break-words font-mono text-[12.5px] leading-[1.75] text-[#c6d0dc]">
      <code>
        {"{\n"}
        {"  "}
        <span className={K}>&quot;candidates&quot;</span>
        <span className={P}>: [{"{"}</span> <span className={K}>&quot;content&quot;</span>
        <span className={P}>: {"{"}</span> <span className={K}>&quot;parts&quot;</span>
        <span className={P}>: [{"{"}</span> <span className={K}>&quot;text&quot;</span>
        <span className={P}>: </span>
        <span className={a.responseText}>&quot;Gemini request served.&quot;</span>
        <span className={P}>{` }] } }]`}</span>
        {"\n"}
        {"  "}
        <span className={K}>&quot;usageMetadata&quot;</span>
        <span className={P}>: {"{"}</span> <span className={K}>&quot;totalTokenCount&quot;</span>
        <span className={P}>: </span>
        <span className={a.number}>25</span>
        <span className={P}>{` }`}</span>
        {"\n"}
        {"}"}
      </code>
    </pre>
  );
}

const tabs: ApiTab[] = [
  {
    id: "chat",
    label: "Chat",
    path: "/v1/chat/completions",
    stats: { ms: "142 MS", tokens: "27 TOKENS", cost: "COST $0.00081" },
    accent: accents.chat,
    request: <ReqChat />,
    response: <ResChat a={accents.chat} />,
  },
  {
    id: "responses",
    label: "Responses",
    path: "/v1/responses",
    stats: { ms: "168 MS", tokens: "31 TOKENS", cost: "COST $0.00093" },
    accent: accents.responses,
    request: <ReqResponses />,
    response: <ResResponses a={accents.responses} />,
  },
  {
    id: "claude",
    label: "Claude",
    path: "/v1/messages",
    stats: { ms: "156 MS", tokens: "29 TOKENS", cost: "COST $0.00087" },
    accent: accents.claude,
    request: <ReqClaude />,
    response: <ResClaude a={accents.claude} />,
  },
  {
    id: "gemini",
    label: "Gemini",
    path: "/v1beta/models/{model}:generateContent",
    stats: { ms: "93 MS", tokens: "25 TOKENS", cost: "COST $0.00075" },
    accent: accents.gemini,
    request: <ReqGemini />,
    response: <ResGemini a={accents.gemini} />,
  },
];

export function ApiPlayground() {
  const [activeId, setActiveId] = useState("chat");
  const tab = tabs.find((t) => t.id === activeId) ?? tabs[0];
  const a = tab.accent;

  return (
    <div className="relative overflow-hidden rounded-[24px] bg-[#0c1118] shadow-2xl ring-1 ring-white/5">
      {/* Tabs + status */}
      <div className="flex items-center justify-between border-b border-white/8 px-6 pt-5">
        <div className="flex items-center gap-7">
          {tabs.map((t) => {
            const active = t.id === activeId;
            return (
              <button
                key={t.id}
                type="button"
                onClick={() => setActiveId(t.id)}
                className={`relative pb-3 text-[13px] transition-colors ${
                  active ? `font-medium ${t.accent.active}` : "text-[#6b7688] hover:text-[#9aa6b8]"
                }`}
              >
                {t.label}
                {active ? (
                  <span className={`absolute inset-x-0 -bottom-px h-0.5 rounded-full ${t.accent.underline}`} />
                ) : null}
              </button>
            );
          })}
        </div>
        <div className="flex items-center gap-2 pb-3 text-[12px] text-[#8b97a8]">
          <span className="size-1.5 rounded-full bg-success" />
          200 OK
        </div>
      </div>

      {/* Endpoint */}
      <div className="flex items-start gap-3 border-b border-white/8 px-6 py-4">
        <span className={`mt-0.5 shrink-0 rounded-md px-2 py-1 text-[11px] font-semibold tracking-wide ${a.post}`}>POST</span>
        <code className="min-w-0 flex-1 break-all font-mono text-[14px] text-[#c6d0dc]">{tab.path}</code>
      </div>

      {/* Request */}
      <div className="min-h-[220px] border-b border-white/8 px-6 py-6">
        <p className="text-[11px] font-semibold tracking-[0.12em] text-[#5c6a7e]">REQUEST</p>
        {tab.request}
      </div>

      {/* Response */}
      <div className="min-h-[160px] border-b border-white/8 px-6 py-6">
        <p className="text-[11px] font-semibold tracking-[0.12em] text-[#5c6a7e]">RESPONSE</p>
        {tab.response}
      </div>

      {/* Stats footer */}
      <div className="flex items-center justify-between px-6 py-4 text-[11px] tracking-wide text-[#5c6a7e]">
        <div className="flex items-center gap-4">
          <span>{tab.stats.ms}</span>
          <span className="text-[#2a3442]">•</span>
          <span>{tab.stats.tokens}</span>
          <span className="text-[#2a3442]">•</span>
          <span>{tab.stats.cost}</span>
        </div>
        <div className="flex items-center gap-3">
          <span>STREAM</span>
          <span className="text-[#2a3442]">•</span>
          <span>SSE</span>
        </div>
      </div>
    </div>
  );
}
