<template>
  <Dropdown
    placement="bottom"
    :trigger="['click']"
    :dropMenuList="spaceList"
    :selectedKeys="selectedKeys"
    @menu-event="handleMenuEvent"
    overlayClassName="app-locale-picker-overlay"
  >
    <span class="cursor-pointer flex items-center">
      <span class="ml-1">{{ getSpaceText }}</span>
    </span>
  </Dropdown>
</template>
<script lang="ts" setup>
  import type { DropMenu } from '/@/components/Dropdown';
  import { ref, watchEffect, unref, computed } from 'vue';
  import { Dropdown } from '/@/components/Dropdown';
  import { useUserStore } from '/@/store/modules/user';

  const userStore = useUserStore();

  const spaces = userStore.getSpaces

  const selectedKeys = ref<number[]>([])

  const spaceList = ref<DropMenu[]>([])
  if (spaces.length > 0) {
    spaces.forEach((v) => {
      spaceList.value = [...spaceList.value, {
        event: v.space_id,
        text: v.space_name,
        disabled: v.status != 1 ,
      }]
    })
  }

  const getSpaceText = computed(() => {
    const key = selectedKeys.value[0];
    if (!key) {
      return '';
    }
    return spaceList.value.find((item) => item.event == key)?.text;
  });

  watchEffect(() => {
    selectedKeys.value = [userStore.getCurrentSpaceId];
  });

  async function toggleSpace(spaceId: number|string) {
    selectedKeys.value = [spaceId as number];
    userStore.setCurrentSpaceId(spaceId as number)
    location.reload();
  }

  function handleMenuEvent(menu: DropMenu) {
    if (unref(userStore.getCurrentSpaceId) === menu.event) {
      return;
    }
    toggleSpace(menu.event);
  }
</script>

<style lang="less">
  .app-locale-picker-overlay {
    .ant-dropdown-menu-item {
      min-width: 160px;
    }
  }
</style>
