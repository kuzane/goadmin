<template>
  <el-card style="margin-bottom:20px;">
    <div slot="header" class="clearfix">
      <span>个人信息</span>
    </div>

    <div class="user-profile">
      <div class="box-avatar">
        <pan-thumb :image="user.avatar" :height="'100px'" :width="'100px'" :hoverable="false">
          <div>Hello</div>
          {{ user.role }}
        </pan-thumb>
      </div>
      <div class="box-center">
        <ul class="list-group list-group-striped">
          <li class="list-group-item">
            <svg-icon icon-class="user" style="margin-right: 8px;color: #606266;" /> 用户名称
            <div class="pull-right">{{ user.name }}</div>
          </li>
          <li class="list-group-item">
            <svg-icon icon-class="phone" style="margin-right: 8px;color: #606266;" /> 手机号码
            <div class="pull-right">{{ user.phone }}</div>
          </li>
          <li class="list-group-item">
            <svg-icon icon-class="email" style="margin-right: 8px;color: #606266;" /> 用户邮箱
            <div class="pull-right">{{ user.email }}</div>
          </li>
          <li class="list-group-item">
            <svg-icon icon-class="peoples" style="margin-right: 8px;color: #606266;" /> 用户角色
            <div class="pull-right">{{ user.role }}</div>
          </li>

          <li class="list-group-item">
            <svg-icon icon-class="textarea" style="margin-right: 8px;color: #606266;" /> 用户描述
            <div class="pull-right">{{ user.description }}</div>
          </li>
        </ul>
      </div>
    </div>
  </el-card>
</template>

<script>
import PanThumb from '@/components/PanThumb'

export default {
  components: { PanThumb },
  props: {
    user: {
      type: Object,
      default: () => {
        return {
          name: '',
          avatar: '',
          phone: '',
          email: '',
          nickname: '',
          description: '',
          role: ''
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.box-avatar {
  margin: 0 auto;
  display: table;
}
.box-center {
  margin: 0 auto;
  .user-name {
    font-weight: bold;
    font-size: 18px;
    margin-bottom: 4px;
    text-align: center;
  }

  .user-role {
    font-size: 14px;
    color: #909399;
    text-align: center;
  }

  .list-group {
    list-style: none;
    padding: 0;
    margin: 0;

    .list-group-item {
      display: flex;
      align-items: center;
      padding: 12px 0;
      font-size: 14px;
      border-bottom: 1px solid #f0f0f0;
      .pull-right {
        margin-left: auto;
        color: #333;
      }
    }
  }
}

.text-muted {
  color: #777;
}

.user-profile {
  .user-name {
    font-weight: bold;
  }

  .box-center {
    padding-top: 10px;
  }

  .user-role {
    padding-top: 10px;
    font-weight: 400;
    font-size: 14px;
  }

  .box-social {
    padding-top: 30px;

    .el-table {
      border-top: 1px solid #dfe6ec;
    }
  }

  .user-follow {
    padding-top: 20px;
  }
}

.user-bio {
  margin-top: 20px;
  color: #606266;

  span {
    padding-left: 4px;
  }

  .user-bio-section {
    font-size: 14px;
    padding: 15px 0;

    .user-bio-section-header {
      border-bottom: 1px solid #dfe6ec;
      padding-bottom: 10px;
      margin-bottom: 10px;
      font-weight: bold;
    }
  }
}
</style>
