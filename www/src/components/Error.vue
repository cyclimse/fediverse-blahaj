<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{ error: Error | null }>();

const shown = ref(true);
watch(props, () => shown.value = true);
</script>

<template>
    <div class="notification is-danger" v-if="error && shown">
        <button class="delete" @click="shown = false"></button>
        <strong>{{ error.name }}</strong> {{ error.message }}
    </div>
</template>
