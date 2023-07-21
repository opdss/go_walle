<template>
  <PageWrapper contentFullHeight>
    <a-card :title="title">
      <a-alert
        v-show="!!alertMsg"
        message="Error"
        :description="alertMsg"
        type="error"
        show-icon
      />
      <BasicForm
        v-show="!alertMsg"
        autoFocusFirstItem
        :labelWidth="200"
        :baseColProps="{span:24}"
        :schemas="formSchema"
        :actionColOptions="{ span: 24 }"
        @submit="handleSubmit"
        ref="formElRef"
      >
      </BasicForm>
      <Loading :loading="loading" :absolute="true" tip="加载中..."/>
    </a-card>
  </PageWrapper>
</template>
<script lang="ts">
import {defineComponent, onMounted, ref} from 'vue';
import {BasicForm, FormActionType} from '/@/components/Form/index';
import {useMessage} from '/@/hooks/web/useMessage';
import {PageWrapper} from '/@/components/Page';
import {formSchema} from './data'
import {detailProject} from '/@/api/project'
import {createDeploy} from '/@/api/deploy'
import {Detail} from '/@/api/project/model'
import {useRoute} from 'vue-router';
import {Alert as AAlert, Card as ACard} from "ant-design-vue"
import {Loading} from '/@/components/Loading';
import {useGo} from "/@/hooks/web/usePage"
import { useTabs } from '/@/hooks/web/useTabs';

export default defineComponent({
  components: {BasicForm, PageWrapper, ACard, AAlert, Loading},
  setup() {
    const {createMessage} = useMessage();
    const route = useRoute();
    const projectDetail = ref<Detail | null>(null)
    const formElRef = ref<Nullable<FormActionType>>(null);
    const title = ref("新建上线单")
    const alertMsg = ref("")
    const loading = ref(false)
    const go = useGo()
    const {  setTitle } = useTabs();

    async function getProjectDetail(projectId: number) {
      if (!projectId) {
        alertMsg.value = "选择项目错误，请重新选择！"
        return
      }
      try {
        loading.value = true
        projectDetail.value = await detailProject(projectId)
      } catch (e) {
        alertMsg.value = "你准备上线的项目不存在或者不允许上线操作，请重新检查项目或者联系相关负责人！"
        return
      } finally {
        loading.value = false
      }
      title.value = `新建上线单-[${alertMsg.value?'项目选择错误':projectDetail.value?.name}]`
      setTitle(title.value)
      if (projectDetail.value && formElRef.value) {
        formElRef.value.appendSchemaByField({
          field: 'project_id',
          label: '',
          defaultValue: projectId,
          required: true,
          show: false,
          component: 'InputNumber',
        }, undefined)
        //渲染服务器选项
        let options: object[] = []
        let checked: number[] = []
        projectDetail.value.servers.forEach((v) => {
          options = [...options, {
            value: v.id,
            label: `${v.name}[${v.user}@${v.host}:${v.port}]`,
          }]
          checked = [...checked, v.id]
        })
        formElRef.value.updateSchema(
          {
            field: 'server_ids',
            defaultValue: checked,
            componentProps: {
              options: options,
              mode: "multiple",
            },
          })
        //branch/tag
        let isTag = projectDetail.value.repo_mode == "tag"
        formElRef.value.updateSchema([
          {
            field: 'tag',
            show: isTag,
            required: isTag,
          },
          {
            field: 'branch',
            show: !isTag,
            required: !isTag,
          }, {
            field: 'commit_id',
            show: !isTag,
            required: !isTag,
          }
        ])
      }
    }

    const handleSubmit = (values: any) => {
      createDeploy(values).then(() => {
        createMessage.success("创建上线单成功！");
        go("/deploy/index")
      })
    }

    onMounted(() => {
      let projectId = parseInt(route.params?.project_id as unknown as string)
      getProjectDetail(projectId)
    })

    return {
      title,
      loading,
      alertMsg,
      formElRef,
      formSchema,
      projectDetail,
      handleSubmit,
    };
  },
});
</script>
