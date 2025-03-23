USE chatroom;

-- 群类型
CREATE TABLE `groups` (
                        id INT AUTO_INCREMENT PRIMARY KEY,
                        name VARCHAR(50) NOT NULL, -- 群名称
                        description VARCHAR(255), -- 群描述
                        avatar VARCHAR(255) -- 群头像
);

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
                                     id BIGINT AUTO_INCREMENT PRIMARY KEY, -- 用户 id（自增整数）
                                     email VARCHAR(255) NOT NULL UNIQUE, -- 邮箱
                                     password VARCHAR(255) NOT NULL, -- 密码
                                     create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);
-- 创建复合索引
CREATE INDEX idx_users_email_id_password_create_at ON users (email, id, password, create_at);


CREATE TABLE IF NOT EXISTS accounts (
                                        id BIGINT PRIMARY KEY,                                -- 账号 id
                                        user_id BIGINT NOT NULL,                              -- 用户 id（外键）
                                        name VARCHAR(255) NOT NULL,                           -- 账号名
                                        avatar VARCHAR(255) NOT NULL,                         -- 账号头像
                                        gender ENUM('男', '女', '未知') NOT NULL DEFAULT '未知', -- 账号性别
                                        signature VARCHAR(1024) NOT NULL DEFAULT '这个用户很懒，什么也没有留下~', -- 账号签名
                                        create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
                                        CONSTRAINT account_unique_name UNIQUE (user_id, name), -- 一个用户的不同账号名不能重复
                                        FOREIGN KEY (user_id) REFERENCES users(id)
                                            ON DELETE CASCADE
                                            ON UPDATE CASCADE  -- 外键约束：删除和更新时级联操作
);


-- 创建账号名和头像索引
CREATE INDEX account_index_name_avatar ON accounts(name, avatar);




-- 创建统一关系表（包含群组和好友的所有字段）
CREATE TABLE relations (
                           id BIGINT AUTO_INCREMENT PRIMARY KEY,
    -- 关系类型标识
                           relation_type ENUM('group', 'friend') NOT NULL,

    -- 群组相关字段（当 relation_type='group' 时有效）
                           group_name VARCHAR(50),
                           group_description VARCHAR(255),
                           group_avatar VARCHAR(255),

    -- 好友相关字段（当 relation_type='friend' 时有效）
                           account1_id BIGINT,
                           account2_id BIGINT,

                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- 添加外键约束（假设有 users 表）
#                            FOREIGN KEY (account1_id) REFERENCES accounts(id),
#                            FOREIGN KEY (account2_id) REFERENCES accounts(id),

    -- 约束条件（需 MySQL 8.0+）
                           CHECK (
                               (relation_type = 'group' AND
                                group_name IS NOT NULL AND
                                account1_id IS NULL AND
                                account2_id IS NULL)
                                   OR
                               (relation_type = 'friend' AND
                                account1_id IS NOT NULL AND
                                account2_id IS NOT NULL AND
                                group_name IS NULL AND
                                group_description IS NULL AND
                                group_avatar IS NULL)
                               )

    -- 确保好友关系不重复（如 1-2 和 2-1 视为相同）
#                            CHECK (
#                                relation_type <> 'friend' OR
#                                (account1_id < account2_id)
#                                ),
#                            UNIQUE KEY (account1_id, account2_id)
);


-- 创建消息通知类型表（示例：可以根据需要将此列放入消息通知表中）
CREATE TABLE IF NOT EXISTS msg_notifications (
                                                 id BIGINT AUTO_INCREMENT PRIMARY KEY, -- 消息通知 id
                                                 msg_type ENUM('system', 'common') NOT NULL, -- 消息通知类型
                                                 content TEXT NOT NULL, -- 消息内容
                                                 create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);


