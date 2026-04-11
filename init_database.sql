-- 创建数据库
CREATE DATABASE IF NOT EXISTS distribution_app 
DEFAULT CHARACTER SET utf8mb4 
DEFAULT COLLATE utf8mb4_unicode_ci;

USE distribution_app;

-- 1. roles 表（角色表）
CREATE TABLE `roles` (
  `id` tinyint unsigned NOT NULL AUTO_INCREMENT COMMENT '角色唯一ID（自增，TINYINT足够容纳所有角色，轻量化）',
  `role_name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称（唯一，如：收货人/送货人/管理员）',
  `description` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '角色描述（可选，如：仅可查看自身订单和配送轨迹）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '角色创建时间（自动填充）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_role_name` (`role_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表（用户角色定义，用于权限管控）';

-- 2. permissions 表（权限表）
CREATE TABLE `permissions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '权限唯一ID（自增）',
  `perm_code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限编码（唯一，如：user:address:list）',
  `perm_name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限名称（如：查询用户地址）',
  `api_path` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '对应接口路径（如：/api/user/addresses）',
  `method` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求方式（GET/POST/PUT/DELETE）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '权限创建时间（自动填充）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_perm_code` (`perm_code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表（接口/操作权限定义）';

-- 3. role_permissions 表（角色-权限关联表）
CREATE TABLE `role_permissions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '关联记录唯一ID（自增）',
  `role_id` tinyint unsigned NOT NULL COMMENT '关联角色ID（逻辑外键，关联roles表id）',
  `perm_id` int unsigned NOT NULL COMMENT '关联权限ID（逻辑外键，关联permissions表id）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '绑定时间（自动填充）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_role_perm` (`role_id`,`perm_id`) USING BTREE,
  KEY `idx_role_id` (`role_id`) USING BTREE,
  KEY `idx_perm_id` (`perm_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色-权限关联表（多对多绑定，实现角色权限配置）';

-- 4. users 表（用户表）
CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '用户唯一ID（自增）',
  `username` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名（登录凭证，可使用手机号）',
  `password` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户密码（SHA256加密后存储，不可明文）',
  `role_id` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '关联角色ID（逻辑外键，关联roles表id；默认2-收货人）',
  `is_del` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否删除：0-未删除 1-已删除（逻辑删除）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '用户创建时间（自动填充）',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '用户信息更新时间（自动更新）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_username` (`username`) USING BTREE,
  KEY `idx_role_id` (`role_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表（用户认证与角色关联，实时物流追踪基础表）';

-- 5. addresses 表（收货地址表）
CREATE TABLE `addresses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '地址唯一ID（自增，无符号避免负数）',
  `user_id` int unsigned NOT NULL COMMENT '归属用户ID（收货人ID，关联用户表user_id）',
  `province` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '省份（如：北京市）',
  `city` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '市（如：北京市）',
  `district` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '区/县（允许为空，如：朝阳区）',
  `street` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '街道/乡镇（允许为空，如：望京街道）',
  `detail` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '详细地址（门牌号/小区/楼栋，如：XX小区1栋101）',
  `receiver` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收件人姓名',
  `phone` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收件人电话（支持手机号/固话，如：13800138000、010-12345678）',
  `latitude` decimal(10,8) NOT NULL COMMENT '纬度（GCJ-02高德/百度坐标系，如：39.90882300）',
  `longitude` decimal(11,8) NOT NULL COMMENT '经度（GCJ-02坐标系，如：116.39747000）',
  `is_default` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否默认地址：0-否，1-是',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（自动填充当前时间）',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间（自动更新）',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_lat_lng` (`latitude`,`longitude`) USING BTREE,
  KEY `idx_user_default` (`user_id`,`is_default`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收货地址表（实时物流追踪专用）';

-- 6. orders 表（配送订单表）
CREATE TABLE `orders` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '订单唯一ID（自增）',
  `user_id` int unsigned NOT NULL COMMENT '关联收货人ID（逻辑外键，关联users表id）',
  `address_id` int unsigned NOT NULL COMMENT '关联收货地址ID（逻辑外键，关联addresses表id）',
  `order_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单编号（唯一，用于前端展示/查询）',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '订单状态：1-待配送 / 2-配送中 / 3-已完成 / 4-已取消',
  `goods_info` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '货物信息（如：生鲜食品 2件）',
  `delivery_user_id` int NOT NULL COMMENT '配送员ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '订单创建时间（自动填充）',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '订单更新时间（自动更新）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_order_no` (`order_no`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE,
  KEY `idx_user_status` (`user_id`,`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配送订单表（实时物流追踪核心业务表）';

-- 7. delivery_assign 表（配送任务分配表）
CREATE TABLE `delivery_assign` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '分配记录唯一ID（自增）',
  `order_id` int unsigned NOT NULL COMMENT '关联订单ID（逻辑外键，关联orders表id）',
  `delivery_user_id` int unsigned NOT NULL COMMENT '关联配送员ID（逻辑外键，关联users表id）',
  `assign_user_id` int unsigned NOT NULL COMMENT '分配人ID（逻辑外键，关联users表id）',
  `assign_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间（自动填充当前时间）',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_order_id` (`order_id`) USING BTREE,
  KEY `idx_delivery_user_id` (`delivery_user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配送任务分配表（记录订单分配给配送员的信息）';

-- 8. location_traces 表（送货人位置轨迹表）
CREATE TABLE `location_traces` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '轨迹记录唯一ID（自增）',
  `order_id` int unsigned NOT NULL COMMENT '关联订单ID（逻辑外键，关联orders表id）',
  `delivery_user_id` int unsigned NOT NULL COMMENT '关联送货人ID（逻辑外键，关联users表id）',
  `longitude` decimal(11,8) NOT NULL COMMENT '经度（GCJ-02高德坐标系）',
  `latitude` decimal(10,8) NOT NULL COMMENT '纬度（GCJ-02高德坐标系）',
  `upload_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '位置上传时间（自动填充当前时间）',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_order_upload` (`order_id`,`upload_time`) USING BTREE,
  KEY `idx_delivery_user` (`delivery_user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='送货人位置轨迹表（实时物流追踪）';

-- 插入初始角色数据
INSERT INTO `roles` (`id`, `role_name`, `description`, `create_time`) VALUES
(1, '收货人', '可查看自身订单和配送轨迹，确认签收等权限', NOW()),
(2, '配送员', '可接受任务、实时上报位置、确认取件和送达等权限', NOW()),
(3, '管理员', '拥有系统所有管理权限', NOW());

-- 插入初始权限数据（示例）
INSERT INTO `permissions` (`perm_code`, `perm_name`, `api_path`, `method`, `create_time`) VALUES
('user:address:list', '查询用户地址列表', '/api/settings/address/list', 'GET', NOW()),
('user:address:detail', '查询地址详情', '/api/settings/address', 'GET', NOW()),
('user:address:create', '创建地址', '/api/settings/address/create', 'POST', NOW()),
('user:address:update', '更新地址', '/api/settings/address', 'PUT', NOW()),
('user:address:delete', '删除地址', '/api/settings/address', 'DELETE', NOW()),
('user:address:default', '设置默认地址', '/api/settings/address/default', 'PUT', NOW()),
('user:logout', '退出登录', '/api/settings/logout', 'POST', NOW()),
('order:list', '查询订单列表', '/api/orders/list', 'GET', NOW()),
('order:create', '创建订单', '/api/orders/create', 'POST', NOW()),
('order:detail', '查询订单详情', '/api/orders/detail', 'GET', NOW()),
('delivery:task:list', '查询配送任务列表', '/api/delivery/tasks', 'GET', NOW()),
('delivery:task:accept', '接受配送任务', '/api/delivery/task/accept', 'POST', NOW()),
('delivery:location:upload', '上传位置信息', '/api/delivery/location', 'POST', NOW()),
('delivery:task:complete', '完成配送任务', '/api/delivery/task/complete', 'POST', NOW());

-- 为角色分配权限（示例）
-- 收货人权限
INSERT INTO `role_permissions` (`role_id`, `perm_id`, `create_time`) 
SELECT 1, id, NOW() FROM permissions WHERE perm_code LIKE 'user:%' OR perm_code LIKE 'order:%';

-- 配送员权限
INSERT INTO `role_permissions` (`role_id`, `perm_id`, `create_time`) 
SELECT 2, id, NOW() FROM permissions WHERE perm_code LIKE 'delivery:%';

-- 管理员拥有所有权限
INSERT INTO `role_permissions` (`role_id`, `perm_id`, `create_time`) 
SELECT 3, id, NOW() FROM permissions;
