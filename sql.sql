/*
 Navicat Premium Dump SQL

 Source Server         : Cloud MySQL
 Source Server Type    : MySQL
 Source Server Version : 50744 (5.7.44)
 Source Host           : 1.95.0.111:3306
 Source Schema         : schisandra-cloud-album

 Target Server Type    : MySQL
 Target Server Version : 50744 (5.7.44)
 File Encoding         : 65001

 Date: 12/11/2024 22:47:39
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sca_auth_permission
-- ----------------------------
DROP TABLE IF EXISTS `sca_auth_permission`;
CREATE TABLE `sca_auth_permission`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `permission_name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '权限名称',
  `parent_id` bigint(20) NULL DEFAULT NULL COMMENT '父ID',
  `type` tinyint(4) NULL DEFAULT 0 COMMENT '类型 0 菜单 1 接口 ',
  `path` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '路径',
  `method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '请求方式',
  `status` tinyint(4) NULL DEFAULT 0 COMMENT '状态 0 启用 1 停用',
  `icon` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '图标',
  `permission_key` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '权限关键字',
  `order` int(11) NULL DEFAULT NULL COMMENT '排序',
  `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注 描述',
  `created_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `update_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '权限表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_auth_permission
-- ----------------------------
INSERT INTO `sca_auth_permission` VALUES (1, 'test', 0, NULL, 'test', 'get', NULL, 'test', 'test', 0, '2024-09-02 22:50:11', '2024-09-02 22:50:11', NULL, 'test', NULL, NULL);
INSERT INTO `sca_auth_permission` VALUES (2, 'test1', 0, NULL, 'test1', 'get', NULL, 'test1', 'test1', 0, '2024-09-02 22:51:22', '2024-09-02 22:51:22', NULL, 'test1', NULL, NULL);
INSERT INTO `sca_auth_permission` VALUES (3, 'test2', 0, NULL, 'test2', 'post', NULL, 'test2', 'test2', 0, '2024-09-02 22:51:22', '2024-09-02 22:51:22', NULL, 'test2', NULL, NULL);

-- ----------------------------
-- Table structure for sca_auth_permission_rule
-- ----------------------------
DROP TABLE IF EXISTS `sca_auth_permission_rule`;
CREATE TABLE `sca_auth_permission_rule`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sca_auth_casbin_rule`(`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) USING BTREE,
  UNIQUE INDEX `idx_sca_auth_permission_rule`(`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 36 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '角色/权限/用户关系表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_auth_permission_rule
-- ----------------------------
INSERT INTO `sca_auth_permission_rule` VALUES (35, 'g', '607492814274629', 'user', '', '', '', '');
INSERT INTO `sca_auth_permission_rule` VALUES (28, 'p', 'user', '/api/auth/comment/cancel_like', 'POST', '', NULL, NULL);
INSERT INTO `sca_auth_permission_rule` VALUES (27, 'p', 'user', '/api/auth/comment/like', 'POST', NULL, NULL, NULL);
INSERT INTO `sca_auth_permission_rule` VALUES (24, 'p', 'user', '/api/auth/comment/list', 'POST', NULL, NULL, NULL);
INSERT INTO `sca_auth_permission_rule` VALUES (22, 'p', 'user', '/api/auth/comment/submit', 'POST', NULL, NULL, NULL);
INSERT INTO `sca_auth_permission_rule` VALUES (1, 'p', 'user', '/api/auth/permission/get_user_permissions', 'POST', '', '', '');
INSERT INTO `sca_auth_permission_rule` VALUES (25, 'p', 'user', '/api/auth/reply/list', 'POST', NULL, NULL, NULL);
INSERT INTO `sca_auth_permission_rule` VALUES (26, 'p', 'user', '/api/auth/reply/reply/submit', 'POST', NULL, NULL, NULL);
INSERT INTO `sca_auth_permission_rule` VALUES (23, 'p', 'user', '/api/auth/reply/submit', 'POST', NULL, NULL, NULL);
INSERT INTO `sca_auth_permission_rule` VALUES (30, 'p', 'user', '/api/captcha/slide/generate', 'GET', NULL, NULL, NULL);

-- ----------------------------
-- Table structure for sca_auth_role
-- ----------------------------
DROP TABLE IF EXISTS `sca_auth_role`;
CREATE TABLE `sca_auth_role`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_name` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '角色名称',
  `role_key` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '角色关键字',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除 0 未删除 1已删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_auth_role
-- ----------------------------
INSERT INTO `sca_auth_role` VALUES (1, '超级管理员', 'root', '2024-08-13 16:58:21', '2024-08-22 15:15:44', 0);
INSERT INTO `sca_auth_role` VALUES (2, '管理员', 'admin', '2024-08-13 16:58:34', '2024-08-13 16:58:34', 0);
INSERT INTO `sca_auth_role` VALUES (3, '普通用户', 'user', '2024-08-13 16:59:00', '2024-08-13 16:59:00', 0);

-- ----------------------------
-- Table structure for sca_auth_user
-- ----------------------------
DROP TABLE IF EXISTS `sca_auth_user`;
CREATE TABLE `sca_auth_user`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `uid` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '唯一ID',
  `username` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '用户名',
  `nickname` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '昵称',
  `email` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '电话',
  `password` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '密码',
  `gender` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '性别',
  `avatar` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '头像',
  `status` tinyint(4) NULL DEFAULT 0 COMMENT '状态 0 正常 1 封禁',
  `introduce` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '介绍',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除 0 未删除 1 已删除',
  `blog` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '博客',
  `location` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '地址',
  `company` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '公司',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uid`(`uid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 51 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_auth_user
-- ----------------------------

-- ----------------------------
-- Table structure for sca_auth_user_device
-- ----------------------------
DROP TABLE IF EXISTS `sca_auth_user_device`;
CREATE TABLE `sca_auth_user_device`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户ID',
  `ip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '登录IP',
  `location` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '地址',
  `agent` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '设备信息',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除',
  `browser` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '浏览器',
  `operating_system` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '操作系统',
  `browser_version` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '浏览器版本',
  `mobile` int(11) NULL DEFAULT NULL COMMENT '是否为手机 0否1是',
  `bot` int(11) NULL DEFAULT NULL COMMENT '是否为bot 0否1是',
  `mozilla` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '火狐版本',
  `platform` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '平台',
  `engine_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '引擎名称',
  `engine_version` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '引擎版本',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 39 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户设备信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_auth_user_device
-- ----------------------------

-- ----------------------------
-- Table structure for sca_auth_user_social
-- ----------------------------
DROP TABLE IF EXISTS `sca_auth_user_social`;
CREATE TABLE `sca_auth_user_social`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户ID',
  `open_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '第三方用户的 open id',
  `source` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '第三方用户来源',
  `status` int(11) NULL DEFAULT 0 COMMENT '状态 0正常 1 封禁',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 23 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '社会用户信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_auth_user_social
-- ----------------------------

-- ----------------------------
-- Table structure for sca_comment_likes
-- ----------------------------
DROP TABLE IF EXISTS `sca_comment_likes`;
CREATE TABLE `sca_comment_likes`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `topic_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '话题ID',
  `user_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户ID',
  `comment_id` bigint(20) NOT NULL COMMENT '评论ID',
  `like_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `comment_id`(`comment_id`) USING BTREE,
  INDEX `topic_id`(`topic_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 81 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '评论点赞表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_comment_likes
-- ----------------------------

-- ----------------------------
-- Table structure for sca_comment_message
-- ----------------------------
DROP TABLE IF EXISTS `sca_comment_message`;
CREATE TABLE `sca_comment_message`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `topic_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '话题Id',
  `from_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '来自人',
  `to_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '送达人',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '消息内容',
  `is_read` int(11) NULL DEFAULT NULL COMMENT '是否已读',
  `created_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新人',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除 0否 1是',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `id`(`id`) USING BTREE,
  INDEX `topic_id`(`topic_id`) USING BTREE,
  INDEX `from_id`(`from_id`) USING BTREE,
  INDEX `to_id`(`to_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '消息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_comment_message
-- ----------------------------

-- ----------------------------
-- Table structure for sca_comment_reply
-- ----------------------------
DROP TABLE IF EXISTS `sca_comment_reply`;
CREATE TABLE `sca_comment_reply`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '评论用户id',
  `topic_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '评论话题id',
  `topic_type` int(11) NULL DEFAULT NULL COMMENT '话题类型',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '评论内容',
  `comment_type` int(11) NULL DEFAULT NULL COMMENT '评论类型 0评论 1 回复',
  `reply_to` bigint(20) NULL DEFAULT NULL COMMENT '回复子评论ID',
  `reply_id` bigint(20) NULL DEFAULT NULL COMMENT '回复父评论Id',
  `reply_user` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '回复人id',
  `author` int(11) NULL DEFAULT 0 COMMENT '评论回复是否作者  0否 1是',
  `likes` bigint(20) NULL DEFAULT 0 COMMENT '点赞数',
  `reply_count` bigint(20) NULL DEFAULT 0 COMMENT '回复数量',
  `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除 0未删除 1 已删除',
  `created_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `update_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新人',
  `browser` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '浏览器',
  `operating_system` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '操作系统',
  `comment_ip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'IP地址',
  `location` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '地址',
  `agent` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '设备信息',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `comment_id`(`id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `reply_id`(`reply_id`) USING BTREE,
  INDEX `topic_id`(`topic_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 146 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '评论表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_comment_reply
-- ----------------------------
INSERT INTO `sca_comment_reply` VALUES (142, '607492814274629', '123', 0, 'waht?', 0, 0, 0, '', 1, 0, 1, '2024-11-04 17:54:50', '2024-11-04 18:33:34', 0, '607492814274629', '', 'Chrome', 'Windows 10', '127.0.0.1', '内网IP|内网IP', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36');
INSERT INTO `sca_comment_reply` VALUES (143, '607492814274629', '123', 0, '什么<img width=\"30px\" height=\"30px\" loading=\"lazy\" src=\"/emoji/qq/gif/106.gif\" alt=\"emoji 106.gif\"/>', 1, 0, 142, '607492814274629', 1, 0, 0, '2024-11-04 17:56:17', '2024-11-04 17:56:17', 0, '607492814274629', '', 'Chrome', 'Windows 10', '127.0.0.1', '内网IP|内网IP', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36');
INSERT INTO `sca_comment_reply` VALUES (144, '607492814274629', '123', 0, '哈哈哈哈哈哈哈哈', 0, 0, 0, '', 1, 0, 0, '2024-11-04 18:42:24', '2024-11-05 16:38:36', 0, '607492814274629', '', 'Chrome', 'Windows 10', '127.0.0.1', '内网IP|内网IP', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36');
INSERT INTO `sca_comment_reply` VALUES (145, '607492814274629', '123', 0, '图像测试<img width=\"30px\" height=\"30px\" loading=\"lazy\" src=\"/emoji/qq/gif/1.gif\" alt=\"emoji 1.gif\"/><img width=\"30px\" height=\"30px\" loading=\"lazy\" src=\"/emoji/qq/gif/1.gif\" alt=\"emoji 1.gif\"/>', 0, 0, 0, '', 1, 0, 0, '2024-11-04 19:03:15', '2024-11-04 19:03:15', 0, '607492814274629', '', 'Chrome', 'Windows 10', '127.0.0.1', '内网IP|内网IP', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36');

-- ----------------------------
-- Table structure for sca_file_folder
-- ----------------------------
DROP TABLE IF EXISTS `sca_file_folder`;
CREATE TABLE `sca_file_folder`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `folder_name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '文件夹名称',
  `parent_folder_id` bigint(20) NULL DEFAULT NULL COMMENT '父文件夹编号',
  `folder_addr` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '文件夹名称',
  `user_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户编号',
  `folder_source` int(11) NULL DEFAULT NULL COMMENT '文件夹来源 0相册 1 评论',
  `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新人',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除 0 未删除 1 已删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文件夹信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_file_folder
-- ----------------------------

-- ----------------------------
-- Table structure for sca_file_info
-- ----------------------------
DROP TABLE IF EXISTS `sca_file_info`;
CREATE TABLE `sca_file_info`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `file_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '文件名',
  `file_size` double NULL DEFAULT NULL COMMENT '文件大小',
  `file_type_id` bigint(20) NULL DEFAULT NULL COMMENT '文件类型编号',
  `upload_time` datetime NULL DEFAULT NULL COMMENT '上传时间',
  `folder_id` bigint(20) NULL DEFAULT NULL COMMENT '文件夹编号',
  `user_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户编号',
  `file_source` int(11) NULL DEFAULT NULL COMMENT '文件来源 0 相册 1 评论',
  `status` int(11) NULL DEFAULT NULL COMMENT '文件状态',
  `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新人',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除 0 未删除 1 已删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文件信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_file_info
-- ----------------------------

-- ----------------------------
-- Table structure for sca_file_recycle_bin
-- ----------------------------
DROP TABLE IF EXISTS `sca_file_recycle_bin`;
CREATE TABLE `sca_file_recycle_bin`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `file_id` bigint(20) NULL DEFAULT NULL COMMENT '文件编号',
  `folder_id` bigint(20) NULL DEFAULT NULL COMMENT '文件夹编号',
  `type` int(11) NULL DEFAULT NULL COMMENT '类型 0 文件 1 文件夹',
  `user_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户编号',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  `original_path` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '原始路径',
  `deleted` int(11) NULL DEFAULT NULL COMMENT '是否被永久删除 0否 1是',
  `file_source` int(11) NULL DEFAULT NULL COMMENT '文件来源 0 相册 1 评论',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文件回收站表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_file_recycle_bin
-- ----------------------------

-- ----------------------------
-- Table structure for sca_file_type
-- ----------------------------
DROP TABLE IF EXISTS `sca_file_type`;
CREATE TABLE `sca_file_type`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `type_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '类型名称',
  `mime_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'MIME 类型',
  `status` int(11) NULL DEFAULT NULL COMMENT '类型状态',
  `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新人',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除 0 未删除 1 已删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文件类型表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_file_type
-- ----------------------------

-- ----------------------------
-- Table structure for sca_message_report
-- ----------------------------
DROP TABLE IF EXISTS `sca_message_report`;
CREATE TABLE `sca_message_report`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户Id',
  `type` int(11) NULL DEFAULT NULL COMMENT '举报类型 0评论 1 相册',
  `comment_id` bigint(20) NULL DEFAULT NULL COMMENT '评论Id',
  `topic_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '话题Id',
  `report_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '举报',
  `report_content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '举报说明内容',
  `report_tag` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '举报标签',
  `status` int(11) NULL DEFAULT NULL COMMENT '状态（0 未处理 1 已处理）',
  `created_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_by` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更新人',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` int(11) NULL DEFAULT 0 COMMENT '是否删除 0否 1是',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '举报信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_message_report
-- ----------------------------

-- ----------------------------
-- Table structure for sca_user_follows
-- ----------------------------
DROP TABLE IF EXISTS `sca_user_follows`;
CREATE TABLE `sca_user_follows`  (
  `follower_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '关注者',
  `followee_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '被关注者',
  `status` tinyint(1) NULL DEFAULT 0 COMMENT '关注状态（0 未互关 1 互关）',
  `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`follower_id`, `followee_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户关注表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_user_follows
-- ----------------------------

-- ----------------------------
-- Table structure for sca_user_level
-- ----------------------------
DROP TABLE IF EXISTS `sca_user_level`;
CREATE TABLE `sca_user_level`  (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `user_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户Id',
  `level_type` tinyint(1) NULL DEFAULT NULL COMMENT '等级类型',
  `level` int(11) NULL DEFAULT NULL COMMENT '等级',
  `level_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '等级名称',
  `exp_start` bigint(20) NULL DEFAULT NULL COMMENT '开始经验值',
  `exp_end` bigint(20) NULL DEFAULT NULL COMMENT '结束经验值',
  `level_description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '等级描述',
  `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户等级表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sca_user_level
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
