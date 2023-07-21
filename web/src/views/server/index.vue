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
            :dropDownActions="dropDownActions(record)"
          />
        </template>
      </template>
    </BasicTable>
    <ServerModal @register="registerModal" @success="handleSuccess" />
    <ModalCheck @register="registerModalCheck" @success="handleSuccess" />
    <ModalSetting @register="registerModalSetting" @success="handleSuccess" />
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import {ActionItem, BasicTable, TableAction, useTable} from '/@/components/Table';
import {useModal} from '/@/components/Modal';
import ServerModal from './Modal.vue';
import {columns, searchFormSchema} from './data';
import {deleteServer, getServerListByPage} from "/@/api/server"
import {useMessage} from '/@/hooks/web/useMessage';
import {useGo} from '/@/hooks/web/usePage';
import ModalCheck from './ModalCheck.vue'
import ModalSetting from './ModalSetting.vue'
import {Status} from "/@/enums/fieldEnum";
import { RoleMaster} from "/@/enums/roleEnum";

export default defineComponent({
  name: 'ServerManagement',
  components: { BasicTable, ServerModal, TableAction, ModalCheck, ModalSetting },
  setup() {
    const go = useGo()
    const {createMessage} = useMessage()
    const [registerModal, { openModal }] = useModal();
    const [registerModalCheck, {openModal: openModalCheck}] = useModal();
    const [registerModalSetting, {openModal: openModalSetting}] = useModal();
    const [registerTable, { reload }] = useTable({
      title: '服务器列表',
      api: getServerListByPage,
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
        title: '操作',
        dataIndex: 'action',
        // slots: { customRender: 'action' },
        fixed: undefined,
      },
    });

    function handleCreate() {
      openModal(true, {
        isUpdate: false,
      });
    }

    function handleEdit(record: Recordable) {
      openModal(true, {
        record,
        isUpdate: true,
      });
    }

    function handleDelete(record: Recordable) {
      deleteServer(record.id).then(()=>{
        createMessage.success("删除成功")
        reload()
      })
    }

    function handleSuccess() {
      reload();
    }

    function handleSetting(record: Recordable) {
      openModalSetting(true, record)
    }

    function handleCheck(record: Recordable) {
      openModalCheck(true, record)
    }

    function handleTerminal(record: Recordable) {
      go(`/server/terminal/${record.id}/${record.user}@${record.host}:${record.port}`)
    }

    function dropDownActions(record: Recordable):ActionItem[]{
      return [
        {
          icon: 'ant-design:code-outlined',
          label: '终端',
          auth:RoleMaster,
          disabled: record.status !== Status.Enabled,
          onClick: handleTerminal.bind(null, record)
        },
        {
          icon:'ant-design:check-square-outlined',
          label:'检测',
          onClick: handleCheck.bind(null, record)
        },
        {
          icon:'ant-design:cloud-upload-outlined',
          label:'配置免登陆',
          auth:RoleMaster,
          disabled: record.status === Status.Enabled,
          tooltip: '将本机ssh公钥同步到目标服务器，实现免登陆连接',
          onClick: handleSetting.bind(null, record)
        },
      ]
    }

    return {
      registerTable,
      registerModal,
      registerModalCheck,
      registerModalSetting,
      handleCreate,
      handleEdit,
      handleDelete,
      handleCheck,
      handleSuccess,
      dropDownActions,
    };
  },
});
</script>
