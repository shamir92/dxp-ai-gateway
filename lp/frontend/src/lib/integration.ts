/**
 * Integration snippets for the token plan.
 *
 * Both protocol styles are served from the SAME base URL, `<base>/v1`:
 *   POST <base>/v1/chat/completions   OpenAI-compatible
 *   POST <base>/v1/messages           Anthropic-compatible
 * (Verified against the live gateway: both return 401 unauthenticated, i.e. the
 * paths exist and only the key is missing. A separate /anthropic base is NOT served.)
 */

export type CodeSample = {
  id: string;
  label: string;
  language: string;
  code: string;
};

export type ToolGuide = {
  id: string;
  name: string;
  blurb: string;
  /** Each step is one line of instruction; `code` blocks are rendered verbatim. */
  steps: string[];
  config: string;
  configLabel: string;
};

export function baseURLFor(endpoint: string): string {
  return `${endpoint.replace(/\/+$/, "")}/v1`;
}

/** `sk-` is presentation only; the gateway accepts the value with or without it. */
export const KEY_PLACEHOLDER = "sk-YOUR_API_KEY";

export function codeSamples(baseURL: string): CodeSample[] {
  return [
    {
      id: "curl",
      label: "curl",
      language: "bash",
      code: `curl ${baseURL}/chat/completions \\
  -H "Authorization: Bearer ${KEY_PLACEHOLDER}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "deepseek-v4.1-flash",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": false
  }'`,
    },
    {
      id: "python",
      label: "Python",
      language: "python",
      code: `# pip install openai
from openai import OpenAI

client = OpenAI(
    base_url="${baseURL}",
    api_key="${KEY_PLACEHOLDER}",
)

response = client.chat.completions.create(
    model="deepseek-v4.1-flash",
    messages=[{"role": "user", "content": "Hello!"}],
)
print(response.choices[0].message.content)`,
    },
    {
      id: "go",
      label: "Go",
      language: "go",
      code: `// go get github.com/openai/openai-go
package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	client := openai.NewClient(
		option.WithBaseURL("${baseURL}"),
		option.WithAPIKey("${KEY_PLACEHOLDER}"),
	)

	chat, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:    openai.String("deepseek-v4.1-flash"),
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("Hello!")},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(chat.Choices[0].Message.Content)
}`,
    },
    {
      id: "node",
      label: "Node.js",
      language: "javascript",
      code: `// npm install openai
import OpenAI from "openai";

const client = new OpenAI({
  baseURL: "${baseURL}",
  apiKey: process.env.DXP_API_KEY ?? "${KEY_PLACEHOLDER}",
});

const response = await client.chat.completions.create({
  model: "deepseek-v4.1-flash",
  messages: [{ role: "user", content: "Hello!" }],
});
console.log(response.choices[0].message.content);`,
    },
    {
      id: "anthropic",
      label: "Anthropic style",
      language: "python",
      code: `# pip install anthropic
from anthropic import Anthropic

client = Anthropic(
    base_url="${baseURL}",
    api_key="${KEY_PLACEHOLDER}",
)

message = client.messages.create(
    model="deepseek-v4.1-flash",
    max_tokens=1024,
    messages=[{"role": "user", "content": "Hello!"}],
)
print(message.content[0].text)`,
    },
  ];
}

export function toolGuides(baseURL: string, apiKey: string): ToolGuide[] {
  return [
    {
      id: "mimo-code",
      name: "MiMo Code",
      blurb: "Point it at the gateway with an OpenAI-compatible provider.",
      steps: [
        "Install MiMo Code and open its settings.",
        "Add a custom OpenAI-compatible provider.",
        "Paste the base URL and your API key, then select a model.",
      ],
      configLabel: "~/.mimocode/config.json",
      config: `{
  "providers": {
    "dxp": {
      "type": "openai",
      "baseURL": "${baseURL}",
      "apiKey": "${apiKey}",
      "models": ["deepseek-v4.1-flash"]
    }
  }
}`,
    },
    {
      id: "hermes",
      name: "Hermes Agent",
      blurb: "Configure the OpenAI-compatible endpoint and default model.",
      steps: [
        "Open the Hermes Agent configuration.",
        "Set the API base to the gateway URL.",
        "Provide the API key and choose a model.",
      ],
      configLabel: "hermes.toml",
      config: `[provider]
type    = "openai"
base_url = "${baseURL}"
api_key  = "${apiKey}"
model    = "deepseek-v4.1-flash"`,
    },
    {
      id: "openclaw",
      name: "OpenClaw",
      blurb: "Uses the OpenAI-compatible interface.",
      steps: [
        "Run the OpenClaw setup wizard (openclaw config).",
        "Choose a custom OpenAI-compatible endpoint when prompted.",
        "Enter the base URL and API key, then pick a model.",
      ],
      configLabel: "~/.openclaw/config.yaml",
      config: `provider: openai
base_url: ${baseURL}
api_key: ${apiKey}
model: deepseek-v4.1-flash`,
    },
    {
      id: "claude-code",
      name: "Claude Code",
      blurb: "Uses the Anthropic-compatible interface on the same base URL.",
      steps: [
        "Export the two variables below in the shell that launches Claude Code.",
        "Anthropic style is served at <base>/v1/messages — the same base as OpenAI.",
        "Start Claude Code; it picks the variables up automatically.",
      ],
      configLabel: "shell profile",
      config: `export ANTHROPIC_BASE_URL="${baseURL}"
export ANTHROPIC_API_KEY="${apiKey}"`,
    },
    {
      id: "codex",
      name: "Codex",
      blurb: "OpenAI-compatible provider configured in the CLI profile.",
      steps: [
        "Add a model provider to the Codex profile.",
        "Set base_url to the gateway and paste the API key.",
        "Reference the provider from your chosen profile.",
      ],
      configLabel: "~/.codex/config.toml",
      config: `model = "deepseek-v4.1-flash"
model_provider = "dxp"

[model_providers.dxp]
name = "DXP AI Gateway"
base_url = "${baseURL}"
env_key = "DXP_API_KEY"

# then: export DXP_API_KEY="${apiKey}"`,
    },
  ];
}