<template>
  <div class="app-container">

    <!-- 搜索/新增/删除框 -->
    <div class="header">
      <el-input placeholder="关键字查询:自动识别字段" :v-model="queryParams.keyword" class="small" style="width: 300px;margin-right: 10px" @keyup.enter.native="getList">
        <template v-slot:append>
          <el-button icon="el-icon-search" @click="getList" />
        </template>
      </el-input>

      <el-button v-if="hasPermission('add_role')" class="small" type="primary" @click="handleCreate">新增</el-button>
      <el-button v-if="hasPermission('del_role')" class="small" type="danger" :disabled="multiple" @click="handleDelete">删除</el-button>
    </div>

    <!-- 表格 -->
    <el-table v-loading="loading" border :data="roleList" style="margin-top: 20px" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="45" align="center" :selectable="isUsernameAdmin" />
      <el-table-column show-overflow-tooltip label="编号" width="50" align="center" prop="id" />
      <el-table-column show-overflow-tooltip label="创建时间" align="center" prop="created_at" width="150">
        <template v-slot="scope">
          <div>{{ formatTimestamp(scope.row.created_at) }}</div>
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip label="更新时间" align="center" prop="updated_at" width="150">
        <template v-slot="scope">
          <div>{{ formatTimestamp(scope.row.updated_at) }}</div>
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip label="角色名称" align="center" prop="rolename" width="140" />
      <el-table-column show-overflow-tooltip label="中文名称" align="center" prop="nickname" width="140" />
      <el-table-column show-overflow-tooltip label="描述" align="center" prop="description" width="150" />
      <el-table-column show-overflow-tooltip label="继承权限" align="center" width="150">
        <template v-slot="scope">
          <div>
            <span v-if="Array.isArray(scope.row.parents)"> {{ scope.row.parents.join(' | ') }} </span>
            <span v-else><!-- 显示空白 --></span>
          </div>
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip label="状态" width="150" align="center">
        <template v-slot="scope">
          <el-switch v-if="hasPermission('upd_role')" v-model="scope.row.status" @change="handleStatusChange(scope.row)" />
        </template>
      </el-table-column>

      <el-table-column show-overflow-tooltip label="用户" align="center" width="379">
        <template v-slot="scope">
          <div>{{ scope.row.users.map(item => item.username).join(' | ') }}</div>
        </template>
      </el-table-column>

      <el-table-column show-overflow-tooltip label="操作" align="center" width="150" class-name="small-padding fixed-width">
        <template v-slot="scope">
          <el-button v-if="hasPermission('upd_role')" size="mini" type="text" icon="el-icon-edit" @click="handleUpdate(scope.row)">编辑</el-button>
          <el-button v-if="hasPermission('del_role') && scope.row.rolename !== 'admin'" size="mini" type="text" icon="el-icon-delete" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <pagination v-show="total>0" :total="total" :page.sync="queryParams.page" :limit.sync="queryParams.perPage" @pagination="getList" />

    <!-- 添加或修改参数配置对话框 -->
    <el-dialog :title="title" :visible.sync="open" width="800px">
      <el-form ref="form" :model="form" :rules="rules" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="角色名称" prop="rolename" style="width: 90%">
              <el-input v-model="form.rolename" :disabled="form.id !== undefined" placeholder="请输入角色名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色昵称" prop="nickname" style="width: 90%">
              <el-input v-model="form.nickname" :disabled="form.id !== undefined" placeholder="请输入角色昵称(中文)" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="继承权限" prop="parents" style="width: 90%">
              <el-select v-model="form.parents" multiple placeholder="请选择" style="width: 100%" @change="$forceUpdate()">
                <el-option v-for="item in parentOptions" :key="item.id" :label="item.rolename" :value="item.rolename" :disabled="!item.status || item.rolename===form.rolename" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="对应用户" prop="users" style="width: 90%">
              <el-select v-model="form.users" multiple placeholder="请选择" style="width: 100%">
                <el-option v-for="item in userOptions" :key="item.id" :label="item.username" :value="item.username" :disabled="!item.status" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="接口权限" prop="endpoints" style="width: 90%">
              <el-cascader v-model="selectedValues" :options="endpointOptions" :props="{ multiple: true }" :show-all-levels="false" clearable filterable placeholder="请选择接口权限" @change="handleEndpointSelect" />
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item label="角色描述" style="width: 95%">
              <el-input v-model="form.description" type="textarea" placeholder="请输入描述内容" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button :loading="loading" type="primary" @click="handleSubmitForm">确 定</el-button>
        <el-button @click="handleCancel">取 消</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { addRole, deleteRole, updateRole, listRole, detailRole, treeEndpoint } from '@/api/system/role'
import { listUser } from '@/api/system/user'
import { formatTimestamp } from '@/utils/timestamp'
import { mapGetters } from 'vuex'

