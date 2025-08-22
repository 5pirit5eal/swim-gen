import pluginVitest from '@vitest/eslint-plugin'
import tseslint from 'typescript-eslint'

export default tseslint.config(
    {
        ignores: ['**/dist/**', '**/coverage/**'],
    },
    ...tseslint.configs.recommended,
    {
        files: ['src/**/__tests__/*'],
        plugins: {
            vitest: pluginVitest,
        },
        rules: {
            ...pluginVitest.configs.recommended.rules,
        },
        languageOptions: {
            globals: {
                ...pluginVitest.environments.env.globals,
            },
        },
    },
)
