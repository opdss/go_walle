<template>
  <div >
    <CollapseContainer title="添加成员" class="m-5">
      <BasicForm @register="registerForm" class="pt-5" @submit="handleCreate"/>
    </CollapseContainer>
    <BasicTable @register="registerTable" class="m-5">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate"> 新增 </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
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
import {defineComponent} from 'vue';
import { BasicTable, useTable, TableAction } from '/@/components/Table';
import { BasicForm, useForm } from '/@/components/Form/index';
import { useModal } from '/@/components/Modal';
import {columns} from './data';
import {getMemberListByPage, deleteMember, storeMember} from "/@/api/member"
import {getUserOptions} from "/@/api/user"
import { useMessage } from '/@/hooks/web/useMessage';
import { CollapseContainer } from '/@/components/Container';
import {RoleEnum} from "/@/enums/roleEnum";

export default defineComponent({
  name: 'UserManagement',
  components: { BasicTable, TableAction,BasicForm , CollapseContainer},
  setup() {
    const { createMessage } = useMessage();
    const [registerModal, { openModal }] = useModal();
    const [registerTable, { reload }] = useTable({
      title: '用户列表',
      api: getMemberListByPage,
      columns,
      showTableSetting: true,
      bordered: true,
      showIndexColumn: false,
      actionColumn: {
        title: '操作',
        dataIndex: 'action',
        // slots: { customRender: 'action' },
        fixed: undefined,
      },
    });
    const [
      registerForm,
      {
        resetFields,
        validate,
      },
    ] = useForm({
      labelWidth: 100,
      submitButtonOptions:{
        text:"添加",
      },
      schemas:  [
        {
          field: 'user_id',
          label: '选择用户',
          required: true,
          component: 'ApiSelect',
          componentProps: {
            api: getUserOptions,
            resultField: 'options',
            labelField: 'text',
            valueField: 'value',
          },
          colProps: {
            span:8,
          }
        },
        {
          field: 'role',
          label: '角色',
          required: true,
          component: 'Select',
          componentProps:{
            options:[
              {
                label:"项目管理员",
                value:RoleEnum.MASTER
              },
              {
                label:"开发者",
                value:RoleEnum.DEVELOPER
              }]
          },
          colProps: {
            span:8,
          }
        },
      ],
    });

    async function handleCreate() {
      try {
        const values = await validate();
        storeMember(values)
          .then(() => {
            createMessage.success("添加成功")
            resetFields()
          })
      }
      catch {
        createMessage.error("添加失败")
      }finally {
      }
    }

    function handleEdit(record: Recordable) {
      openModal(true, {
        record,
        isUpdate: true,
      });
    }

    function handleDelete(record: Recordable) {
      deleteMember(record.id).then(()=>{
        createMessage.success("删除成功")
        reload()
      })
    }

    function handleSuccess() {
      reload();
    }

    return {
      registerTable,
      registerModal,
      registerForm,
      handleCreate,
      handleEdit,
      handleDelete,
      handleSuccess,
    };
  },
});
</script>
