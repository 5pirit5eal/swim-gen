<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore()
</script>

<template>
    <div class="profile-view">
        <h1>Profile</h1>
        <p v-if="auth.user">Welcome, {{ auth.user.email }}</p>
    </div>
</template>

<style scoped>
.profile-view {
    padding: 2rem;
}
</style>
