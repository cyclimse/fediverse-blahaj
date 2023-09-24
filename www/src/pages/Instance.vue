<script setup lang="ts">
import { getCurrentInstance, ref } from 'vue'
import { ApiError, Instance } from '../api/generated'
import Error from '../components/Error.vue'
import CrawlsForInstance from '../components/CrawlsForInstance.vue'

const { $api } = getCurrentInstance()!.appContext.config.globalProperties;
const props = defineProps<{ id: string }>();

const instance = ref<Instance | null>(null);
const error = ref<Error | null>(null);

try {
    const resp = await $api.getInstanceById(props.id);
    if ("id" in resp) {
        instance.value = resp;
    } else {
        console.error(resp);
        error.value = { name: resp.code.toString(), message: resp.message };
    }
} catch (e) {
    if (e instanceof ApiError) error.value = { name: "Error", message: e.message };
}
</script>


<template>
    <section class="section">
        <div class="container">
            <Error :error="error" />
            <Suspense fallback="Loading...">
                <div v-if="instance">
                    <h1 class="title">{{ instance.domain }}</h1>
                    <p class="subtitle">{{ instance.description }}</p>

                    <p><strong>Created:</strong> {{ instance }}</p>

                    <h2 class="title is-4">Crawls</h2>
                    <CrawlsForInstance :instanceId="instance.id" :page=1 />
                </div>


            </Suspense>
        </div>
    </section>
</template>