# -- 账号对群组或好友关系的设置
CREATE TABLE IF NOT EXISTS settings (
                                        account_id BIGINT NOT NULL, -- 账号id（外键）
                                        relation_id BIGINT NOT NULL, -- 关系 id（外键）
                                        nick_name VARCHAR(255) NOT NULL, -- 昵称，默认是账户名或群组名
                                        is_not_disturb BOOLEAN NOT NULL DEFAULT FALSE, -- 是否免打扰
                                        is_pin BOOLEAN NOT NULL DEFAULT FALSE, -- 是否置顶
                                        pin_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 置顶时间
                                        is_show BOOLEAN NOT NULL DEFAULT TRUE, -- 是否显示
                                        last_show TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 最后一次显示时间
                                        is_leader BOOLEAN NOT NULL DEFAULT FALSE, -- 是否是群主，仅对群组有效
                                        is_self BOOLEAN NOT NULL DEFAULT FALSE, -- 是否是自己对自己的关系，仅对好友有效
                                        FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE ON UPDATE CASCADE,
                                        FOREIGN KEY (relation_id) REFERENCES relations(id) ON DELETE CASCADE ON UPDATE CASCADE
);

# -- 好友申请
CREATE TABLE IF NOT EXISTS applications (
#                                             id int primary key not null ,
                                            account1_id BIGINT NOT NULL, -- 申请者账号 id（外键）
                                            account2_id BIGINT NOT NULL, -- 被申请者账号 id（外键）
                                            apply_msg TEXT NOT NULL, -- 申请信息
                                            refuse_msg TEXT NOT NULL, -- 拒绝信息
                                            status ENUM('已申请', '已同意', '已拒绝') NOT NULL DEFAULT '已申请', -- 申请状态
                                            create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
                                            update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 更新时间
                                            PRIMARY KEY (account1_id, account2_id),
                                            FOREIGN KEY (account1_id) REFERENCES accounts(id) ON DELETE CASCADE ON UPDATE CASCADE,
                                            FOREIGN KEY (account2_id) REFERENCES accounts(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- 文件记录
CREATE TABLE IF NOT EXISTS files (
                                     id BIGINT AUTO_INCREMENT PRIMARY KEY, -- 文件 id
                                     file_name VARCHAR(255) NOT NULL, -- 文件名称
                                     file_type ENUM('img', 'file') NOT NULL, -- 文件类型
                                     file_size BIGINT NOT NULL, -- 文件大小 byte
                                     `key` VARCHAR(255) NOT NULL, -- 文件 key 用于删除文件
                                     url VARCHAR(255) NOT NULL, -- 文件 url
                                     relation_id BIGINT, -- 关系 id（外键）
                                     account_id BIGINT, -- 发送账号 id（外键）
                                     create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
                                     FOREIGN KEY (relation_id) REFERENCES relations(id) ON DELETE CASCADE ON UPDATE CASCADE,
                                     FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE ON UPDATE CASCADE
);
#
# -- 文件关系id索引
# CREATE INDEX file_relation_id ON files (relation_id);
#

CREATE TABLE IF NOT EXISTS messages (
                                        id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                        notify_type ENUM('system', 'common') NOT NULL,
                                        msg_type ENUM('text', 'file') NOT NULL,
                                        msg_content TEXT NOT NULL,
                                        msg_extend JSON,
                                        file_id BIGINT,
                                        account_id BIGINT,
                                        rly_msg_id BIGINT,
                                        relation_id BIGINT NOT NULL,
                                        create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        is_revoke BOOLEAN NOT NULL DEFAULT FALSE,
                                        is_top BOOLEAN NOT NULL DEFAULT FALSE,
                                        is_pin BOOLEAN NOT NULL DEFAULT FALSE,
                                        pin_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        read_ids JSON NOT NULL, -- 移除 DEFAULT '[]'
                                        msg_content_tsy TEXT,
                                        FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE ON UPDATE CASCADE,
                                        FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE SET NULL ON UPDATE CASCADE,
                                        FOREIGN KEY (rly_msg_id) REFERENCES messages(id) ON DELETE CASCADE ON UPDATE CASCADE,
                                        FOREIGN KEY (relation_id) REFERENCES relations(id) ON DELETE CASCADE ON UPDATE CASCADE
);

# -- 消息
# CREATE TABLE IF NOT EXISTS messages (
#                                         id BIGINT AUTO_INCREMENT PRIMARY KEY, -- 消息 id
#                                         notify_type ENUM('system', 'common') NOT NULL, -- 消息通知类型
#                                         msg_type ENUM('text', 'file') NOT NULL, -- 消息类型
#                                         msg_content TEXT NOT NULL, -- 消息内容
#                                         msg_extend JSON, -- 消息扩展信息
#                                         file_id BIGINT, -- 文件 id（外键）
#                                         account_id BIGINT, -- 发送账号 id（外键）
#                                         rly_msg_id BIGINT, -- 回复消息 id，没有则为 null（外键）
#                                         relation_id BIGINT NOT NULL, -- 关系 id（外键）
#                                         create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
#                                         is_revoke BOOLEAN NOT NULL DEFAULT FALSE, -- 是否撤回
#                                         is_top BOOLEAN NOT NULL DEFAULT FALSE, -- 是否置顶
#                                         is_pin BOOLEAN NOT NULL DEFAULT FALSE, -- 是否 pin
#                                         pin_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- pin时间
#                                         read_ids JSON NOT NULL DEFAULT '[]', -- 已读用户 id 集合（使用 JSON 数组替代 bigint[]）
#                                         msg_content_tsy TEXT, -- 消息分词
#                                         FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE ON UPDATE CASCADE,
#                                         FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE SET NULL ON UPDATE CASCADE,
#                                         FOREIGN KEY (rly_msg_id) REFERENCES messages(id) ON DELETE CASCADE ON UPDATE CASCADE,
#                                         FOREIGN KEY (relation_id) REFERENCES relations(id) ON DELETE CASCADE ON UPDATE CASCADE
# );
#
# -- 创建时间索引
# CREATE INDEX msg_create_at ON messages (create_at);
#
# -- 群通知
# CREATE TABLE IF NOT EXISTS group_notify (
#                                             id BIGINT AUTO_INCREMENT PRIMARY KEY, -- 群通知 id
#                                             relation_id BIGINT, -- 关系 id（外键）
#                                             msg_content TEXT NOT NULL, -- 消息内容
#                                             msg_expand JSON, -- 消息扩展信息
#                                             account_id BIGINT, -- 发送账号 id（外键）
#                                             create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
#                                             read_ids JSON NOT NULL DEFAULT '[]', -- 已读用户 id 集合
#                                             msg_content_tsv TEXT, -- 消息分词
#                                             FOREIGN KEY (relation_id) REFERENCES relations(id) ON DELETE CASCADE ON UPDATE CASCADE,
#                                             FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE ON UPDATE CASCADE
# );
#
# -- 分词索引
# CREATE INDEX group_notify_msg_content_tsv ON group_notify (msg_content_tsv);

# -- 创建更新时间戳的触发器
# ALTER TABLE applications
#     ADD COLUMN update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;
#
# -- 创建更新 applications 表的更新时间戳触发器
# CREATE TRIGGER update_application_timestamp
#     BEFORE UPDATE ON applications
#     FOR EACH ROW
# BEGIN
#     -- 如果需要手动更新时间戳
#     IF OLD.update_at <> NEW.update_at THEN
#         SET NEW.update_at = CURRENT_TIMESTAMP;
#     END IF;
# END;
# #
# -- 创建更新 settings 表 last_show 字段的触发器
# CREATE TRIGGER update_show_timestamp
#     BEFORE UPDATE ON settings
#     FOR EACH ROW
# BEGIN
#     IF NEW.is_show THEN
#         SET NEW.last_show = CURRENT_TIMESTAMP;
#     END IF;
# END;
# #
# -- 创建更新 settings 表 pin_time 字段的触发器
# CREATE TRIGGER update_pin_timestamp
#     BEFORE UPDATE ON settings
#     FOR EACH ROW
# BEGIN
#     IF NEW.is_pin THEN
#         SET NEW.pin_time = CURRENT_TIMESTAMP;
#     END IF;
# END;
#
# -- 创建更新 messages 表 pin_time 字段的触发器
# CREATE TRIGGER update_message_pin_timestamp
#     BEFORE UPDATE ON messages
#     FOR EACH ROW
# BEGIN
#     IF NEW.is_pin THEN
#         SET NEW.pin_time = CURRENT_TIMESTAMP;
#     END IF;
# END;
