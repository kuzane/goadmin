<template>
  <div class="login-container">
    <el-form ref="form" :model="form" :rules="loginRules" class="login-form" autocomplete="on" label-position="left">
      <div class="title-container">
        <h3 class="title">使用邮箱重置密码</h3>
      </div>

      <div class="email-container">
        <el-form-item prop="email" class="email-input">
          <span class="svg-container">
            <svg-icon icon-class="email" />
          </span>
          <el-input ref="email" v-model="form.email" placeholder="请输入邮箱" name="email" type="text" tabindex="1" autocomplete="on" />
        </el-form-item>
        <el-button class="btn-getcaptcha" :disabled="captchaCd" @click="handleCaptcha('form')">
          <span v-if="!captchaCd">点击获取验证码</span>
          <span v-else>{{ waitingTime }}s后重新获取</span>
        </el-button>
      </div>

      <el-form-item prop="captcha">
        <span class="svg-container">
          <svg-icon icon-class="number" />
        </span>
        <el-input ref="captcha" v-model="form.captcha" placeholder="请输入验证码" name="captcha" />
      </el-form-item>

      <el-button :loading="loading" type="primary" style="width: 100%; margin-bottom: 30px" @click.native.prevent="handleSubmit">提交</el-button>
    </el-form>
  </div>
</template>

<script>
import { forgotPwd, getEmailCaptcha } from "@/api/system/user"
export default {
  name: 'Login',
  data() {
    const validateEmail = (rule, value, callback) => {
      if (!value) {
        callback(new Error("邮箱不能为空"));
      }
      // 使用正则表达式进行验证邮箱验证
      if (
        !/^([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z]{2,4})+$/.test(
          value
        )
      ) {
        callback(new Error("邮箱格式不正确"));
      }
      // 自定义校验规则 需要调用callback()函数！
      callback();
    };
    return {
      form: {
        email: '',
        captcha: ''
      },
      waitingTime: 60,
      captchaCd: false,
      loading: false,
      loginRules: {
        email: [{ required: true, validator: validateEmail, trigger: "blur" }],
        captcha: [
          { required: true, message: "验证码不能为空", trigger: "blur" },
        ]
      },
    }
  },
  methods: {
    handleCaptcha(form) {
      this.$refs[form].validateField("email", async (valid) => {
        if (!valid) {
          getEmailCaptcha({ email: this.form.email }).then(() => {
            this.msgSuccess("请求成功");
          })
          this.captchaCd = true
          const timer = setInterval(() => {
            this.waitingTime--;
            if (this.waitingTime <= 0) {
              this.waitingTime = 60;
              this.captchaCd = false;
              clearInterval(timer);
            }
          }, 1000);
        } else {
          return false;
        }
      });
    },
    handleSubmit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.loading = true
          this.form.captcha = String(this.form.captcha)
          forgotPwd(this.form).then(() => {
            this.msgSuccess("请求成功")
            this.$router.push({ path: '/login' })
          })
          this.loading = false
        } else {
          return false
        }
      })
    }

  }
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */

$bg: #283443;
$light_gray: #fff;
$cursor: #fff;

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container .el-input input {
    color: $cursor;
  }
}

/* reset element-ui css */
.login-container {
  .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;

    input {
      background: transparent;
      border: 0px;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: $light_gray;
      height: 47px;
      caret-color: $cursor;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px $bg inset !important;
        -webkit-text-fill-color: $cursor !important;
      }
    }
  }

  .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.1);
    border-radius: 5px;
    color: #454545;
  }

  .btn-getcaptcha {
    border: 0px;
    border-radius: 0px;
    padding: 12px 5px 12px 15px;
    color: $light_gray;
    height: 47px;
    caret-color: $cursor;
  }
}
</style>

<style lang="scss" scoped>
$bg: #2d3a4b;
$dark_gray: #889aa4;
$light_gray: #eee;

.login-container {
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  overflow: hidden;

  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;
  }

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .email-container {
    display: flex;
    .email-input {
      flex: 1;
      border-right: none;
      border-top-right-radius: 0; // 去掉左上圆角
      border-bottom-right-radius: 0; // 去掉左下圆角
    }
    .btn-getcaptcha {
      width: 150px;
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 5px;
      border-top-left-radius: 0; // 去掉左上圆角
      border-bottom-left-radius: 0; // 去掉左下圆角
      height: 51px; // 确保按钮与输入框高度一致
      background: rgba(0, 0, 0, 0.1);
      color: $light_gray;
      box-shadow: none;
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: $light_gray;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }

  @media only screen and (max-width: 470px) {
    .thirdparty-button {
      display: none;
    }
  }
}
</style>
