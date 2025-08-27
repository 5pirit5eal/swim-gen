import express from "express";
import dotenv from "dotenv";
import axios from "axios";
import cors from "cors";
import rateLimit from "express-rate-limit";
import { MemoryStore } from "./memory-store";
import * as authModule from "./auth";

dotenv.config();

const app = express();
app.set("trust proxy", 1);
const port = process.env.PORT || 8080;

// Middleware to handle JSON bodies
app.use(express.json());

// Enable CORS
const corsOptions = {
  origin: process.env.FRONTEND_URL,
  optionsSuccessStatus: 200, // some legacy browsers (IE11, various SmartTVs) choke on 204
};
app.use(cors(corsOptions));

// Rate limiting
export const store = new MemoryStore();
export const apiLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100, // Limit each IP to 100 requests per windowMs
  standardHeaders: true,
  legacyHeaders: true,
  store,
});

// Apply the rate limiting middleware to API calls only
app.use("/api", apiLimiter);

// Generic proxy handler for all API requests
async function proxyRequest(req: express.Request, res: express.Response) {
  const { method, originalUrl, body } = req;

  // The frontend will call paths like /api/query. We strip /api before forwarding.
  const backendPath = originalUrl.replace(/^\/api/, "");
  const targetUrl = `${process.env.BACKEND_URL}${backendPath}`;

  console.log(`Proxying request: ${method} ${originalUrl} -> ${targetUrl}`);

  try {
    const authHeaders = await authModule.getAuthHeaders();

    const response = await axios({
      method,
      url: targetUrl,
      data: body,
      headers: {
        ...authHeaders,
        "Content-Type": "application/json",
      },
    });

    res.status(response.status).json(response.data);
  } catch (error) {
    console.error(`Error proxying request to ${targetUrl}:`, error);
    if (error && typeof error === "object" && "response" in error && error.response) {
      const axiosError = error as { response: { status: number; data: unknown } };
      res.status(axiosError.response.status).json(axiosError.response.data);
    } else {
      res.status(500).json({ message: "Error proxying request to backend" });
    }
  }
}

app.get("/health", (req, res) => {
  res.status(200).send("OK");
});

// Block the /api/scrape endpoint
app.use("/api/scrape", (req, res) => {
  res.status(403).send("This endpoint is not available.");
});

// All API routes from the frontend are prefixed with /api
app.use("/api", proxyRequest);

// Avoid binding a port during tests
if (process.env.NODE_ENV !== "test") {
  app.listen(port, () => {
    console.log(`BFF server listening on port ${port}`);
  });
}

// Exported for testing/mocking in unit tests
export { app, store as testingStore };
