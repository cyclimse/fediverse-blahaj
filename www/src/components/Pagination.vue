<script setup lang="ts">
import { computed, defineProps } from 'vue'

const props = defineProps<{ page: number, pageSize: number, total: number }>();

const firstPage = 1;
const lastPage = computed(() => Math.ceil(props.total / props.pageSize));

const paginationButtons = computed(() => {
    var buttons = [firstPage, firstPage + 1, props.page - 1, props.page, props.page + 1, lastPage.value - 1, lastPage.value];
    buttons = buttons.filter(b => b >= firstPage && b <= lastPage.value);
    const uniqueButtons = new Set(buttons);
    buttons = Array.from(uniqueButtons);
    return buttons.map((b, i) => {
        if (i === 0) return { page: b, ellipsis: false };
        if (i === buttons.length - 1) return { page: b, ellipsis: false };
        if (buttons[i + 1] - buttons[i] > 1) return { page: b, ellipsis: true };
        return { page: b, ellipsis: false };
    });
});
</script>

<template>
    <nav class="pagination" role="navigation" aria-label="pagination">
        <router-link class="pagination-previous" :class="{ 'is-disabled': page <= firstPage }"
            :to="{ query: { page: page > firstPage ? page - 1 : page } }">Previous</router-link>
        <router-link class="pagination-next" :class="{ 'is-disabled': page >= lastPage }"
            :to="{ query: { page: page < lastPage ? page + 1 : page } }">Next</router-link>
        <ul class="pagination-list">
            <li v-for="button in paginationButtons" :key="button.page">
                <router-link class="pagination-link" :class="{ 'is-current': button.page === page }"
                    :to="{ query: { page: button.page } }">{{ button.page
                    }}</router-link>
                <span class="pagination-ellipsis" v-if="button.ellipsis">…</span>
            </li>
        </ul>
    </nav>
</template>
