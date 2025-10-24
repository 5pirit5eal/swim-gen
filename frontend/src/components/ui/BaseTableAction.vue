<template>
    <div class="action-container">
        <!-- The content passed to this slot will be the action buttons -->
        <slot name="actions"></slot>
    </div>
</template>

<style scoped>
.action-container {
    /* Position absolutely relative to the nearest positioned ancestor */
    position: absolute;
    top: 0;
    bottom: 0;

    /* Place it just to the left of the parent cell's boundary */
    right: 100%;

    /* Vertically center the content inside */
    display: flex;
    align-items: center;

    /* Add some spacing for the buttons */
    padding: 0 0.75rem;

    /* Hide by default, will be shown on parent hover */
    opacity: 0;
    transition: opacity 0.2s ease-in-out;
    z-index: 10;

    /* Use a CSS variable for the background, defaulting to transparent */
    background-color: var(--action-bg-color, transparent);
    border-top-left-radius: 0.5rem;
    border-bottom-left-radius: 0.5rem;
    border: 1px solid var(--color-border);
    border-right: none;
}
</style>
