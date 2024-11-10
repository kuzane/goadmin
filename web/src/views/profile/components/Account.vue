<template>
  <el-form ref="form" :model="data" :rules="rules" label-width="80px">
    <el-form-item label="用户手机" prop="phone">
      <el-input v-model="data.phone" maxlength="11" />
    </el-form-item>
    <el-form-item label="用户邮箱" prop="email">
      <el-input v-model="data.email" maxlength="50" />
    </el-form-item>
    <el-form-item label="用户描述" prop="description">
      <el-input v-model="data.description" />
    </el-form-item>
    <el-form-item>
      <el-button :loading="loading" type="primary" size="mini" @click="submit">保存</el-button>
      <el-button type="danger" size="mini" @click="close">关闭</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import { changeProfile } from '@/api/system/user'

export default {
  props: {
    user: {
      type: Object,
      default: () => {
        return {
          phone: '',
          email: '',
          description: ''
        }
      }
    }
  },
  data() {
    return {
      data: {},
      loading: false,
      rules: {
        email: [
          { required: true, message: '用户邮箱不能为空', trigger: 'blur' },
          {
            type: 'email',
            message: "'请输入正确的邮箱地址",
            trigger: ['blur', 'change']
          }
        ],
        phone: [
          { required: true, message: '用户手机不能为空', trigger: 'blur' },
          {
            pattern: /^1[3|4|5|6|7|8|9][0-9]\d{8}$/,
            message: '请输入正确的手机号码',
            trigger: 'blur'
          }
        ],
        description: [
          { required: true, message: '用户描述不能为空', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getData()
  },
  methods: {
    getData() {
      this.data = {
        email: this.user.email,
        phone: this.user.phone,
        description: this.user.description
      }
    },
    submit() {
      this.$refs['form'].validate(valid => {
        if (valid) {
          this.loading = true
          changeProfile(this.data).then(() => {
            this.msgSuccess('修改成功')
            this.$store.dispatch('user/getInfo')
            this.$router.push({ path: '/' })
          })
          this.loading = false
        }
      })
    },
    close() {
      this.$store.dispatch('tagsView/delView', this.$route)
      this.$router.push({ path: '/' })
    }
  }
}
</script>
