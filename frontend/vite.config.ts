import { fileURLToPath, URL } from 'node:url'
import path from 'node:path'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import imagemin from 'vite-plugin-imagemin'
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
      imagemin({
        gifsicle: {
          optimizationLevel: 3,
          interlaced: false,
        },
        optipng: {
          optimizationLevel: 3,
        },
        mozjpeg: {
          quality: 75,
        },
        pngquant: {
          quality: [0.65, 0.8],
          speed: 4,
        },
        svgo: {
          plugins: [
            {
              name: 'removeViewBox',
            },
            {
              name: 'removeEmptyAttrs',
              active: false,
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
