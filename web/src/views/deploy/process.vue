<template>
  <div>
    <div ref="terminalRef" style="width: 100%; height: 100%; padding: 1px 2px; background: black">
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import 'xterm/css/xterm.css';

export default defineComponent({
  name: 'DeployManagement',
  components: {  },
  setup() {
    const terminalRef = ref<HTMLElement | null>(null);
    var $term:Nullable<Terminal>
    var $fitAddon:Nullable<FitAddon>

    const initTerm = () => {
      $term = new Terminal({
        cursorBlink: true,
        //cursorStyle: "underline",
        screenReaderMode: true,
      });
      if (terminalRef.value) {
        $term.open(terminalRef.value);
      }
      $fitAddon = new FitAddon();
      $term.loadAddon($fitAddon);
      //监听窗口变动并且同步
      // window.addEventListener('resize', resizeTerm);
      //
      // $term.value.writeln('等待连接。。。');
      // $term.value.onData((message: string) => {
      //   if (!isWsOpen()) {
      //     initWs();
      //   }
      //   wsSend({
      //     type: msgType.CMD,
      //     cmd: message,
      //     cols: 0,
      //     rows: 0,
      //   });
      //});
    };

    onMounted(() => {
      console.log('onMounted start');
      initTerm();
      if ($term) {
        $term.focus();
      }
    });

    return {
      terminalRef,
    };
  },
});
</script>
