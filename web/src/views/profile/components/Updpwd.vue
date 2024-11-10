<template>
  <el-form ref="form" :model="data" :rules="rules" label-width="80px">
    <el-form-item label="旧密码" prop="old_password">
      <el-input v-model="data.old_password" placeholder="请输入旧密码" type="password" />
    </el-form-item>
    <el-form-item label="新密码" prop="new_password">
      <el-input v-model="data.new_password" placeholder="请输入新密码" type="password" />
    </el-form-item>
    <el-form-item label="确认密码" prop="confirm_password">
      <el-input v-model="data.confirm_password" placeholder="请确认密码" type="password" />
    </el-form-item>
    <el-form-item>
      <el-button :loading="loading" type="primary" size="mini" @click="submit">保存</el-button>
      <el-button type="danger" size="mini" @click="close">关闭</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import { changePassword } from '@/api/system/user'

export default {
  data() {
    const equalToPassword = (rule, value, callback) => {
      if (this.data.new_password !== value) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }
    return {
      loading: false,
      data: {
        old_password: undefined,
        new_password: undefined,
        confirm_password: undefined
      },
      // 表单校验
      rules: {
        old_password: [
          { required: true, message: '旧密码不能为空', trigger: 'blur' }
        ],
        new_password: [
          { required: true, message: '新密码不能为空', trigger: 'blur' },
          { validator: this.validatePasswordComplexity, trigger: 'blur' }
        ],
        confirm_password: [
          { required: true, message: '确认密码不能为空', trigger: 'blur' },
          { required: true, validator: equalToPassword, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    submit() {
      this.$refs['form'].validate((valid) => {
        if (valid) {
          this.loading = true
          changePassword(this.data).then(() => {
            this.msgSuccess('修改成功')
            this.$store.dispatch('user/logout').then(() => {
              location.reload()
            })
            this.loading = false
          })
        }
      })
    },
    validatePasswordComplexity(rule, value, callback) {
      // 正则表达式检查密码复杂度
      const regex =
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{6,}$/
      if (regex.test(value)) {
        callback()
      } else {
        callback(
          new Error('密码必须至少6位且同时包含大小写字母、数字和[@$!%*?&]等特殊字符中的一位')
        )
      }
    },
    close() {
      this.$store.dispatch('tagsView/delView', this.$route)
      this.$router.push({ path: '/' })
    }
  }
}
</script>
