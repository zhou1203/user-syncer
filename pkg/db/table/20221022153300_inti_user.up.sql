CREATE TABLE `user`
(
    USER_ID         varchar(60)  NOT NULL,
    LOGIN_NO        varchar(20)  NOT NULL DEFAULT '',
    USER_NAME       varchar(50)  NOT NULL,
    ORG_ID          varchar(100) NOT NULL,
    EMAIL           varchar(200) NOT NULL,
    MOBILE          varchar(50)  NOT NULL DEFAULT '',
    APP_ACCT_STATUS char(1)      NOT NULL DEFAULT '',
    PRIMARY KEY (USER_ID)
)

