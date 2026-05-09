import { fileURLToPath, URL } from 'node:url'
import path from 'node:path'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { ViteImageOptimizer } from 'vite-plugin-image-optimizer'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [
      vue(),
      mode === 'development' && vueDevTools(),
      VueI18nPlugin({
        include: [path.resolve(__dirname, './src/locales/**')]
      }),
      ViteImageOptimizer({
        png: {
          quality: 80,
        },
        jpeg: {
          quality: 75,
        },
        svg: {
          plugins: [
            {
              name: 'removeViewBox',
            },
          ],
        },
      }),
    ].filter(Boolean),
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      proxy: {
        '/api': {
          target: env.BFF_APP_API_URL || 'http://localhost:8080',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
        },
      },
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes('node_modules')) {
              // Extract the package name from the path
              const parts = id.split('node_modules/')[1].split('/');
              // Handle scoped packages
              const pkgName = parts[0].startsWith('@') ? parts.slice(0, 2).join('/') : parts[0];
              return `vendor-${pkgName}`;
            }
          },
        },
      },
    },
  }
})
