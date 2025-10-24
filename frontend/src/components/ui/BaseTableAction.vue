<script setup lang="ts">
defineProps<{
    isFirst: boolean
    isLast: boolean
}>()

const emit = defineEmits<{
    (e: 'add'): void
    (e: 'remove'): void
    (e: 'move-up'): void
    (e: 'move-down'): void
}>()
</script>

<template>
    <div class="action-container">
        <div class="action-grid">
            <button class="action-btn move-up" :disabled="isFirst" @click.stop="emit('move-up')"
                title="Move row up"></button>
            <button class="action-btn add" @click.stop="emit('add')" title="Add new row"></button>
            <button class="action-btn move-down" :disabled="isLast" @click.stop="emit('move-down')"
                title="Move row down"></button>
            <button class="action-btn remove" @click.stop="emit('remove')" title="Remove row"></button>
        </div>
    </div>
</template>

<style scoped>
.action-container {
    position: absolute;
    top: 0;
    bottom: 0;
    right: 101%;
    display: flex;
    align-items: center;
    padding: 0 0.5rem;
    opacity: 0;
    transition:
        opacity 0.2s ease-in-out,
        transform 0.2s ease-in-out;
    z-index: 10;
    background-color: var(--color-primary);
    border-top-left-radius: 0.5rem;
    border-bottom-left-radius: 0.5rem;
    outline: solid 0.0625rem var(--color-primary);
    transform: translateX(0.3125rem);
}

.action-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(2, 1fr);
    gap: 0.3rem;
}

.action-btn {
    width: 1.5rem;
    height: 1.5rem;
    border: none;
    background-color: transparent;
    cursor: pointer;
    position: relative;
    border-radius: 0.25rem;
    transition: var(--color-primary-hover) 0.2s;
}

.action-btn:hover:not(:disabled) {
    background-color: var(--color-primary-hover);
}

.action-btn:disabled {
    cursor: not-allowed;
    opacity: 0.3;
}

/* --- Icons using Pseudo-elements --- */

/* Plus Icon */
.add::before,
.add::after {
    content: '';
    position: absolute;
    background-color: white;
    transition: background-color 0.2s;
}

.add::before {
    top: 50%;
    left: 20%;
    width: 60%;
    height: 0.125rem;
    transform: translateY(-50%);
}

.add::after {
    top: 20%;
    left: 50%;
    width: 0.125rem;
    height: 60%;
    transform: translateX(-50%);
}

/* Minus Icon */
.remove::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 20%;
    width: 60%;
    height: 0.125rem;
    transform: translateY(-50%);
    background-color: white;
    transition: background-color 0.2s;
}

/* Arrow Up Icon */
.move-up::before {
    content: '';
    position: absolute;
    top: 40%;
    left: 50%;
    width: 0;
    height: 0;
    border-left: 0.375rem solid transparent;
    border-right: 0.375rem solid transparent;
    border-bottom: 0.5rem solid white;
    transform: translateX(-50%);
    transition: border-color 0.2s;
}

/* Arrow Down Icon */
.move-down::before {
    content: '';
    position: absolute;
    top: 60%;
    left: 50%;
    width: 0;
    height: 0;
    border-left: 0.375rem solid transparent;
    border-right: 0.375rem solid transparent;
    border-top: 0.5rem solid white;
    transform: translate(-50%, -50%);
    transition: border-color 0.2s;
}

@media (max-width: 740px) {
    .action-container {
        padding: 0 0.25rem;
        transform: translateX(0.125rem);
    }

    .action-grid {
        gap: 0.15rem;
    }

    .action-btn {
        width: 1.2rem;
        height: 1.2rem;
    }

    .move-up::before {
        border-left-width: 0.25rem;
        border-right-width: 0.25rem;
        border-bottom-width: 0.375rem;
    }

    .move-down::before {
        border-left-width: 0.25rem;
        border-right-width: 0.25rem;
        border-top-width: 0.375rem;
    }
}
</style>
