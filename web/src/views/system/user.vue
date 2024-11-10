<template>
  <div class="app-container">

    <div class="header">
      <el-input v-model="queryParams.keyword" placeholder="关键字查询:自动识别字段" class="small" style="width: 300px;margin-right: 10px" @keyup.enter.native="getList">
        <template v-slot:append>
          <el-button icon="el-icon-search" @click="getList" />
        </template>
      </el-input>

      <el-button v-if="hasPermission('add_user')" class="small" type="primary" @click="handleCreate">新增</el-button>
      <el-button v-if="hasPermission('del_user')" class="small" type="danger" :disabled="multiple" @click="handleDelete">删除</el-button>
    </div>

    <el-table v-loading="loading" border :data="userList" style="margin-top: 20px" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="45" align="center" :selectable="isUsernameAdmin" />
      <el-table-column show-overflow-tooltip label="编号" width="50" align="center" prop="id" />
      <el-table-column show-overflow-tooltip label="创建时间" align="center" prop="created_at" width="180">
        <template v-slot="scope">
          <div>{{ formatTimestamp(scope.row.created_at) }}</div>
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip label="更新时间" align="center" prop="updated_at" width="180">
        <template v-slot="scope">
          <div>{{ formatTimestamp(scope.row.updated_at) }}</div>
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip label="用户名称" align="center" prop="username" width="140" />
      <el-table-column show-overflow-tooltip label="用户昵称" align="center" prop="nickname" width="140" />
      <el-table-column show-overflow-tooltip label="手机号码" align="center" prop="phone" width="150" />
      <el-table-column show-overflow-tooltip label="邮箱" align="center" prop="email" width="150" />
      <el-table-column show-overflow-tooltip label="描述" align="center" prop="description" width="150" />
      <el-table-column show-overflow-tooltip label="角色" align="center" width="161">
        <template v-slot="scope">
          <div>
            <span v-if="isObjectArray(scope.row.roles)">
              {{ scope.row.roles.map(role => role.rolename).join(' | ') }}
            </span>
            <span v-else>
              {{ scope.row.roles.join(' | ') }}
            </span>
          </div>
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip label="状态" width="68" align="center">
        <template v-slot="scope">
          <el-switch v-model="scope.row.status" @change="handleStatusChange(scope.row)" />
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip label="操作" align="center" width="240" class-name="small-padding fixed-width">
        <template v-slot="scope">
          <el-button v-if="hasPermission('upd_user')" size="mini" type="text" icon="el-icon-edit" @click="handleUpdate(scope.row)">编辑</el-button>
          <el-button v-if="hasPermission('del_user') && scope.row.username !== 'admin'" size="mini" type="text" icon="el-icon-delete" @click="handleDelete(scope.row)">删除</el-button>
          <el-button v-if="hasPermission('reset_user')" size="mini" type="text" icon="el-icon-key" @click="handleResetPwd(scope.row)">重置</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="queryParams.page" :limit.sync="queryParams.perPage" @pagination="getList" />
    <!-- 添加或修改参数配置对话框 -->
    <el-dialog :title="title" :visible.sync="open" width="800px">
      <el-form ref="form" :model="form" :rules="rules" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="账户" prop="username" style="width: 90%">
              <el-input v-model="form.username" :disabled="form.id !== undefined" placeholder="请输入用户名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="昵称" prop="nickname" style="width: 90%">
              <el-input v-model="form.nickname" :disabled="form.id !== undefined" placeholder="请输入用户昵称(中文名称)" />
            </el-form-item>
          </el-col>
          <el-col :span="12" />
          <el-col :span="12">
            <el-form-item label="手机" prop="phone" style="width: 90%">
              <el-input v-model="form.phone" placeholder="请输入手机号码" maxlength="11" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email" style="width: 90%">
              <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色" prop="roles" style="width: 90%">
              <el-select v-model="form.roles" multiple placeholder="请选择" style="width: 100%" @change="$forceUpdate()">
                <el-option v-for="item in roleOptions" :key="item.id" :label="item.rolename" :value="item.rolename" :disabled="!item.status" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item v-if="form.id == undefined" label="密码" prop="password" style="width: 90%">
              <el-input v-model="form.password" placeholder="请输入密码,为空则系统默认生成" type="password" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="描述" style="width: 95%">
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
import { addUser, deleteUser, updateUser, listUser, detailUser } from '@/api/system/user'
import { listRole } from '@/api/system/role'
import { formatTimestamp } from '@/utils/timestamp'
import { mapGetters } from 'vuex'