export default {
  name: 'Role',
  data() {
    return {
      // 遮罩层
      loading: true,
      // 选中数组
      ids: [],
      // 非单个禁用
      single: true,
      // 非多个禁用
      multiple: true,
      // 弹出层标题
      title: '',
      // 是否显示弹出层
      open: false,
      // 表单参数
      form: {},
      // 权限继承选项
      parentOptions: [],
      // 接口选项
      endpointOptions: [],
      // 接口选项绑定value
      selectedValues: [],
      // 用户选项
      userOptions: [],
      // 表格数据
      roleList: null,
      // 总条数
      total: 0,
      // 查询参数
      queryParams: {
        page: 1,
        perPage: 10,
        keyword: undefined,
        rolename: undefined,
        nickname: undefined
      },
      // 表单校验
      rules: {
        rolename: [
          { required: true, message: '角色名称不能为空', trigger: 'blur' }
        ],
        nickname: [
          { required: true, message: '角色昵称不能为空', trigger: 'blur' }
        ],
        endpoints: [
          { required: true, message: '角色需要分配权限', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters(['hasPermission'])
  },
  created() {
    this.getList()
  },
  methods: {
    /* 处理返回时间的时间戳转化 */
    formatTimestamp(timestamp) { // 转化时间戳
      return formatTimestamp(timestamp)
    },
    /** 列表查询 */
    getList() {
      this.loading = true
      listRole(this.queryParams).then(
        (response) => {
          this.roleList = response.items
          this.total = response.total
          this.loading = false
        }
      )
    },
    /** 表单重置 */
    reset() {
      this.selectedValues = []
      this.form = {
        id: undefined,
        rolename: undefined,
        nickname: undefined,
        description: undefined,
        status: undefined,
        parents: [],
        users: [],
        endpoints: []
      }
    },
    /* 处理更新 */
    handleUpdate(row) {
      this.reset()

      this.parentOptions = this.roleList
      treeEndpoint().then(res => {
        this.endpointOptions = res
      })
      this.queryParams.perPage = 1000
      listUser(this.queryParams).then(res => {
        this.userOptions = res.items
      })

      const id = row.id
      detailRole(id).then(res => {
        this.open = true
        this.form = res
        this.title = '修改角色'
        this.form.users = res.users.map(item => item.username)
        this.selectedValues = res.endpoints.map(item => [item.module, item.kind, item.identity])
        this.form.endpoints = res.endpoints.map(item => item.identity)
        this.parents = res.parents
      })
    },
    /* 处理删除 */
    handleDelete(row) {
      const ids = row.id || this.ids
      this.$confirm(
        '是否确认删除编号为"' + ids + '"的数据项?',
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
        .then(() => {
          return deleteRole(ids)
        })
        .then(() => {
          this.getList()
          this.msgSuccess('删除成功')
        })
        .catch(() => { })
    },
    /* 处理行的状态变化 */
    handleStatusChange(row) {
      const text = row.status ? '启用' : '停用'
      this.reset()
      let data = this.form
      data = row
      this.$confirm(
        '确认要"' + text + '""' + row.rolename + '"角色吗?',
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
        .then(() => {
          data.users = data.users.map(item => item.username)
          return updateRole(row.id, data)
        })
        .then(() => {
          this.msgSuccess(text + '成功')
        })
        .catch(() => {
          row.status = !row.status
        })
    },
    /* 处理selection的行选择 */
    handleSelectionChange(selection) {
      this.ids = selection.map((item) => item.id)
      this.single = selection.length !== 1
      this.multiple = !selection.length
    },
    /* 添加 */
    handleCreate() {
      this.reset()
      this.parentOptions = this.roleList
      treeEndpoint().then(res => {
        this.endpointOptions = res
      })
      this.queryParams.perPage = 1000
      listUser(this.queryParams).then(res => {
        this.userOptions = res.items
      })
      this.open = true
      this.title = '添加角色'
    },

    /* 处理树行选择 */
    handleEndpointSelect(value) {
      this.form.endpoints = value.map(item => item[item.length - 1])
    },
    /* 处理对话框的提交 */
    handleSubmitForm() {
      this.$refs['form'].validate((valid) => {
        if (valid) {
          this.loading = true
          if (this.form.id !== undefined) {
            updateRole(this.form.id, this.form).then(res => {
              this.msgSuccess(`更新成功`)
              this.open = false
              this.getList()
            })
          } else {
            this.loading = true
            addRole(this.form).then(res => {
              this.msgSuccess(`角色[${res.rolename}]新增成功`)
              this.open = false
              this.getList()
            })
          }
        }
      })
    },
    /* 处理对话框的取消 */
    handleCancel() {
      this.reset()
      this.open = false
    },
    /* 判断这行是否是admin用户还是其他用户*/
    isUsernameAdmin(row, index) {
      return row.rolename !== 'admin'
    }
  }
}
</script>
