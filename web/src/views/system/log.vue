<template>
  <div class="app-container">

    <!-- 搜索/新增/删除框 -->
    <div class="header">
      <el-input v-model="queryParams.keyword" placeholder="关键字查询:自动识别字段" class="small" style="width: 300px;margin-right: 10px" @keyup.enter.native="getList">
        <template v-slot:append>
          <el-button icon="el-icon-search" @click="getList" />
        </template>
      </el-input>

      <el-button v-if="hasPermission('del_log')" class="small" type="danger" icon="el-icon-delete" :disabled="multiple" @click="handleDelete">删除</el-button>
      <el-button v-if="hasPermission('reset_log')" class="small" type="danger" icon="el-icon-delete" @click="handleClean">清空</el-button>
    </div>

    <!-- 表格 -->
    <el-table v-loading="loading" border :data="logList" style="margin-top: 20px" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="45" align="center" />
      <el-table-column show-overflow-tooltip label="编号" width="160" align="center" prop="id" />
      <el-table-column show-overflow-tooltip label="用户" align="center" prop="username" width="150" />
      <el-table-column show-overflow-tooltip label="地址" align="center" prop="ip_addr" width="150" />
      <el-table-column show-overflow-tooltip label="请求路径" align="center" prop="path" width="150" />
      <el-table-column show-overflow-tooltip label="请求方式" align="center" prop="method" width="150" />
      <el-table-column show-overflow-tooltip label="请求状态" align="center" prop="status" width="150" />
      <el-table-column show-overflow-tooltip label="发起时间" align="center" prop="start_at" width="200">
        <template v-slot="scope">
          <div>{{ formatTimestamp(scope.row.start_at) }}</div>
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip label="耗时(ms)" align="center" prop="duration" width="150" />
      <el-table-column show-overflow-tooltip label="系统" align="center" prop="client_os" width="150" />
      <el-table-column show-overflow-tooltip label="浏览器" align="center" prop="browser" width="199" />
    </el-table>

    <!-- 分页 -->
    <pagination v-show="total>0" :total="total" :page.sync="queryParams.page" :limit.sync="queryParams.perPage" @pagination="getList" />

  </div>
</template>

<script>
import { listLog, deleteLog, emptyLog } from '@/api/system/log'
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
      // 表单参数
      form: {},
      // 表格数据
      logList: null,
      // 总条数
      total: 0,
      // 查询参数
      queryParams: {
        page: 1,
        perPage: 10,
        keyword: undefined
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
    /* 处理返回时间的时间戳转化*/
    formatTimestamp(timestamp) { // 转化时间戳
      return formatTimestamp(timestamp)
    },
    /* 处理selection的行选择*/
    handleSelectionChange(selection) {
      this.ids = selection.map((item) => item.id)
      this.single = selection.length !== 1
      this.multiple = !selection.length
    },
    /** 列表查询 */
    getList() {
      this.loading = true
      listLog(this.queryParams).then(
        (res) => {
          this.logList = res.items
          this.total = res.total
          this.loading = false
        }
      )
    },
    /* 处理删除*/
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
        .then(function() {
          return deleteLog(ids)
        })
        .then(() => {
          this.getList()
          this.msgSuccess('删除成功')
        })
        .catch(function() { })
    },
    /** 清空按钮操作 */
    handleClean() {
      this.$confirm('是否确认清空所有登录日志数据项?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(function() {
        return emptyLog()
      }).then(res => {
        this.getList()
        this.msgSuccess('清空成功')
      }).catch(function() { })
    }
  }
}
</script>
