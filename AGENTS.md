# Agents 模块设计

## 1. 概念定义
主要的业务内容是物流app，例如顺丰app，具体是用户可以寄快递：从地址簿选择寄件地址和输入收件地址、商品信息，数据库的orders表加入这条数据。通过收件人的手机号匹配到收件人。配送员端通过查询orders表查询status为1(待接单)的订单来接单，通过开启所选的订单的配送任务将订单的状态改为配送中，并将数据加入到delivery_assign表。通过开启goroutine用websocket和map使一个配送员可对应多个用户，并实时将自己的位置信息上传至location_traces表，用户查看订单详情可以显示配送的路径和配送员的实时位置。订单详情显示配送的信息，包括位置信息，位置长时间不动上报异常。配送员送达后将订单状态改为已完成。
### 1.1 Agent 角色定义
- **配送员**：直接执行物流任务的人员，负责从取件到送达的全流程操作 role_id=2
- **用户/收货人**：收货人 role_id=1

### 1.2 核心术语
- **任务**：需要配送员执行的物流订单
- **节点**：物流过程中的关键时间点和位置

### 1.3 模块定位（区分收货人和配送员的功能模块）
- **配送员模块**：负责配送员的注册、登录、信息管理和任务执行
- **收货人模块**：负责收货人修改个人信心、添加订单、查看订单状态，确认签收等收货人的基本权限

## 2. API 接口设计

### 2.1 认证相关接口

#### 用户注册

**路径** : `/api/auth/register`  
**方法** : `POST`  
**功能描述** : 用户注册，创建新用户账号

**请求体 (JSON)** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| username | string | 是 | 用户名（手机号） |
| password | string | 是 | 用户密码 |
| role_id | int | 否 | 角色ID（默认2-收货人） |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 用户登录

**路径** : `/api/auth/login`  
**方法** : `POST`  
**功能描述** : 用户登录，获取JWT token

**请求体 (JSON)** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| username | string | 是 | 用户名（手机号） |
| password | string | 是 | 用户密码 |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "13800138000",
      "role_id": 1
    }
  }
}
```

### 2.2 收货人信息管理接口

#### 获取地址列表

**路径** : `/api/settings/address/list`  
**方法** : `GET`  
**功能描述** : 获取当前用户的所有收货地址

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "province": "北京市",
      "city": "北京市",
      "district": "朝阳区",
      "street": "望京街道",
      "detail": "XX小区1栋101",
      "receiver": "张三",
      "phone": "13800138000",
      "latitude": 39.908823,
      "longitude": 116.397470,
      "is_default": 1
    }
  ]
}
```

#### 获取地址详情

**路径** : `/api/settings/address`  
**方法** : `GET`  
**功能描述** : 获取指定地址的详细信息

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| id | int | 是 | 地址ID |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "province": "北京市",
    "city": "北京市",
    "district": "朝阳区",
    "street": "望京街道",
    "detail": "XX小区1栋101",
    "receiver": "张三",
    "phone": "13800138000",
    "latitude": 39.908823,
    "longitude": 116.397470,
    "is_default": 1
  }
}
```

#### 新建地址

**路径** : `/api/settings/address/create`  
**方法** : `POST`  
**功能描述** : 创建新的收货地址

**请求体 (JSON)** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| province | string | 是 | 省份 |
| city | string | 是 | 城市 |
| district | string | 否 | 区县 |
| street | string | 否 | 街道 |
| detail | string | 是 | 详细地址 |
| receiver | string | 是 | 收件人姓名 |
| phone | string | 是 | 收件人电话 |
| is_default | int | 否 | 是否默认地址（0-否，1-是） |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 修改地址

**路径** : `/api/settings/address`  
**方法** : `PUT`  
**功能描述** : 修改收货地址信息

**请求体 (JSON)** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| id | int | 是 | 地址ID |
| province | string | 是 | 省份 |
| city | string | 是 | 城市 |
| district | string | 否 | 区县 |
| street | string | 否 | 街道 |
| detail | string | 是 | 详细地址 |
| receiver | string | 是 | 收件人姓名 |
| phone | string | 是 | 收件人电话 |
| is_default | int | 否 | 是否默认地址（0-否，1-是） |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 删除地址

**路径** : `/api/settings/address`  
**方法** : `DELETE`  
**功能描述** : 删除指定的收货地址

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| id | int | 是 | 地址ID |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 设置默认地址

**路径** : `/api/settings/address/default`  
**方法** : `PUT`  
**功能描述** : 设置指定地址为默认地址

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| id | int | 是 | 地址ID |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 退出登录

**路径** : `/api/settings/logout`  
**方法** : `POST`  
**功能描述** : 用户退出登录

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

### 2.3 订单管理接口

#### 获取订单列表

**路径** : `/api/orders/list`  
**方法** : `GET`  
**功能描述** : 获取当前用户的订单列表

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| current_page | int | 否 | 当前页码（默认1） |
| per_page | int | 否 | 每页数量（默认10） |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "order_no": "ORD123456789",
        "status": 1,
        "goods_info": "生鲜食品 2件",
        "receiver_name": "李四",
        "receiver_phone": "13900139000",
        "create_time": "2024-01-01 10:00:00"
      }
    ],
    "total": 10,
    "current_page": 1,
    "per_page": 10
  }
}
```

