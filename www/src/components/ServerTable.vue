<script setup lang="ts">
import { ref, watch, getCurrentInstance } from 'vue'
import { useRouter } from 'vue-router'
import { ApiError, Server } from '../api/generated'

defineProps<{ msg: string }>()

type Hoverable<T> = T & { hovered: boolean };

const { $api } = getCurrentInstance()!.appContext.config.globalProperties;
const page = ref(1);
const error = ref<Error | null>(null);
const servers = ref<Hoverable<Server>[]>([]);

let loadPage = async () => {
  try {
    const resp = await $api.listServers(undefined, page.value);
    if ("results" in resp) {
      servers.value = resp.results.map(s => ({ ...s, hovered: false }));
    } else {
      console.error(resp);
      error.value = { name: resp.code.toString(), message: resp.message };
    }
  } catch (e) {
    if (e instanceof ApiError) error.value = { name: "Error", message: e.message };
    else throw e;
  }
}

const router = useRouter();

let goToServer = (server: Server) => {
  router.push({ path: "/servers/" + server.id });
}

watch(page, loadPage, { immediate: true });

loadPage();
</script>

<template>
  <div class="notification is-danger" v-if="error">
    <button class="delete" @click="error = null"></button>
    <strong>{{ error.name }}</strong> {{ error.message }}
  </div>
  <nav class="pagination" role="navigation" aria-label="pagination">
    <a class="pagination-previous" :class="{ 'is-disabled': page <= 1 }" @click="page--"
      title="This is the first page">Previous</a>
    <a class="pagination-next" @click="page++" title="This is the last page">Next page</a>
    <ul class="pagination-list">
      <li>
        <a class="pagination-link is-current" aria-label="Page 1" aria-current="page">1</a>
      </li>
      <li>
        <a class="pagination-link" aria-label="Goto page 2">2</a>
      </li>
      <li>
        <a class="pagination-link" aria-label="Goto page 3">3</a>
      </li>
    </ul>
  </nav>
  <table class="table">
    <thead>
      <tr>
        <th>Domain</th>
        <th>Status</th>
        <th>Software</th>
        <th>Total Users</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="server in servers" :key="server.id" @click="goToServer(server)" @mouseover="server.hovered = true"
        @mouseleave="server.hovered = false" :class="{ 'is-hovered': server.hovered }">
        <td>
          <a :href='"https://" + server.domain'>{{ server.domain }}</a>
        </td>
        <td>{{ server.status }}</td>
        <td>{{ server.software }}</td>
        <td>{{ server.total_users }}</td>
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
