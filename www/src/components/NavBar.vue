
<script setup lang="ts">
import { ref } from 'vue'

const burgerActive = ref(false);

const links = {
    "Home": "/",
    "Instances": "/instances",
    "About": "/about",
};
</script>


<template>
    <nav class="navbar is-primary" role="navigation" aria-label="main navigation">
        <div class="navbar-brand">
            <a class="navbar-item" href="/">
                <img src="/logo.png" width="28" height="28">
            </a>

            <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false" @click="burgerActive = !burgerActive">
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
            </a>
        </div>

        <div class="navbar-menu" :class="{ 'is-active': burgerActive }">
            <div class="navbar-start">
                <a v-for="(link, name) in links" :key="name" :href="link" class="navbar-item">
                    {{ name }}
                </a>
            </div>
        </div>
    </nav>
</template>
