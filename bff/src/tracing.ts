/**
 * OpenTelemetry SDK initialization for the BFF service.
 *
 * IMPORTANT: This file must be imported before any other modules in main.ts
 * so that auto-instrumentation can patch Express and HTTP before they are loaded.
 */

import { NodeSDK } from "@opentelemetry/sdk-node";
import { getNodeAutoInstrumentations } from "@opentelemetry/auto-instrumentations-node";
import { TraceExporter } from "@google-cloud/opentelemetry-cloud-trace-exporter";
import { resourceFromAttributes } from "@opentelemetry/resources";
import { ATTR_SERVICE_NAME } from "@opentelemetry/semantic-conventions";
import { ParentBasedSampler, TraceIdRatioBasedSampler } from "@opentelemetry/sdk-trace-node";

const serviceName = process.env.OTEL_SERVICE_NAME || "swim-gen-bff";

// Validate samplerArg: must be a finite number in [0, 1]; default to 1.0 otherwise.
const _rawSamplerArg = parseFloat(process.env.OTEL_TRACES_SAMPLER_ARG || "1.0");
const samplerArg = Number.isFinite(_rawSamplerArg)
  ? Math.min(1.0, Math.max(0.0, _rawSamplerArg))
  : 1.0;

const sampler =
  samplerArg < 1.0
    ? new ParentBasedSampler({ root: new TraceIdRatioBasedSampler(samplerArg) })
    : undefined;

const sdk = new NodeSDK({
  resource: resourceFromAttributes({
    [ATTR_SERVICE_NAME]: serviceName,
  }),
  traceExporter: new TraceExporter(),
  instrumentations: [
    getNodeAutoInstrumentations({
      // Only instrument HTTP and Express; disable noisy/unnecessary instrumentations
      "@opentelemetry/instrumentation-fs": { enabled: false },
      "@opentelemetry/instrumentation-dns": { enabled: false },
      "@opentelemetry/instrumentation-net": { enabled: false },
    }),
  ],
  ...(sampler && { sampler }),
});

// start() is async; await it and log any startup failure before the server handles requests.
sdk.start().catch((err: unknown) => {
  console.error("OTel SDK failed to start:", err);
});

// Graceful shutdown
process.on("SIGTERM", () => {
  sdk
    .shutdown()
    .then(() => console.log("OTel SDK shut down"))
    .catch((err: unknown) => console.error("OTel SDK shutdown error", err))
    .finally(() => process.exit(0));
});

export { sdk };
