import { Store, Options, IncrementResponse } from "express-rate-limit";

export class MemoryStore implements Store {
  windowMs!: number;
  hits: { [key: string]: number[] } = {};

  init(options: Options): void {
    this.windowMs = options.windowMs;
  }

  async increment(key: string): Promise<IncrementResponse> {
    const now = Date.now();
    if (!this.hits[key]) {
      this.hits[key] = [];
    }

    // clean up old hits
    this.hits[key] = this.hits[key].filter((timestamp) => timestamp > now - this.windowMs);

    // add new hit
    this.hits[key].push(now);

    const totalHits = this.hits[key].length;
    const resetTime = new Date(now + this.windowMs);

    return {
      totalHits,
      resetTime,
    };
  }

  async decrement(key: string): Promise<void> {
    if (this.hits[key]) {
      this.hits[key].pop();
    }
  }

  async resetKey(key: string): Promise<void> {
    delete this.hits[key];
  }

  async resetAll(): Promise<void> {
    this.hits = {};
  }
}
