// src/__tests__/setup.ts
import dotenv from 'dotenv';
import { beforeAll } from 'vitest';

// Load environment variables from .env.test file for the test environment
beforeAll(() => {
    dotenv.config({ path: '.env.test', override: true });
})
