# Agents 模块设计

## 1. 概念定义

### 1.1 Agent 角色定义
- **配送员**：直接执行物流任务的人员，负责从取件到送达的全流程操作
- **收货人**：收货人

### 1.2 核心术语
- **任务**：需要配送员执行的物流订单
- **节点**：物流过程中的关键时间点和位置

### 1.3 模块定位（区分收货人和配送员的功能模块）
- **配送员模块**：负责配送员的注册、登录、信息管理和任务执行
- **收货人模块**：负责收货人修改个人信心、添加订单、查看订单状态，确认签收等收货人的基本权限

## 2. API 接口设计

### 2.1 认证相关接口

| 方法 | 路径 | 功能描述 | 请求体 (JSON) | 成功响应 (200 OK) |
|------|------|----------|---------------|-------------------|


### 2.2 收货人信息管理接口

| 方法 | 路径 | 功能描述 | 请求体 (JSON) | 成功响应 (200 OK) |
|------|------|----------|---------------|-------------------|
Get|/api/settings/address/list|地址列表
Get|/api/settings/address/:id|地址详情
Post|/api/settings/address/create|新建地址
Put|/api/seetings/address|修改地址
Delete|/api/settings/address|删除地址
Put|/api/settings/address/:id/default|设置默认地址
Post|/api/settings/logout|退出登陆

### 2.3 任务管理接口

| 方法 | 路径 | 功能描述 | 请求体 (JSON) | 成功响应 (200 OK) |
|------|------|----------|---------------|-------------------|


### 2.4 位置管理接口

| 方法 | 路径 | 功能描述 | 请求体 (JSON) | 成功响应 (200 OK) |
|------|------|----------|---------------|-------------------|


### 2.5 订单管理接口

| 方法 | 路径 | 功能描述 | 请求体 (JSON) | 成功响应 (200 OK) |
|------|------|----------|---------------|-------------------|

## 3. 业务逻辑

### 3.1 核心业务流程

#### 3.1.1 配送员登录流程

1. 接收登录请求（手机号和密码）
2. 生成JWT token
3. 更新最后登录时间
4. 返回token和配送员信息

#### 3.1.2 任务执行流程
1. 配送员接受任务
2. 前往发货点取件
3. 输入取件码，更新任务状态为「已取件」
4. 开始配送，实时上报位置
5. 到达收货点，输入签收码，更新任务状态为「已签收」
6. 获取收货人签名（可选）
7. 更新任务状态为「已完成」
8. 系统更新物流订单状态
9. 通知收货人

#### 3.1.3 收货人操作流程
1. 收货人登录系统
2. 修改个人信息（如手机号、地址）
3. 添加新订单（包括订单详情、配送地址、联系电话等）
4. 查看订单状态（包括取件状态、配送状态、签收状态等）
5. 确认签收订单（输入签收码，更新订单状态为「已签收」）

#### 3.1.4 异常处理流程
1. 配送员遇到异常情况（如交通堵塞、找不到地址）
2. 上报异常信息，包括类型、原因和预计影响
3. 系统记录异常信息
4. 任务调度器评估异常影响：
   - 轻微异常：提醒配送员注意，调整预计时间
   - 严重异常：重新分配任务给其他配送员
5. 通知相关用户（发货人、收货人）

## 4. 权限管理

### 4.1 RBAC 权限设计

#### 4.1.1 角色定义

- **配送员**：基础执行角色，拥有执行任务的基本权限
- **收货人**：修改个人信心、添加订单、查看订单状态，确认签收等收货人的基本权限

## 5. 数据库设计

### 5.1 表结构设计

#### 5.1.1 users 表

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

#### 5.1.2 roles 表

CREATE TABLE `roles` (
  `id` tinyint unsigned NOT NULL AUTO_INCREMENT COMMENT '角色唯一ID（自增，TINYINT足够容纳所有角色，轻量化）',
  `role_name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称（唯一，如：收货人/送货人/管理员）',
  `description` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '角色描述（可选，如：仅可查看自身订单和配送轨迹）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '角色创建时间（自动填充）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_role_name` (`role_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表（用户角色定义，用于权限管控）';

#### 5.1.3 permissions 表

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

#### 5.1.4 role_permissions 表

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

#### 5.1.5 orders 表

CREATE TABLE `orders` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '订单唯一ID（自增）',
  `user_id` int unsigned NOT NULL COMMENT '关联收货人ID（逻辑外键，关联users表id）',
  `address_id` int unsigned NOT NULL COMMENT '关联收货地址ID（逻辑外键，关联addresses表id）',
  `order_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单编号（唯一，用于前端展示/查询）',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '订单状态：1-待配送 / 2-配送中 / 3-已完成 / 4-已取消',
  `goods_info` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '货物信息（如：生鲜食品 2件）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '订单创建时间（自动填充）',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '订单更新时间（自动更新）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_order_no` (`order_no`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE,
  KEY `idx_user_status` (`user_id`,`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配送订单表（实时物流追踪核心业务表）';

#### 5.1.6 addresses 表

CREATE TABLE `addresses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '地址唯一ID（自增，无符号避免负数）',
  `user_id` varchar(8) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '归属用户ID（收货人ID，关联用户表user_id）',
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

#### 5.1.7 delivery_assign 表

CREATE TABLE `delivery_assign` (
  `id` int NOT NULL,
  `order_id` int NOT NULL,
  `delivery_user_id` int NOT NULL,
  `assign_user_id` int NOT NULL,
  `assign_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

#### 5.1.8 location_traces 

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
