/*
 Navicat Premium Data Transfer

 Source Server         : localmysql
 Source Server Type    : MySQL
 Source Server Version : 80300
 Source Host           : 192.168.3.10:3306
 Source Schema         : goadmin

 Target Server Type    : MySQL
 Target Server Version : 80300
 File Encoding         : 65001

 Date: 10/11/2024 11:37:31
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for endpoint
-- ----------------------------
DROP TABLE IF EXISTS `endpoint`;
CREATE TABLE `endpoint` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` bigint DEFAULT NULL,
  `updated_at` bigint DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `method` varchar(255) DEFAULT NULL,
  `module` longtext,
  `kind` varchar(255) DEFAULT NULL COMMENT '类别(系统模块下的子类别)',
  `identity` varchar(64) DEFAULT NULL COMMENT '接口的唯一标识,做为和Roles的外键关联',
  `remark` varchar(255) DEFAULT NULL COMMENT '描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_path_method` (`path`,`method`),
  UNIQUE KEY `idx_identity` (`identity`)
) ENGINE=InnoDB AUTO_INCREMENT=159 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of endpoint
-- ----------------------------
BEGIN;
INSERT INTO `endpoint` VALUES (16, 1730982877, 1731209826, '/sys/users', 'POST', '系统管理', '用户管理', 'add_user', '增加用户');
INSERT INTO `endpoint` VALUES (17, 1730982877, 1731209826, '/sys/users/:id', 'DELETE', '系统管理', '用户管理', 'del_user', '删除用户');
INSERT INTO `endpoint` VALUES (18, 1730982877, 1731209826, '/sys/users/:id', 'PATCH', '系统管理', '用户管理', 'upd_user', '修改用户');
INSERT INTO `endpoint` VALUES (19, 1730982877, 1731209826, '/sys/users/:id', 'GET', '系统管理', '用户管理', 'get_user', '用户详情');
INSERT INTO `endpoint` VALUES (20, 1730982877, 1731209826, '/sys/users', 'GET', '系统管理', '用户管理', 'list_user', '用户列表');
INSERT INTO `endpoint` VALUES (21, 1730982877, 1731209826, '/sys/roles', 'POST', '系统管理', '角色管理', 'add_role', '增加角色');
INSERT INTO `endpoint` VALUES (22, 1730982877, 1731209826, '/sys/roles/:id', 'DELETE', '系统管理', '角色管理', 'del_role', '删除角色');
INSERT INTO `endpoint` VALUES (23, 1730982877, 1731209826, '/sys/roles/:id', 'PATCH', '系统管理', '角色管理', 'upd_role', '修改角色');
INSERT INTO `endpoint` VALUES (24, 1730982877, 1731209826, '/sys/roles/:id', 'GET', '系统管理', '角色管理', 'get_role', '角色详情');
INSERT INTO `endpoint` VALUES (25, 1730982877, 1731209826, '/sys/roles', 'GET', '系统管理', '角色管理', 'list_role', '角色列表');
INSERT INTO `endpoint` VALUES (26, 1730982877, 1731209826, '/sys/roles/apis', 'GET', '系统管理', '角色管理', 'tree_api', '接口树');
INSERT INTO `endpoint` VALUES (27, 1730982877, 1731209826, '/sys/logs/:id', 'DELETE', '系统管理', '日志管理', 'del_log', '删除日志');
INSERT INTO `endpoint` VALUES (28, 1730982877, 1731209826, '/sys/logs', 'GET', '系统管理', '日志管理', 'list_log', '日志列表');
INSERT INTO `endpoint` VALUES (29, 1730982877, 1731209826, '/sys/logs', 'DELETE', '系统管理', '日志管理', 'reset_log', '清空日志');
COMMIT;

-- ----------------------------
-- Table structure for role_endpoints
-- ----------------------------
DROP TABLE IF EXISTS `role_endpoints`;
CREATE TABLE `role_endpoints` (
  `role_rolename` varchar(20) NOT NULL,
  `endpoint_identity` varchar(64) NOT NULL COMMENT '接口的唯一标识,做为和Roles的外键关联',
  PRIMARY KEY (`role_rolename`,`endpoint_identity`),
  KEY `fk_role_endpoints_endpoint` (`endpoint_identity`),
  CONSTRAINT `fk_role_endpoints_endpoint` FOREIGN KEY (`endpoint_identity`) REFERENCES `endpoint` (`identity`),
  CONSTRAINT `fk_role_endpoints_role` FOREIGN KEY (`role_rolename`) REFERENCES `roles` (`rolename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of role_endpoints
-- ----------------------------
BEGIN;
INSERT INTO `role_endpoints` VALUES ('admin', 'add_role');
INSERT INTO `role_endpoints` VALUES ('role-editor', 'add_role');
INSERT INTO `role_endpoints` VALUES ('admin', 'add_user');
INSERT INTO `role_endpoints` VALUES ('usereditor', 'add_user');
INSERT INTO `role_endpoints` VALUES ('admin', 'del_log');
INSERT INTO `role_endpoints` VALUES ('log-editor', 'del_log');
INSERT INTO `role_endpoints` VALUES ('admin', 'del_role');
INSERT INTO `role_endpoints` VALUES ('admin', 'del_user');
INSERT INTO `role_endpoints` VALUES ('admin', 'get_role');
INSERT INTO `role_endpoints` VALUES ('role-editor', 'get_role');
INSERT INTO `role_endpoints` VALUES ('admin', 'get_user');
INSERT INTO `role_endpoints` VALUES ('usereditor', 'get_user');
INSERT INTO `role_endpoints` VALUES ('admin', 'list_log');
INSERT INTO `role_endpoints` VALUES ('log-editor', 'list_log');
INSERT INTO `role_endpoints` VALUES ('log-view', 'list_log');
INSERT INTO `role_endpoints` VALUES ('admin', 'list_role');
INSERT INTO `role_endpoints` VALUES ('role-editor', 'list_role');
INSERT INTO `role_endpoints` VALUES ('usereditor', 'list_role');
INSERT INTO `role_endpoints` VALUES ('admin', 'list_user');
INSERT INTO `role_endpoints` VALUES ('role-editor', 'list_user');
INSERT INTO `role_endpoints` VALUES ('usereditor', 'list_user');
INSERT INTO `role_endpoints` VALUES ('admin', 'reset_log');
INSERT INTO `role_endpoints` VALUES ('log-editor', 'reset_log');
INSERT INTO `role_endpoints` VALUES ('admin', 'tree_api');
INSERT INTO `role_endpoints` VALUES ('role-editor', 'tree_api');
INSERT INTO `role_endpoints` VALUES ('admin', 'upd_role');
INSERT INTO `role_endpoints` VALUES ('role-editor', 'upd_role');
INSERT INTO `role_endpoints` VALUES ('admin', 'upd_user');
INSERT INTO `role_endpoints` VALUES ('usereditor', 'upd_user');
COMMIT;

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` bigint DEFAULT NULL,
  `updated_at` bigint DEFAULT NULL,
  `rolename` varchar(20) NOT NULL,
  `nickname` varchar(20) NOT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `parents` longtext COMMENT '当前角色的父级角色,用于继承权限',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_roles_rolename` (`rolename`),
  UNIQUE KEY `uni_roles_nickname` (`nickname`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of roles
-- ----------------------------
BEGIN;
INSERT INTO `roles` VALUES (1, 1730723535, 1731204918, 'admin', '管理员', 1, '管理员权限', NULL);
INSERT INTO `roles` VALUES (14, 1731199969, 1731204356, 'role-editor', '角色编辑者', 1, '角色编辑者', NULL);
INSERT INTO `roles` VALUES (31, 1731202873, 1731202886, 'log-editor', '日志管理者', 1, '日志管理者', NULL);
INSERT INTO `roles` VALUES (32, 1731204408, 1731204408, 'log-view', '日志查看', 1, '具备查看日志列表', NULL);
INSERT INTO `roles` VALUES (37, 1731204897, 1731204897, 'usereditor', '用户维护者', 1, '用户维护者', NULL);
COMMIT;

-- ----------------------------
-- Table structure for server_configs
-- ----------------------------
DROP TABLE IF EXISTS `server_configs`;
CREATE TABLE `server_configs` (
  `skey` varchar(191) NOT NULL,
  `value` longtext,
  PRIMARY KEY (`skey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of server_configs
-- ----------------------------
BEGIN;
INSERT INTO `server_configs` VALUES ('email-captcha', '\n<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>验证码邮件</title>\n    <style>\n        body {\n            font-family: Arial, sans-serif;\n            background-color: #f4f4f4;\n            padding: 20px;\n            margin: 0;\n        }\n        .container {\n            background-color: #ffffff;\n            padding: 20px;\n            border-radius: 8px;\n            box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);\n            width: 100%;\n            max-width: 600px;\n            margin: 0 auto;\n        }\n        h2 {\n            color: #333333;\n        }\n        .code {\n            font-size: 32px;\n            font-weight: bold;\n            color: #4CAF50;\n            margin: 20px 0;\n        }\n        .note {\n            font-size: 14px;\n            color: #888888;\n        }\n        .footer {\n            font-size: 12px;\n            color: #888888;\n            margin-top: 20px;\n            text-align: center;\n        }\n        .footer a {\n            color: #4CAF50;\n            text-decoration: none;\n        }\n    </style>\n</head>\n<body>\n    <div class=\"container\">\n        <h2>您好，{{.Name}}！</h2>\n        <p>感谢您请求验证码。请使用以下验证码重置您的账户密码：</p>\n        \n        <div class=\"code\">{{.Content}}</div>\n        \n        <p class=\"note\">该验证码有效期为 5 分钟，请尽快使用。</p>\n\n        <p class=\"note\">如果您没有请求此验证码，请忽略此邮件。</p>\n        \n        <div class=\"footer\">\n            <p><a rel=\"noopener\" href=\"{{.Domain}}\" target=\"_blank\">访问goAdmin</a></p>\n        </div>\n    </div>\n</body>\n</html>\n');
INSERT INTO `server_configs` VALUES ('email-password', '\n<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>您的账号密码</title>\n    <style>\n        body {\n            font-family: Arial, sans-serif;\n            background-color: #f4f4f4;\n            padding: 20px;\n            margin: 0;\n        }\n        .container {\n            background-color: #ffffff;\n            padding: 20px;\n            border-radius: 8px;\n            box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);\n            width: 100%;\n            max-width: 600px;\n            margin: 0 auto;\n        }\n        h2 {\n            color: #333333;\n        }\n        .password {\n            font-size: 24px;\n            font-weight: bold;\n            color: #ff5733;\n            background-color: #f8f8f8;\n            padding: 10px;\n            border-radius: 5px;\n            margin: 20px 0;\n            word-wrap: break-word;\n        }\n        .note {\n            font-size: 14px;\n            color: #888888;\n        }\n        .footer {\n            font-size: 12px;\n            color: #888888;\n            margin-top: 20px;\n            text-align: center;\n        }\n        .footer a {\n            color: #4CAF50;\n            text-decoration: none;\n        }\n    </style>\n</head>\n<body>\n    <div class=\"container\">\n        <h2>您好, {{.Name}}！</h2>\n        <p>您的账号密码为：</p>\n        \n        <div class=\"password\">{{.Content}}</div>\n        \n        <p class=\"note\">为了确保您的账户安全，强烈建议您尽快修改密码。您可以在登录后访问“个人中心”进行修改。</p>\n\n        <p class=\"note\">如果您没有请求此邮件，请忽略此邮件。</p>\n        \n        <div class=\"footer\">\n            <p><a rel=\"noopener\" href=\"{{.Domain}}\" target=\"_blank\">访问goAdmin</a></p>\n        </div>\n    </div>\n</body>\n</html>\n');
INSERT INTO `server_configs` VALUES ('jwt-secret', 'X5PVCRIH4L6MRWA4ONQ6HFF2GQQQKIA2XAWVJCALZRRSCU3ROWCQ====');
COMMIT;

-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
  `role_rolename` varchar(20) NOT NULL,
  `user_username` varchar(50) NOT NULL COMMENT '用户名',
  PRIMARY KEY (`role_rolename`,`user_username`),
  KEY `fk_user_roles_user` (`user_username`),
  CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_rolename`) REFERENCES `roles` (`rolename`),
  CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_username`) REFERENCES `users` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
BEGIN;
INSERT INTO `user_roles` VALUES ('admin', 'admin');
INSERT INTO `user_roles` VALUES ('log-view', 'guest');
INSERT INTO `user_roles` VALUES ('log-editor', 'logadmin');
INSERT INTO `user_roles` VALUES ('role-editor', 'roleadmin');
INSERT INTO `user_roles` VALUES ('usereditor', 'usereditor');
COMMIT;

-- ----------------------------
-- Table structure for userlogs
-- ----------------------------
DROP TABLE IF EXISTS `userlogs`;
CREATE TABLE `userlogs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` bigint DEFAULT NULL,
  `updated_at` bigint DEFAULT NULL,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `ip_addr` varchar(128) NOT NULL COMMENT 'ip地址',
  `start_at` varchar(18) NOT NULL COMMENT '请求开始时间',
  `path` varchar(128) NOT NULL COMMENT '请求路径',
  `method` varchar(50) NOT NULL COMMENT '请求方法',
  `status` int NOT NULL COMMENT '请求状态',
  `duration` int NOT NULL COMMENT '请求耗时(ms)',
  `browser` varchar(128) DEFAULT NULL COMMENT '客户端',
  `client_os` varchar(128) DEFAULT NULL COMMENT '客户端操作系统',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of userlogs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` bigint DEFAULT NULL,
  `updated_at` bigint DEFAULT NULL,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '用户密码',
  `nickname` varchar(50) DEFAULT NULL COMMENT '中文名',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `phone` varchar(15) NOT NULL COMMENT '手机号',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `access_token` varchar(255) DEFAULT NULL,
  `refresh_token` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_username` (`username`),
  UNIQUE KEY `uni_users_email` (`email`),
  UNIQUE KEY `uni_users_phone` (`phone`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO `users` VALUES (1, 1730723525, 1731209785, 'admin', '$2a$10$d0IxldJ1cdCGnMz4VrkxaOmOH3KFObabpAgwClO4TI5m3S5jQph7a', '超级管理员', 'admin@qq.com', '18888888888', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', 1, '测试账户', '', '');
INSERT INTO `users` VALUES (2, 1731204566, 1731204583, 'logadmin', '$2a$10$d0IxldJ1cdCGnMz4VrkxaOmOH3KFObabpAgwClO4TI5m3S5jQph7a', '日志管理员', 'logadmin@qq.com', '19999999999', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', 1, '日志管理', '', '');
INSERT INTO `users` VALUES (3, 1731204813, 1731204827, 'roleadmin', '$2a$10$d0IxldJ1cdCGnMz4VrkxaOmOH3KFObabpAgwClO4TI5m3S5jQph7a.', '角色管理员', 'roleadmin@qq.com', '13333333333', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', 1, '角色管理', '0', '');
INSERT INTO `users` VALUES (4, 1731204972, 1731205019, 'usereditor', '$2a$10$d0IxldJ1cdCGnMz4VrkxaOmOH3KFObabpAgwClO4TI5m3S5jQph7a', '用户管理员', 'usereditor@qq.com', '15555555555', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', 1, '用户管理权限', '', '');
INSERT INTO `users` VALUES (5, 1731112995, 1731204474, 'guest', '$2a$10$d0IxldJ1cdCGnMz4VrkxaOmOH3KFObabpAgwClO4TI5m3S5jQph7a', '访客用户', 'guest@qq.com', '15888888888', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', 1, '访客用户', '', '');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
