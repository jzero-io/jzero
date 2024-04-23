CREATE TABLE credential (
                      id bigint AUTO_INCREMENT,
                      name varchar(255) NOT NULL DEFAULT '' COMMENT 'The credential name',
                      type tinyint(1) NULL DEFAULT 0 COMMENT 'The credential type, 0:normal,1:vip, for test golang keyword',
                      create_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                      update_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      PRIMARY KEY (id)
) ENGINE = InnoDB COLLATE utf8mb4_general_ci COMMENT 'credential table';