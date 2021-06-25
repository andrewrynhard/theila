<!--
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
-->
<template>
  <div class="flex flex-col">
    <div class="flex-0 px-3 py-2 mb-2">
      <t-breadcrumbs>{{ $route.params.node }} {{ $route.params.service }} Logs</t-breadcrumbs>
    </div>
    <t-alert v-if="err" title="Failed to Fetch Logs" type="error">
      {{ err }}.
    </t-alert>
    <div v-else class="flex-1 pb-3 logs">
      <div class="flex border-b border-talos-gray-300 dark:border-talos-gray-600 p-4 gap-1">
        <div class="flex-1">{{ logs.length }} lines</div>
        <div class="flex items-center justify-center gap-2 text-talos-gray-800 hover:text-talos-gray-600 dark:text-talos-gray-400 dark:hover:text-talos-gray-300">
          <Switch
            v-model="follow"
            class="inline-flex justify-center items-center w-5 h-5 rounded-md border-2 border-talos-gray-800 dark:border-talos-gray-400 outline-none"
            >
            <check-icon v-if="follow" class="w-4 h-4 inline-block"/>
          </Switch>
          <div @click="() => { follow = !follow }" class="cursor-pointer uppercase text-sm select-none font-bold">Follow Logs</div>
        </div>
      </div>
      <div class="flex-1 flex flex-col overflow-auto w-full h-full text-xs" ref="logView" style="font-family: monospace">
        <div v-for="line, index in logs" :key="index" class="log-line">
          <div class="inline-block mr-2 text-talos-gray-500 font-bold select-none">{{ index + 1 }}</div>
          <p>{{ line }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script type="ts">
import { useRoute } from 'vue-router';
import TBreadcrumbs from '../../components/TBreadcrumbs.vue';
import TAlert from '../../components/TAlert.vue';
import { ref, watch, onUnmounted, onUpdated, computed } from 'vue';
import { MachineService, Options, subscribe } from '../../api/grpc';
import { Source } from '../../common/theila.pb';
import { Switch } from '@headlessui/vue';
import { CheckIcon } from '@heroicons/vue/solid';

export default {
  components: {
    CheckIcon,
    Switch,
    TBreadcrumbs,
    TAlert,
  },

  setup() {
    const logs = ref([]);
    const route = useRoute();
    const follow = ref(true);
    const logView = ref(null);

    let line = 0;
    let stream = ref(null);

    const reset = () => {
      line = 0;
      logs.value = [""];
    };

    const scrollToBottom = () => {
      if(!logView.value)
        return;

      logView.value.scrollTop = logView.value.scrollHeight;
    }

    onUpdated(() => {
      if(follow.value)
        scrollToBottom();
    });

    const init = () => {
      if(stream.value)
        stream.value.shutdown();

      reset();

      stream.value = subscribe(MachineService.Logs, {
        namespace: "system",
        id: route.params.service,
        follow: true,
        tailLines: 250,
      }, (resp) => {
        if(resp.error) {
          reset();
          return;
        }

        const chunk = atob(resp.bytes);

        for(const s of chunk) {
          if(s === "\n") {
            logs.value.push("");
            line++;

            continue;
          }

          logs.value[line] += s;
        }

        scrollToBottom();
      }, new Options(Source.Talos, {
        nodes: [route.params.node],
      }));
    }

    init();

    watch(() => route.params.service, () => {
      init();
    });

    watch(follow, (val) => {
      if(val)
        scrollToBottom();
    });

    onUnmounted(() => {
      if(stream.value)
        stream.value.shutdown();
    });

    return {
      logView,
      follow,
      err: computed(() => {
        return stream.value ? stream.value.err : null;
      }),
      logs: computed(() => {
        if(logs.value[logs.value.length - 1] == "") {
          return logs.value.slice(0, logs.value.length - 2);
        }

        return logs.value;
      }),
    }
  }
}
</script>

<style scoped>
.logs {
  @apply flex flex-col w-full h-full overflow-hidden bg-white text-talos-gray-100 border rounded-md border-talos-gray-300 dark:border-talos-gray-600 dark:bg-talos-gray-900 text-talos-gray-900 dark:text-talos-gray-100;
}

.log-line {
  @apply hover:bg-talos-gray-100 dark:hover:bg-talos-gray-700 whitespace-pre-line flex px-3;
}
</style>