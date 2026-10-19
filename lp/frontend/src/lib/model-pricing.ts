/** One row of the token plan pricing table. "n.a." means the model has no such cache tier. */
export type ModelPricing = {
  model: string;
  cachedInput: number;
  input: number;
  output: number;
  cacheWrite5m: number | null;
  cacheWrite1h: number | null;
};

/** Prices are credits per token. */
export const MODEL_PRICING: ModelPricing[] = [
  { model: "DeepSeek-V4.1-Flash", cachedInput: 2, input: 100, output: 420, cacheWrite5m: null, cacheWrite1h: null },
  { model: "DeepSeek Pro", cachedInput: 15, input: 450, output: 1400, cacheWrite5m: null, cacheWrite1h: null },
  { model: "DeepSeek-V4.1-Flash-Fast", cachedInput: 33, input: 216, output: 858, cacheWrite5m: null, cacheWrite1h: null },
  { model: "GLM-5.3-Flash", cachedInput: 20, input: 100, output: 350, cacheWrite5m: null, cacheWrite1h: null },
  { model: "GLM-5.3", cachedInput: 180, input: 1000, output: 3100, cacheWrite5m: null, cacheWrite1h: null },
  { model: "Kimi K3", cachedInput: 200, input: 2100, output: 10500, cacheWrite5m: 2100, cacheWrite1h: 4200 },
  { model: "MiniMax M3 Standard", cachedInput: 42, input: 200, output: 850, cacheWrite5m: null, cacheWrite1h: null },
  { model: "MiMo v2.6 Pro", cachedInput: 2.5, input: 300, output: 600, cacheWrite5m: null, cacheWrite1h: null },
  { model: "MiMo v2.6 Flash", cachedInput: 2, input: 100, output: 200, cacheWrite5m: null, cacheWrite1h: null },
];