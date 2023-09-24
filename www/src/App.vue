<script setup lang="ts">
import { Suspense } from 'vue'
import NavBar from './components/NavBar.vue'
</script>

<template>
  <header>
    <NavBar />
  </header>
  <Suspense>
    <template #default>
      <router-view></router-view>
    </template>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<style scoped></style>
