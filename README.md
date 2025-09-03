# 1.项目特色
- 1.实现用户的双token,由后端实现,续期的ak通过响应报文header发送给前端
- 2.使用内置函数自动生成接口(权限)表，所有的权限围绕接口(权限)进行
- 3.用户管理角色，角色关联接口(权限)，实现的RBAC，角色权限可以继承
- 4.前端菜单和相关按钮也是直接关联到接口(权限)
- 5.只要授予角色有关的list_xxx权限，那么此角色就会有相应的菜单栏
- 6.授予角色相关的add_xxx,del_xxx那么就会出现相应新增和删除按
- 7.系统对admin用户进行了单独授权，admin为超级管理员

前端项目基于：https://github.com/PanJiaChen/vue-element-admin

# 2.快速体验
- 所有账户密码均为: 123456
```sh
cd deploy && docker-compose up -d
```

# 3.构建部署

## 3.1.docker本地构建坚信
- 构建容器示例命令：
```sh
docker build . -t goadmin:latest
```

## 3.2.github ci构建
默认提交变更会将本项目生成两个镜像推送到docker hub上,其一为镜像tag为commit id，其二镜像tag为latest，拉取镜像为
```sh
docker pull kuzane/goadmin:latest
```

## 3.4.k8s部署
- 利用构建的镜像进行goadmin部署，默认镜像使用的`kuzane/goadmin:latest` 
```sh
 kustomize build kustomize/ |kubectl apply -f -
```


# 4.项目开发
需要在项目目录下创建一个 `.env`文件，配置下必要的参数
```env
SERVER_PORT=8000
SERVER_HOST=http://localhost:8000
DATABASE_DRIVER=mysql
DATABASE_DATASOURCE=root:123456@tcp(localhost:3306)/goadmin?parseTime=true
#如果需要开启邮件功能,需要添加如下配置
SMTP_DSN=smtp_user:smtp_pass:smtp_host:smtp_port:true
```
- 前后端不分离启动项目

需先构建前端 `cd web && npm install && npm run build:prod`，并在`.env`中开启`LOAD_WEB`
```env
LOAD_WEB=true
```



