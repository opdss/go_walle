<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate"> 新增 </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                onClick: handleEdit.bind(null, record),
              },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                popConfirm: {
                  title: '是否确认删除',
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts">
  import { defineComponent } from 'vue';
  import { BasicTable, useTable, TableAction } from '/@/components/Table';
  import { columns, searchFormSchema } from './data';
  import { useMessage } from '/@/hooks/web/useMessage';
  import {getProjectListByPage, deleteProject} from "/@/api/project"
  import { useGo } from '/@/hooks/web/usePage';

  export default defineComponent({
    name: 'ProjectManagement',
    components: { BasicTable, TableAction },
    setup() {
      const go = useGo()
      const { createMessage } = useMessage();
      const [registerTable, { reload }] = useTable({
        title: '项目管理',
        api: getProjectListByPage,
        columns,
        formConfig: {
          labelWidth: 120,
          schemas: searchFormSchema,
        },
        useSearchForm: true,
        showTableSetting: true,
        bordered: true,
        showIndexColumn: false,
        actionColumn: {
          width: 80,
          title: '操作',
          dataIndex: 'action',
          // slots: { customRender: 'action' },
          fixed: "right",
        },
      });

      function handleCreate() {
       go("/project/create")
      }

      function handleEdit(record: Recordable) {
        go(`/project/update/${record.id}`)
      }

      function handleDelete(record: Recordable) {
        deleteProject(record.id).then(()=>{
          createMessage.success("删除成功")
          reload()
        })
      }

      function handleSuccess() {
        reload();
      }

      return {
        registerTable,
        handleCreate,
        handleEdit,
        handleDelete,
        handleSuccess,
      };
    },
  });
</script>
