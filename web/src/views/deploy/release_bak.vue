<template>
  <PageWrapper contentFullHeight>
    <a-card title="新建上线单" v-loading="loadingRef">
      <alert v-if="alertMsg.msg" :type="alertMsg.type">
        <template  #description>
          {{ alertMsg.msg }}
          <br/>
          <a-button type="info" v-if="alertMsg.state==DeployStatus.Audit" @click="startRelease">
            开始上线
          </a-button>
        </template>
      </alert>

      <a-tabs
        :style="{ minHeight: '400px' }"
        v-model:activeKey="activeKey"
      >
        <a-tab-pane class="text-white bg-gray-700" v-for="server in servers" :key="server.id"  :tab="server.host">
          <pre v-for="record in server.records" :key="record.id">
[{{server.user}}@{{server.host}}]$ {{record.command}}
{{ record.output }}
          </pre>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </PageWrapper>
</template>
<script lang="ts" setup>
import {onMounted, reactive, ref, watchEffect} from 'vue';
import {PageWrapper} from '/@/components/Page';
import {Alert, Card as ACard,TabPane as ATabPane, Tabs as ATabs} from 'ant-design-vue';
import {ListItem as DeployListItem, RecordItem} from "/@/api/deploy/model";
import {useRoute} from "vue-router";
import {DeployStatus, DeployStatusShowMsg} from "/@/enums/fieldEnum";
import {useWebSocket} from "@vueuse/core";
import {getDeployConsoleWs, detailDeploy, startDeploy} from "/@/api/deploy";
import {useUserStore} from "/@/store/modules/user";

type server = {
  id: number;
  user: string;
  host: string;
  port: number;
  name: string;
  records: { [key: number]: RecordItem };
}
type serverItems = { [key: number]: server }
type stateMsg = { type: "success" | "info" | "error" | "warning", msg: string, state: number }

const userStore = useUserStore();
const route = useRoute();
const loadingRef = ref(false);
const activeKey = ref<number>(0)
const errMsg = ref<string>("")
const deployDetail = ref<DeployListItem>()
const servers = ref<serverItems>({0:{id: 0, user: "local", host: "127.0.0.1", port: 0, name: "localhost", records: {}}})
const deployId = parseInt(route.params?.id as unknown as string)
const alertMsg = reactive<stateMsg>({type: "error", msg: "", state: 0})

const {status, open: wsOpen} = useWebSocket(getDeployConsoleWs(deployId), {
  autoReconnect: false,
  heartbeat: true,
  immediate: false,
  protocols: [userStore.getToken, userStore.getCurrentSpaceId.toString()],
  onError: (_, event) => {
    alertMsg.type = "error"
    alertMsg.msg = event.toString()
  },
  onMessage: (_, event) => {
    let data = JSON.parse(event.data)
    if (data.type == 'records') {
      data.records.forEach(v => {
        servers.value[v.server_id].records[v.id] = v
      })
    }
    if (data.type == 'append') {
      //console.log(data)
      data.records.forEach(v => {
        if (servers.value[v.server_id].records[v.id]) {
          servers.value[v.server_id].records[v.id].success += v.success
        } else {
          servers.value[v.server_id].records[v.id] = v
        }
      })
      //console.log(servers.value)
    }
  }
})

async function init(id: number) {
  try {
    loadingRef.value = true
    deployDetail.value = await detailDeploy(id, true)
  } catch (e) {
    alertMsg.msg = (e as unknown as Error).toString()
    return
  } finally {
    loadingRef.value = false
  }
  if (deployDetail.value.servers.length == 0) {
    alertMsg.type = "error"
    alertMsg.msg = "该上线单发布服务器为空，请检查在执行操作！"
    return
  }
  deployDetail.value.servers.forEach((v) => {
    servers.value[v.id] = {
      id: v.id,
      user: v.user,
      host: v.host,
      port: v.port,
      name: v.name,
      records: {},
    }
  })

  let sm = DeployStatusShowMsg(deployDetail?.value.status)
  alertMsg.msg = sm[0]
  alertMsg.type = sm[1]
  alertMsg.state = deployDetail.value.status
  //打开websocket
  if (deployDetail.value.status != DeployStatus.Waiting && deployDetail.value.status != DeployStatus.Audit) {
    wsOpen()
  }
  //wsOpen()
}

function startRelease() {
  loadingRef.value = true
  startDeploy(deployId, true).then(()=>{
    loadingRef.value = false
    if (status.value != "OPEN") {
      wsOpen()
    }
  }).catch((e) => {
    alertMsg.msg = e.toString()
    loadingRef.value = false
  })
}

watchEffect(() => {
  if (status.value == "CONNECTING") {
    loadingRef.value = true
  } else {
    loadingRef.value = false
  }
})

onMounted(() => {
  init(deployId)
})
</script>

