<template>
  <div id="terminal_box" v-loading="loading" style="width: 100%; height: 100%">
    <div ref="terminalRef" style="width: 100%; height: 100%; padding: 1px 2px; background: black">
    </div>

    <div id="showTip">
      {{ serverHost }}
      <Icon icon="ant-design:logout-outlined" :size="18" style="cursor: pointer"
            @click="toggleConnect"/>
    </div>

  </div>
</template>
<script lang="ts">
export default {
  name: 'ServerTerminal',
};
</script>
<script lang="ts" setup>
import {ref, onMounted, watchEffect, onBeforeUnmount} from 'vue';
import {Terminal} from 'xterm';
import {FitAddon} from 'xterm-addon-fit';
import {useRoute} from 'vue-router';
import 'xterm/css/xterm.css';
import {Icon} from '/@/components/Icon/index';
import {useTabs} from '/@/hooks/web/useTabs';
import {useWebSocket} from '@vueuse/core';
import {useUserStore} from '/@/store/modules/user';
import {getServerSshWs} from '/@/api/server'

interface wsMsg {
  typ: "cmd" | "resize";
  cmd?: string;
  col?: number;
  row?: number
}

let $terminal: Nullable<Terminal>
let $fitAddon: Nullable<FitAddon>
const terminalRef = ref<HTMLElement | null>(null);
const loading = ref(false);
const {setTitle} = useTabs();
const route = useRoute();

const serverId = parseInt(route.params?.id as string)
const serverHost = route.params?.server
setTitle("终端:" + serverHost)
const userStore = useUserStore();

//const { status, send, close:wsClose, open: wsOpen } = useWebSocket("wss://echo.websocket.events", {
const {status, send, close: wsClose, open: wsOpen} = useWebSocket(getServerSshWs(serverId), {
  autoReconnect: false,
  heartbeat: {
    message: '{"typ":"ping"}',
    interval: 3000
  },
  immediate: false,
  protocols: [userStore.getToken, userStore.getCurrentSpaceId.toString()],
  onConnected: () => {
    resizeTerm();
  },
  onError: (ws, event) => {
    termWrite(`\x1b[31mwebsocket连接失败：${ws.url}\x1b[m\r\n`);
  },
  onDisconnected: () => {
    termWrite(`\x1b[31mwebsocket关闭\x1b[m\r\n`);
  },
  onMessage: (ws, event) => {
    termWrite(event.data)
  }
})

const showMsg = (msg: string) => {
  console.log(msg);
};

//终端写入消息
const termWrite = (msg: string) => {
  if ($terminal) {
    $terminal.element && $terminal.focus();
    $terminal.write(msg);
  } else {
    showMsg('xterm已关闭，无法写入消息');
  }
};

//ws发送消息
const wsSend = (msg: wsMsg) => {
  if (status.value === 'OPEN') {
    send(JSON.stringify(msg));
  } else {
    termWrite("\x1b[31m连接已关闭，无法发送指令！\x1b[m\r\n")
  }
};

//监听窗口变化
const resizeTerm = () => {
  if ($fitAddon && $terminal) {
    $fitAddon.fit();
    $terminal.scrollToBottom();
    let dm = $fitAddon.proposeDimensions();
    if (dm && dm.cols != NaN) {
      wsSend({
        typ: 'resize',
        col: dm.cols,
        row: dm.rows,
      });
    }
  }
};

watchEffect(() => {
  if (status.value == "CONNECTING") {
    loading.value = true
  } else {
    loading.value = false
  }
})

const toggleConnect = () => {
  if (status.value == "OPEN") {
    wsClose()
  } else if (status.value == "CLOSED") {
    wsOpen()
  }
}

const initTerm = () => {
  $terminal = new Terminal({
    cursorBlink: true,
    //cursorStyle: "underline",
    screenReaderMode: true,
  });
  if (!$terminal || !terminalRef.value) {
    console.log("terminal初始化错误！")
    return
  }
  $terminal.open(terminalRef.value);
  $fitAddon = new FitAddon();
  $terminal.loadAddon($fitAddon);
  //监听窗口变动并且同步
  window.addEventListener('resize', resizeTerm);

  $terminal.writeln('等待连接。。。');
  $terminal.onData((message: string) => {
    wsSend({
      typ: "cmd",
      cmd: message,
    });
  });
};

onMounted(() => {
  initTerm();
  if (!serverId) {
    termWrite(`\x1b[31m服务器选择错误，无法连接终端！\x1b[m\r\n`);
    return
  }
  if ($terminal) {
    $terminal.focus();
  }
  wsOpen()
});
onBeforeUnmount(() => {
  //监听窗口变动并且同步
  window.removeEventListener('resize', resizeTerm);
  if (status.value === 'OPEN') {
    wsClose()
  }
  $fitAddon = null;
  $terminal = null;
});
</script>

<style scoped>
#showTip {
  z-index: 10000;
  position: absolute;
  right: 20px;
  top: 2px;
  background: #bbb;
  padding: 2px 3px;
}

#showTip:hover {
  background: #ffff;
}
</style>
