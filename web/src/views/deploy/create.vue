<template>
  <PageWrapper contentFullHeight>
    <a-card title="新建上线单" v-loading="loadingRef" >
      <a-empty v-show="envEmpty"  class="pt-10">
        <template #description>
          <span> 没有项目可发布，去创建一个吧 </span>
        </template>
        <a-button type="primary" @click="go('/project/create')">创建项目</a-button>
      </a-empty>
      <a-tabs
        size="large"
        type="card"
        :style="{ minHeight: '400px' }"
        v-model:activeKey="activeKey"
        @tabClick="tabClick"
      >
        <a-tab-pane v-for="val in envList" :key="val.value" :tab="val.text" class="px-1">
          <a-empty v-show="projectEmpty" class="pt-10">
            <template #description>
              <span> 没有项目可发布，去创建一个吧 </span>
            </template>
            <a-button type="primary" @click="go('/project/create')">创建项目</a-button>
          </a-empty>
          <a-list v-show="!projectEmpty"
                  :grid="{ gutter: 10, xs: 1, sm: 2, md: 4, lg: 4, xl: 6, xxl: 8 }"
                  :data-source="projectList"
          >
              <template #renderItem="{ item }">
                  <a-list-item>
                    <a-card :hoverable="true"  @click="selectProject(item.value)">
                      <h2>{{ item.text }}</h2>
<!--                      <p >-->
<!--                        最近发布：<br/>-->
<!--                        发布人：wuxin<br/>-->
<!--                        发布时间：2014-12-44<br/>-->
<!--                        发布版本：v1.23.4<br/>-->
<!--                        发布标题：测试发布哇俄方哇俄方哇饿乏味乏味乏味乏味我乏味  我我-->
<!--                      </p>-->
                    </a-card>
                  </a-list-item>
              </template>
          </a-list>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </PageWrapper>
</template>
<script lang="ts">
import {defineComponent, onMounted, ref} from 'vue';
import {PageWrapper} from '/@/components/Page';

import {
  Card as ACard,
  Col as ACol,
  Empty as AEmpty,
  List as AList,
  ListItem as AListItem,
  Row as ARow,
  TabPane as ATabPane,
  Tabs as ATabs
} from 'ant-design-vue';
import {getEnvironmentOptions} from "/@/api/environment";
import {getProjectOptions} from "/@/api/project";
import {OptionItem} from "/@/api/model/baseModel"
import {createLocalStorage} from '/@/utils/cache';
import {DEF_ENV_CACHE_KEY} from "/@/enums/cacheEnum";
import {useGo} from "/@/hooks/web/usePage";

type envProjList = { [envId: number]: { list: OptionItem[], isGet: boolean } }

export default defineComponent({
  components: {PageWrapper, ACard, ATabs, ATabPane, AList, AListItem, ARow, ACol, AEmpty},
  setup() {
    const lCache = createLocalStorage();
    const envList = ref<OptionItem[]>([]);
    const activeKey = ref<number>(lCache.get(DEF_ENV_CACHE_KEY, 0) as number)
    const projectList = ref<OptionItem[]>([])
    const go = useGo()
    const envEmpty = ref(false)
    const projectEmpty = ref(false)
    const loadingRef = ref(false);


    const projectListByEnv: envProjList = {}

    const setDefEnvId = (envId: number) => {
      lCache.set(DEF_ENV_CACHE_KEY, envId)
      activeKey.value = envId
      setProjectList(envId)
    }

    async function initEnvList() {
      let res
      try {
        loadingRef.value = true
        res = await getEnvironmentOptions({page_size: 1000})
      } catch (e) {
      } finally {
        loadingRef.value = false
      }
      if (!res || res.options.length == 0) {
        envEmpty.value = true
        return
      }
      envList.value = res.options
      if (activeKey.value == 0) {
        setDefEnvId(res.options[0].value as number)
        return
      }
      let exists = false
      envList.value.forEach((v) => {
        if (v.value as number == activeKey.value) {
          exists = true
        }
      })
      if (!exists) {
        setDefEnvId(res.options[0].value as number)
        return
      }
      await setProjectList(activeKey.value)
    }


    async function setProjectList(envId: number) {
      if (projectListByEnv[envId]?.isGet) {
        projectList.value = projectListByEnv[envId].list
      } else {
        let res
        try {
          loadingRef.value = true
          res = await getProjectOptions({page_size: 1000, environment_id: envId})
        } finally {
          loadingRef.value = false
        }
        projectListByEnv[envId] = {
          isGet: true,
          list: res ? (res.options as OptionItem[]) : []
        }
        projectList.value = projectListByEnv[envId].list
      }
      if (projectList.value.length == 0) {
        projectEmpty.value = true
      } else {
        projectEmpty.value = false
      }
    }

    const selectProject = (projId) => {
      go(`/deploy/create/${projId}`)
    }


    const tabClick = (envId) => {
      setDefEnvId(envId)
    }

    onMounted(() => {
      initEnvList()
    })
    return {
      go,
      envEmpty,
      projectEmpty,
      projectList,
      activeKey,
      selectProject,
      tabClick,
      envList,
      loadingRef,
    };
  },
});
</script>