#### 创建订单

**路径** : `/api/orders/create`  
**方法** : `POST`  
**功能描述** : 创建新的配送订单

**请求体 (JSON)** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| sender_address_id | int | 是 | 寄件人地址ID |
| goods_info | string | 否 | 货物信息 |
| receiver_name | string | 是 | 收货人姓名 |
| receiver_phone | string | 是 | 收货人电话 |
| receiver_province | string | 是 | 收货人省份 |
| receiver_city | string | 是 | 收货人城市 |
| receiver_district | string | 否 | 收货人区县 |
| receiver_street | string | 否 | 收货人街道 |
| receiver_detail | string | 是 | 收货人详细地址 |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 获取订单详情

**路径** : `/api/orders/detail`  
**方法** : `GET`  
**功能描述** : 获取订单详细信息

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| order_no | string | 是 | 订单编号 |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "order_no": "ORD123456789",
    "status": 1,
    "goods_info": "生鲜食品 2件",
    "sender_province": "北京市",
    "sender_city": "北京市",
    "sender_district": "朝阳区",
    "sender_detail": "XX小区1栋101",
    "sender_receiver": "张三",
    "sender_phone": "13800138000",
    "receiver_name": "李四",
    "receiver_phone": "13900139000",
    "receiver_province": "上海市",
    "receiver_city": "上海市",
    "receiver_district": "浦东新区",
    "receiver_detail": "YY小区2栋202",
    "delivery_user_id": null,
    "create_time": "2024-01-01 10:00:00",
    "update_time": "2024-01-01 10:00:00"
  }
}
```

### 2.4 配送员任务管理接口

#### 获取待接单订单列表

**路径** : `/api/rider/waiting`  
**方法** : `GET`  
**功能描述** : 配送员获取待接单的订单列表

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| current_page | int | 否 | 当前页码（默认1） |
| per_page | int | 否 | 每页数量（默认10） |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "order_no": "ORD123456789",
        "status": 1,
        "goods_info": "生鲜食品 2件",
        "sender_detail": "XX小区1栋101",
        "receiver_detail": "YY小区2栋202"
      }
    ],
    "total": 5,
    "current_page": 1,
    "per_page": 10
  }
}
```

#### 接单

**路径** : `/api/rider/accept`  
**方法** : `POST`  
**功能描述** : 配送员接单，更新订单状态为已接单

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| order_no | string | 是 | 订单编号 |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 取件

**路径** : `/api/rider/pickup`  
**方法** : `POST`  
**功能描述** : 配送员取件，更新订单状态为已取件待配送

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| order_no | string | 是 | 订单编号 |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 开始配送

