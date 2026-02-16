-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（加密）',
    `avatar` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
    `bio` VARCHAR(500) DEFAULT NULL COMMENT '个人简介',
    `capacity` BIGINT NOT NULL DEFAULT 10737418240 COMMENT '总容量（默认10GB）',
    `used` BIGINT NOT NULL DEFAULT 0 COMMENT '已用空间',
    `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否激活',
    `last_login` TIMESTAMP NULL DEFAULT NULL COMMENT '最后登录时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_username` (`username`),
    UNIQUE INDEX `idx_email` (`email`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 创建文件表
-- 如果之前没有创建files表，现在添加
CREATE TABLE IF NOT EXISTS `files` (
    `id` VARCHAR(36) NOT NULL COMMENT '文件ID（UUID）',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '所属用户ID',
    `filename` VARCHAR(255) NOT NULL COMMENT '文件名',
    `path` VARCHAR(500) NOT NULL COMMENT '存储路径',
    `size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
    `hash` VARCHAR(64) DEFAULT NULL COMMENT '文件哈希（用于秒传）',
    `mime_type` VARCHAR(100) DEFAULT NULL COMMENT 'MIME类型',
    `extension` VARCHAR(20) DEFAULT NULL COMMENT '文件扩展名',
    `is_dir` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否为目录',
    `parent_id` VARCHAR(36) DEFAULT NULL COMMENT '父目录ID',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，2-冻结',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间（回收站）',
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_parent_id` (`parent_id`),
    INDEX `idx_hash` (`hash`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件表';

-- 创建分享表
CREATE TABLE IF NOT EXISTS `shares` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分享ID',
    `token` VARCHAR(64) NOT NULL COMMENT '分享令牌',
    `file_id` VARCHAR(36) NOT NULL COMMENT '文件ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '分享用户ID',
    `password` VARCHAR(255) DEFAULT NULL COMMENT '访问密码（加密）',
    `expire_time` TIMESTAMP NULL DEFAULT NULL COMMENT '过期时间',
    `max_downloads` INT DEFAULT 0 COMMENT '最大下载次数（0表示无限制）',
    `download_count` INT NOT NULL DEFAULT 0 COMMENT '已下载次数',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-有效，2-已取消',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_token` (`token`),
    INDEX `idx_file_id` (`file_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_expire_time` (`expire_time`),
    CONSTRAINT `fk_shares_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_shares_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分享表';

-- 创建收藏表
CREATE TABLE IF NOT EXISTS `favorites` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `file_id` VARCHAR(36) NOT NULL COMMENT '文件ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_user_file` (`user_id`, `file_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_file_id` (`file_id`),
    CONSTRAINT `fk_favorites_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_favorites_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';

-- 创建操作日志表
CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `user_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '用户ID',
    `action` VARCHAR(50) NOT NULL COMMENT '操作类型',
    `target_type` VARCHAR(20) NOT NULL COMMENT '目标类型',
    `target_id` VARCHAR(36) DEFAULT NULL COMMENT '目标ID',
    `details` JSON DEFAULT NULL COMMENT '操作详情',
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_action` (`action`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 创建分享表
CREATE TABLE IF NOT EXISTS `shares` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分享ID',
    `token` VARCHAR(64) NOT NULL COMMENT '分享令牌',
    `file_id` VARCHAR(36) NOT NULL COMMENT '文件ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '分享用户ID',
    `password` VARCHAR(255) DEFAULT NULL COMMENT '访问密码（加密）',
    `expire_time` DATETIME NULL DEFAULT NULL COMMENT '过期时间',
    `max_downloads` INT NOT NULL DEFAULT 0 COMMENT '最大下载次数（0表示无限制）',
    `download_count` INT NOT NULL DEFAULT 0 COMMENT '已下载次数',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-有效，2-已取消',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_token` (`token`),
    INDEX `idx_file_id` (`file_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_expire_time` (`expire_time`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分享表';

-- 创建收藏表
CREATE TABLE IF NOT EXISTS `favorites` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `file_id` VARCHAR(36) NOT NULL COMMENT '文件ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_user_file` (`user_id`, `file_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_file_id` (`file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';

-- 为回收站查询添加索引
CREATE INDEX IF NOT EXISTS idx_files_user_deleted ON files(user_id, deleted_at);

-- 为分享过期查询添加索引
CREATE INDEX IF NOT EXISTS idx_shares_expire ON shares(expire_time) WHERE status = 1;