-- 优化文件表索引
ALTER TABLE files ADD INDEX idx_user_status (user_id, status);
ALTER TABLE files ADD INDEX idx_parent_user (parent_id, user_id);
ALTER TABLE files ADD INDEX idx_created_at (created_at) WHERE deleted_at IS NULL;

-- 优化分享表索引
ALTER TABLE shares ADD INDEX idx_user_status (user_id, status);
ALTER TABLE shares ADD INDEX idx_token_status (token, status);
ALTER TABLE shares ADD INDEX idx_expire_status (expire_time, status) WHERE status = 1;

-- 优化收藏表索引
ALTER TABLE favorites ADD INDEX idx_user_created (user_id, created_at);