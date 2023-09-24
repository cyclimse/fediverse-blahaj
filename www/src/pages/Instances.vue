<script setup lang="ts">
import InstanceTable from '../components/InstanceTable.vue'

const props = defineProps<{ page: number }>();
</script>

<template>
    <section class="section">
        <div class="container">
            <InstanceTable :page="props.page" />
        </div>
    </section>
</template>
