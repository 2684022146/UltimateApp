USE distribution_app;

-- 添加测试用户 yzx，密码为123456（MD5加密）
INSERT INTO `users` (`username`, `password`, `role_id`, `is_del`, `create_time`, `update_time`) 
VALUES (
    'yzx',
    'e10adc3949ba59abbe56e057f20f883e', -- MD5('123456')
    1, -- 角色ID：1-收货人
    0, -- 未删除
    NOW(),
    NOW()
);

-- 为测试用户添加一个默认地址
INSERT INTO `addresses` (
    `user_id`, `province`, `city`, `district`, `street`, `detail`, 
    `receiver`, `phone`, `latitude`, `longitude`, `is_default`, `create_time`, `update_time`
) VALUES (
    (SELECT id FROM users WHERE username = 'yzx'),
    '北京市',
    '北京市',
    '朝阳区',
    '望京街道',
    '望京SOHO T1 C座 2801',
    '张三',
    '13800138000',
    39.998823,
    116.487470,
    1, -- 默认地址
    NOW(),
    NOW()
);

-- 查看用户信息
SELECT u.id, u.username, u.role_id, r.role_name 
FROM users u 
JOIN roles r ON u.role_id = r.id 
WHERE u.username = 'yzx';

-- 查看用户地址
SELECT * FROM addresses WHERE user_id = (SELECT id FROM users WHERE username = 'yzx');