**路径** : `/api/rider/start`  
**方法** : `POST`  
**功能描述** : 配送员开始配送，更新订单状态为配送中

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| order_no | string | 是 | 订单编号 |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 完成订单

**路径** : `/api/rider/complete`  
**方法** : `POST`  
**功能描述** : 配送员完成配送，更新订单状态为已完成

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| order_no | string | 是 | 订单编号 |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

### 2.5 位置管理接口

#### 上传位置

**路径** : `/api/rider/location`  
**方法** : `POST`  
**功能描述** : 配送员上传实时位置信息

**请求体 (JSON)** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| order_id | int | 是 | 订单ID |
| longitude | float | 是 | 经度 |
| latitude | float | 是 | 纬度 |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 获取配送轨迹

**路径** : `/api/orders/trace`  
**方法** : `GET`  
**功能描述** : 获取订单的配送轨迹

**查询参数** : 
| 参数名 | 类型 | 必填 | 描述 | 
|--------|------|------|------| 
| order_no | string | 是 | 订单编号 |

**成功响应 (200 OK)** : 
```json
{
  "code": 200,
  "msg": "success",
  "data": [
    {
      "longitude": 116.397470,
      "latitude": 39.908823,
      "upload_time": "2024-01-01 10:00:00"
    },
    {
      "longitude": 116.400000,
      "latitude": 39.910000,
      "upload_time": "2024-01-01 10:05:00"
    }
  ]
}
```
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
3. 更新任务状态为「已接单」
4. 开始配送，实时上报位置
5. 到达目的地，更新任务状态为「已完成」
7. 更新任务状态为「已完成」
8. 系统更新物流订单状态
9. 通知收货人

#### 3.1.3 收货人操作流程
1. 收货人登录系统
2. 修改个人信息（如手机号、地址）
3. 添加新订单（包括订单详情、配送地址、联系电话等）
4. 查看订单状态（包括取件状态、配送状态、签收状态等）

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

create table orders(
    id                 int unsigned auto_increment comment '订单唯一ID（自增）'primary key,
    sender_user_id     int unsigned                               not null comment '关联寄件人ID（逻辑外键，关联users表id）',
    sender_address_id  int unsigned                               not null comment '关联寄件人地址ID（逻辑外键，关联addresses表id',
    order_no           varchar(32)                                not null comment '订单编号（唯一，用于前端展示/查询）',
    status             tinyint unsigned default '1'               not null comment '订单状态：1-待接单/2-已接单待取件/ 3-配送中 / 4-已完成 / 5-已取消 /6-已取件待配送',
    goods_info         varchar(255)                               null comment '货物信息（如：生鲜食品 2件）',
    receiver_name      varchar(32)                                not null comment '收货人姓名',
    receiver_phone     varchar(16)                                not null comment '收货人电话',
    receiver_province  varchar(20)                                not null comment '收货人省份',
    receiver_city      varchar(20)                                not null comment '收货人城市',
    receiver_district  varchar(20)                                null comment '收货人区县',
    receiver_street    varchar(50)                                null comment '收货人街道',
    receiver_detail    varchar(128)                               not null comment '收货人详细地址',
    receiver_latitude  decimal(10, 8)                             not null comment '收货人地址纬度',
    receiver_longitude decimal(11, 8)                             not null comment '收货人地址经度',
    delivery_user_id   int unsigned                               null comment '关联配送员ID（逻辑外键，关联users表id）',
    create_time        datetime         default CURRENT_TIMESTAMP not null comment '订单创建时间（自动填充）',
    update_time        datetime         default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '订单更新时间（自动填充，每次更新时自动更新）',
    constraint uk_order_no
        unique (order_no)
)
    comment '配送订单表（实时物流追踪核心业务表，包含寄件人和收件人信息）';
#### 5.1.6 addresses 表

CREATE TABLE `addresses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '地址唯一ID（自增，无符号避免负数）',
  `user_id` int unsigned NOT NULL NOT NULL COMMENT '归属用户ID（收货人ID，关联用户表user_id）',
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
