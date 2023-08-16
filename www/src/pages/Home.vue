<script setup lang="ts">
import { defineComponent, Suspense } from 'vue'
import ServerTable from '../components/ServerTable.vue'
</script>

<template>
    <section class="section">
        <div class="container">
            <Suspense fallback="Loading...">
                <ServerTable msg="toto" />
            </Suspense>
        </div>
    </section>
</template>