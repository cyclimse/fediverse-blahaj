<script setup lang="ts">
import { ref, getCurrentInstance, watch, defineProps } from 'vue'
import { ApiError, Crawl } from '../api/generated'
import Error from './Error.vue'
import Pagination from './Pagination.vue';

const props = defineProps<{ page: number, instanceId: string }>();

const { $api } = getCurrentInstance()!.appContext.config.globalProperties;

// Select based on screen height
const perPage = 10;

const totalCrawls = ref(0);
const crawls = ref<Crawl[]>([]);
const error = ref<Error | null>(null);

const loadPage = async () => {
    try {
        const resp = await $api.listCrawlsForInstance(props.instanceId, props.page, perPage);
        if ("results" in resp) {
            crawls.value = resp.results;
            totalCrawls.value = resp.total;
        } else {
            console.error(resp);
            error.value = { name: resp.code.toString(), message: resp.message };
        }
    } catch (e) {
        if (e instanceof ApiError) error.value = { name: "Error", message: e.message };
        else throw e;
    }
};

loadPage();
watch(props, loadPage);
</script>

<template>
    <Error :error="error" />
    <Pagination :page="props.page" :pageSize="perPage" :total="totalCrawls" />
    <table class="table" style="width: 100%">
        <thead>
            <tr>

            </tr>
        </thead>
        <tbody>
            <tr v-for="crawl in crawls" :key="crawl.id">
                {{ crawl }}
            </tr>
        </tbody>
    </table>
</template>
