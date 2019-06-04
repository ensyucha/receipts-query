CREATE TABLE IF NOT EXISTS systems (
    password TEXT COMMENT '管理员密码',
    unusedusage INTEGER COMMENT '未分配余额',
    apicode TEXT COMMENT 'API Code'
);

CREATE TABLE IF NOT EXISTS users (
    userid INTEGER PRIMARY KEY AUTO_INCREMENT COMMENT '用户id',
    username VARCHAR(100) unique COMMENT '用户名',
    nickname VARCHAR(100) unique COMMENT '用户昵称',
    password TEXT COMMENT '用户密码',
    usages INTEGER COMMENT '用户余额',
    total INTEGER DEFAULT 0 COMMENT '用户历史使用额度'
);

CREATE TABLE IF NOT EXISTS logs (
    logid INTEGER PRIMARY KEY AUTO_INCREMENT COMMENT '日志id',
    logtime TEXT COMMENT '日志记录时间',
    logtype TEXT COMMENT '日志类型',
    loginfo TEXT COMMENT '日志信息',
    username VARCHAR(100) COMMENT '记录日志者'
);


DELETE FROM systems WHERE apicode = 'empty';

INSERT INTO systems(password, unusedusage, apicode) VALUES ('123456', 0, 'empty');