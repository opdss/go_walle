<template>
  <PageWrapper contentFullHeight>
    <Card  v-loading="loading">
      <Tabs v-model:activeKey="activeKey">
        <TabPane
          v-for="item in tabsFormSchema"
          :key="item.key"
          :tab="item.tab"
          v-bind="omit(item, ['Form', 'key', 'RegisterFn'])"
        >
          <BasicForm @register="item.Form[0]" />
        </TabPane>
        <template #rightExtra>
          <a-button @click="handleReset" class="mr-5">重置</a-button>
          <a-button @click="handleSubmit" type="primary">提交</a-button>
        </template>
      </Tabs>
    </Card>
  </PageWrapper>
</template>
<script lang="ts" setup>
import {onMounted, ref} from 'vue';
import {useMessage} from '/@/hooks/web/useMessage';
import {PageWrapper} from '/@/components/Page';
import {formSchemas} from './data'
import {createProject, detailProject, updateProject} from '/@/api/project'
import {CreateReq, UpdateReq} from '/@/api/project/model'
import {Card, Tabs, TabPane} from 'ant-design-vue'
import { omit } from 'lodash-es';
import { deepMerge } from '/@/utils';
import { useGo } from '/@/hooks/web/usePage';
import { BasicForm, FormSchema, useForm, FormProps, UseFormReturnType } from '/@/components/Form';
import {useRoute} from "vue-router";

type TabsFormType = {
  key: string;
  tab: string;
  forceRender?: boolean;
  Form: UseFormReturnType;
};
// 公共属性
const baseFormConfig: Partial<FormProps> = {
  showActionButtonGroup: false,
  labelWidth: 180,
  baseColProps: { span: 24 },
};

const go = useGo()
const route = useRoute();
const projectId = parseInt(route.params?.id as string)
const isUpdate = route.name == "ProjectUpdate"
const isDetail =  route.name == "ProjectDetail"
const loading = ref(false);
const {createMessage} = useMessage();
const activeKey = ref<string>("")
const tabsFormSchema:TabsFormType[] = []
for (let k  in formSchemas) {
  let form = useForm(Object.assign({ schemas: formSchemas[k] }, baseFormConfig) as FormProps)
  if (activeKey.value == "") {
    activeKey.value = k
    // if (projectId > 0) {
    //   const {appendSchemaByField}  = form[1]
    //   appendSchemaByField({
    //     label:"id",
    //     field:"id",
    //     defaultValue:projectId,
    //     show:false,
    //     component:"InputNumber",
    //   }, "")
    // }
  }
  tabsFormSchema.push({
    key: k,
    tab:k,
    forceRender: true,
    Form: form,
  })
}

async function handleReset() {
  for (const item of tabsFormSchema) {
    const { resetFields } = item.Form[1];
    await resetFields();
  }
}

async function handleSubmit() {
  let lastKey = '';
  loading.value = true;
  try {
    let values = {};
    for (const item of tabsFormSchema) {
      lastKey = item.key;
      const { validate, getFieldsValue } = item.Form[1];
      await validate();
      values = deepMerge(values, getFieldsValue());
    }
    console.log("values", values)
    if (isUpdate) {
      //提交后台
      values['id'] = projectId
      updateProject(values as UpdateReq).then(() => {
        createMessage.success("更新成功");
        go("/project/index")
      }).catch()
    } else {
      //提交后台
      createProject(values as CreateReq).then(() => {
        createMessage.success("创建成功");
        go("/project/index")
      }).catch()
    }

  } catch (e) {
    // 验证失败或出错，切换到对应标签页
    activeKey.value = lastKey;
  } finally {
    loading.value = false;
  }
}

async function initProject(projectId:number) {
  try{
    let res = await detailProject(projectId)
    res.server_ids = []
    res.servers.forEach((val) => {
      res.server_ids.push(val.id)
    })
    for (const item of tabsFormSchema) {
      const { setFieldsValue, setProps  } = item.Form[1];
      await setFieldsValue(res)
      if (isDetail) {
        await setProps({disabled: true})
      }
    }
  } catch {

  } finally{

  }
}

onMounted(()=>{
  if (isUpdate || isDetail) {
    initProject(projectId)
  }
})
</script>
