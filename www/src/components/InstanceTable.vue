<script setup lang="ts">
import { ref, getCurrentInstance, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ApiError, Instance } from '../api/generated'
import Error from './Error.vue'
import Pagination from './Pagination.vue'

const props = defineProps<{ page: number }>();

type Hoverable<T> = T & { hovered: boolean };

const { $api } = getCurrentInstance()!.appContext.config.globalProperties;

// Select based on screen height
const perPage = Math.floor(window.innerHeight / 60);

const totalInstances = ref(0);
const instances = ref<Hoverable<Instance>[]>([]);
const error = ref<Error | null>(null);

const loadPage = async () => {
  try {
    const resp = await $api.listInstances(undefined, props.page, perPage);
    if ("results" in resp) {
      instances.value = resp.results.map(s => ({ ...s, hovered: false }));
      totalInstances.value = resp.total;
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

const router = useRouter();
// it's necessary to use programmatic navigation because the router-link
// will mess up the table row css
let goToInstance = (instance: Instance) => {
  router.push({ path: "/instances/" + instance.id });
};
</script>

<template>
  <Error :error="error" />
  <Pagination :page="props.page" :pageSize="perPage" :total="totalInstances" />
  <table class="table" style="width: 100%">
    <thead>
      <tr>
        <th>Domain</th>
        <th>Status</th>
        <th>Software</th>
        <th>Total Users</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="   instance    in    instances   " :key="instance.id" @click="goToInstance(instance)"
        @mouseover="instance.hovered = true" @mouseleave="instance.hovered = false"
        :class="{ 'is-hovered': instance.hovered }">
        <td>
          <a :href='"https://" + instance.domain'>{{ instance.domain }}</a>
        </td>
        <td>{{ instance.status }}</td>
        <td>{{ instance.software }}</td>
        <td>{{ instance.total_users }}</td>
      </tr>
    </tbody>
  </table>
</template>

<style scoped>
tr.is-hovered {
  cursor: pointer;
  box-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
  transition: background-color 0.3s, box-shadow 0.3s, opacity 0.3s, color 0.3s;
}
</style>
