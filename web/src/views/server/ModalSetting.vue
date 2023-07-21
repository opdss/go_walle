<template>
  <BasicModal
    v-bind="$attrs"
    :destroyOnClose="true"
    @register="register"
    :canFullscreen="false"
    title="自动设置免登陆公钥"
    autoclose
    @ok="handleSubmit"
    :maskClosable="false"
  >
      <Alert type="info" v-if="showInfo" :description="showMsg"></Alert>
      <template v-if="!showInfo">
        <Alert v-if="!errMsg"
               message="设置成功"
               description="现在可以顺畅的在这个服务器上部署代码了"
               type="success"
               show-icon></Alert>
        <Alert v-if="errMsg"
               message="设置失败"
               :description="`失败原因：${errMsg}`"
               type="error"
               show-icon></Alert>
      </template>
    <BasicForm
      style="margin-top: 20px"
      layout="vertical" @register="registerForm"/>

  </BasicModal>
</template>
<script lang="ts">
import {defineComponent, ref} from 'vue';
import {BasicModal, useModalInner} from '/@/components/Modal';
import {BasicForm, FormSchema, useForm} from '/@/components/Form/index';
import {setLoginServer} from '/@/api/server'
import {SetLogin} from '/@/api/server/model'
import {Alert} from 'ant-design-vue'

const schemas: FormSchema[] = [
  {
    field: 'id',
    component: 'InputNumber',
    label: 'id',
    show: false,
    defaultValue: 0,
    required: true,
  },
  {
    field: 'password',
    component: 'InputPassword',
    label: '输入密码',
    colProps: {
      span: 24,
    },
    required: true,
  },
];
export default defineComponent({
  components: {BasicModal, Alert, BasicForm},
  setup() {
    const loading = ref(true);
    const errMsg = ref("")
    const showMsg = ref('')
    const showInfo = ref(true)
    const [
      registerForm,
      {
        setFieldsValue,
        validate
      },
    ] = useForm({
      labelWidth: 120,
      schemas,
      showActionButtonGroup: false,
      actionColOptions: {
        span: 24,
      },
    });
    const [register, {setModalProps}] = useModalInner((data) => {
      if (!data) {
        errMsg.value = "出错了，请刷新重试！"
        return
      }
      showInfo.value = true
      errMsg.value = ""
      showMsg.value = `请输入[${data.user}@${data.host}:${data.port}]的登陆密码，此密码仅用来登陆服务器设置公钥实现后续的免密码登陆，不会做任何保存操作！如果该服务器不支持密码远程登陆，则请自行设置免密登陆`
      setFieldsValue({id: data.id as number})
    });

    async function handleSubmit() {
      try {
        let values = await validate();
        loading.value = true;
        setModalProps({loading: true, confirmLoading: true});
        setLoginServer(values as SetLogin).then(() => {
          errMsg.value = ""
        }).catch(e => {
          errMsg.value = e.toString()
        }).finally(() => {
          showInfo.value = false
          loading.value = false;
          setModalProps({loading: false, confirmLoading: false});
        })
      } catch (e) {
      }
    }

    return {register, loading, errMsg, showMsg, showInfo, registerForm, handleSubmit};
  },
});
</script>