export default {
  name: 'User',
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
      // 角色选项
      roleOptions: [],
      // 用户表格数据
      userList: null,
      // 总条数
      total: 0,
      // 查询参数
      queryParams: {
        page: 1,
        perPage: 10,
        keyword: undefined,
        username: undefined,
        phone: undefined
      },
      // 表单校验
      rules: {
        username: [
          { required: true, message: '用户名称不能为空', trigger: 'blur' }
        ],
        nickname: [
          { required: true, message: '用户昵称不能为空', trigger: 'blur' }
        ],
        roles: [
          { required: true, message: '用户角色不能为空', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '邮箱地址不能为空', trigger: 'blur' },
          {
            type: 'email',
            message: "'请输入正确的邮箱地址",
            trigger: ['blur', 'change']
          }
        ],
        phone: [
          { required: true, message: '手机号码不能为空', trigger: 'blur' },
          {
            pattern: /^1[3|4|5|6|7|8|9][0-9]\d{8}$/,
            message: '请输入正确的手机号码',
            trigger: 'blur'
          }
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
    /** 查询用户列表 */
    getList() {
      this.loading = true
      listUser(this.queryParams).then(
        (response) => {
          this.userList = response.items
          this.total = response.total
          this.loading = false
        }
      )
    },
    /** 检查数组是否包含对象 */
    isObjectArray(arr) {
      return Array.isArray(arr) && arr.length > 0 && typeof arr[0] === 'object'
    },
    /** 表单重置 */
    reset() {
      this.form = {
        id: undefined,
        username: undefined,
        nickname: undefined,
        password: undefined,
        phone: undefined,
        email: undefined,
        description: undefined,
        roles: undefined,
        status: true
      }
    },
    /* 处理返回时间的时间戳转化*/
    formatTimestamp(timestamp) { // 转化时间戳
      return formatTimestamp(timestamp)
    },
    /* 处理用户更新*/
    handleUpdate(row) {
      this.reset()
      this.queryParams.perPage = 1000
      listRole(this.queryParams).then(res => {
        this.roleOptions = res.items
      })

      const id = row.id
      detailUser(id).then(res => {
        this.open = true
        this.title = '修改用户'
        this.form = res
        this.form.roles = res.roles.map(item => item.rolename)
        this.form.password = res.password
      })
    },
    /* 处理用户删除*/
    handleDelete(row) {
      const userIds = row.id || this.ids
      this.$confirm(
        '是否确认删除用户编号为"' + userIds + '"的数据项?',
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
        .then(() => {
          return deleteUser(userIds)
        })
        .then(() => {
          this.getList()
          this.msgSuccess('删除成功')
        })
        .catch(() => { })
    },
    /* 处理重置按钮,用于重置用户密码 */
    handleResetPwd(row) {
      this.reset()
      let data = this.form
      data = row
      this.$prompt('请输入"' + row.username + '"的新密码', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      })
        .then(({ value }) => {
          data.password = value
          data.roles = data.roles.map(item => item.rolename)
          updateUser(row.id, data).then(() => {
            this.msgSuccess('密码重置修改成功，新密码是：' + value)
          })
        })
        .catch(() => { })
    },
    /* 处理行的状态变化*/
    handleStatusChange(row) {
      const text = row.status ? '启用' : '停用'
      this.reset()
      let data = this.form
      data = row
      this.$confirm(
        '确认要"' + text + '""' + row.username + '"用户吗?',
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
        .then(() => {
          data.roles = data.roles.map(item => item.rolename)
          return updateUser(row.id, data)
        })
        .then(() => {
          this.msgSuccess(text + '成功')
        })
        .catch(() => {
          row.status = !row.status
        })
    },
    /* 处理selection的行选择*/
    handleSelectionChange(selection) {
      this.ids = selection.map((item) => item.id)
      this.single = selection.length !== 1
      this.multiple = !selection.length
    },
    /* 添加用户*/
    handleCreate() {
      this.reset()
      this.queryParams.perPage = 1000
      listRole(this.queryParams).then(res => {
        this.roleOptions = res.items
      })
      this.open = true
      this.title = '添加用户'
      this.form.password = ''
    },
    /* 处理对话框的提交按钮*/
    handleSubmitForm() {
      this.$refs['form'].validate((valid) => {
        if (valid) {
          this.loading = true
          if (this.form.id !== undefined) {
            updateUser(this.form.id, this.form).then(res => {
              this.msgSuccess(`用户更新成功`)
              this.open = false
              this.getList()
            })
          } else {
            this.loading = true
            addUser(this.form).then(res => {
              this.msgSuccess(`用户[${res.username}]新增成功`)
              this.open = false
              this.getList()
            })
          }
        }
      })
    },
    /* 处理对话框的取消按钮*/
    handleCancel() {
      this.reset()
      this.open = false
    },
    /* 判断这行是否是admin用户还是其他用户*/
    isUsernameAdmin(row, index) {
      return row.username !== 'admin'
    }
  }
}
</script>
