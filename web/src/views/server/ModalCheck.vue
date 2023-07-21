<template>
  <BasicModal
    v-bind="$attrs"
    destroyOnClose
    @register="register"
    :canFullscreen="false"
    title="检验服务器连接"
    autoclose
    @ok="handleSubmit"
    :minHeight="100"

  >
    <template v-if="loading">
      <div style="line-height: 100px;text-align: center;">连接检查中，请稍等……</div>
    </template>
    <template v-if="!loading">
      <Alert v-if="!errMsg"
             message="连接成功"
             description="现在可以顺畅的在这个服务器上部署代码了"
             type="success"
             show-icon></Alert>
      <Alert v-if="errMsg"
             message="连接失败"
             :description="`失败原因：${errMsg}`"
             type="error"
             show-icon></Alert>
    </template>
  </BasicModal>
</template>
<script lang="ts">
import {defineComponent, ref} from 'vue';
import {BasicModal, useModalInner} from '/@/components/Modal';
import {checkServer} from '/@/api/server'
import {Alert} from 'ant-design-vue'

export default defineComponent({
  components: {BasicModal, Alert},
  emits: ['success', 'register'],
  setup(_, {emit}) {
    const loading = ref(true);
    const errMsg = ref("")
    const [register, {setModalProps, closeModal}] = useModalInner(async (data) => {
      if (!data) {
        errMsg.value = "出错了，请刷新重试！"
        return
      }
      loading.value = true;
      setModalProps({loading: true, confirmLoading: true});
      await checkServer(data.id)
        .then(()=>{
          errMsg.value = ""
        })
        .catch(e => {
          errMsg.value = e.toString()
        })
        .finally(() => {
          loading.value = false;
          setModalProps({loading: false, confirmLoading: false});
        })
    });
    const handleSubmit = () => {
      try {
        emit("success")
      } catch (e) {

      }finally {
        closeModal()
      }
    }
    return {register, loading, errMsg, handleSubmit};
  },
});
</script>
