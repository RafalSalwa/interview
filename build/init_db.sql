use interview;

create table plan
(
    id          int auto_increment
        primary key,
    name        varchar(255)                         not null,
    description varchar(1000)                        not null,
    is_active   tinyint(1) default 0                 not null,
    created_at  datetime   default CURRENT_TIMESTAMP null,
    updated_at  datetime                             null,
    deleted_at  datetime                             null,
    constraint plan_name_uq
        unique (name)
)
    collate = utf8mb4_unicode_ci;

create index plan_created_at_index
    on plan (created_at);

create index u_plan_idx
    on plan (name);

create table user
(
    id                int auto_increment
        primary key,
    username          varchar(180)                         not null,
    password          varchar(255)                         not null,
    first_name        varchar(255)                         null,
    last_name         varchar(255)                         null,
    email             varchar(255)                         not null,
    phone_no          varchar(11)                          null,
    roles             json                                 null,
    verification_code VARCHAR(128)                         not null,
    is_verified       tinyint(1) default 0                 not null,
    is_active         tinyint(1) default 0                 not null,
    created_at        datetime   default CURRENT_TIMESTAMP not null,
    last_updated      datetime                             null,
    updated_at        datetime                             null,
    last_login        datetime                             null,
    deleted_at        datetime                             null
)
    collate = utf8mb4_unicode_ci;

create table payment
(
    id           int auto_increment
        primary key,
    user_id      int                                not null,
    plan_id      int                                not null,
    purchased_at datetime default CURRENT_TIMESTAMP null,
    started_at   datetime                           null,
    ends_at      datetime                           null,
    constraint subscriptions_user_id_uq
        unique (user_id),
    constraint payment_plan_id_fk
        foreign key (plan_id) references plan (id),
    constraint payment_user_id_fk
        foreign key (user_id) references user (id)
)
    collate = utf8mb4_unicode_ci;

create table subscription
(
    id           int auto_increment
        primary key,
    user_id      int                                not null,
    plan_id      int                                not null,
    purchased_at datetime default CURRENT_TIMESTAMP null,
    started_at   datetime                           null,
    ends_at      datetime                           null,
    constraint subscriptions_user_id_uq
        unique (user_id),
    constraint subscription_plan_id_fk
        foreign key (plan_id) references plan (id),
    constraint subscription_user_id_fk
        foreign key (user_id) references user (id)
)
    collate = utf8mb4_unicode_ci;

create index plan_created_at_index
    on subscription (purchased_at);

create index subscriptions_user_idx
    on subscription (user_id);

create table voucher
(
    id         int auto_increment
        primary key,
    code       varchar(25)                        not null,
    valid_for  varchar(12)                        null,
    created_at datetime default CURRENT_TIMESTAMP not null
);

create table orders
(
    id              int auto_increment
        primary key,
    user_id         int                                not null,
    subscription_id int                                not null,
    price           int                                not null,
    created_at      datetime default CURRENT_TIMESTAMP null,
    deleted_at      datetime                           null,
    voucher_id      int                                null,
    constraint subscriptions_user_id_uq
        unique (user_id),
    constraint orders_subscriptions_id_fk
        foreign key (subscription_id) references subscription (id),
    constraint orders_user_id_fk
        foreign key (user_id) references user (id),
    constraint orders_voucher_id_fk
        foreign key (voucher_id) references voucher (id)
)
    collate = utf8mb4_unicode_ci;



INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (1, 'user1', '$2a$13$5NmP4C6LG2wUzrkvZ/uVdue7QlZQNP/2FTFHo3/6QKmfWJD7YpAIa', 'username1', 'lastname1',
        'username1@domain1.com', '123456789', '{
        "roles": [
            "ROLE_ADMIN",
            "ROLE_USER"
        ]
    }', '624813', 1, 1, '2023-05-01 09:16:12', '2023-05-19 14:09:25', '2023-05-19 14:09:25', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (2, 'another_user', '$2a$13$5NmP4C6LG2wUzrkvZ/uVdue7QlZQNP/2FTFHo3/6QKmfWJD7YpAIa', 'another_user_firstname',
        'another_lastname', 'user@user.com', '987654321', '{
        "roles": "ROLE_USER"
    }', '176278', 1, 1, '2023-05-01 22:14:09', null, '2023-05-02 22:14:09', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (3, 'bad_user', '$2a$13$5NmP4C6LG2wUzrkvZ/uVdue7QlZQNP/2FTFHo3/6QKmfWJD7YpAIa', 'bad, bad,', 'very bad user',
        'trickortreat@strange.com', '123333777', '{
        "roles": "ROLE_USER"
    }', '347816', 1, 0, '2023-05-01 04:33:20', '2023-05-01 19:11:34', '2023-05-02 15:47:04', '2023-05-02 17:12:48');
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (5, 'vip', '$2a$13$5NmP4C6LG2wUzrkvZ/uVdue7QlZQNP/2FTFHo3/6QKmfWJD7YpAIa', 'very important', 'user',
        'memberplus@vip.com', '777333444', '{
        "roles": [
            "ROLE_VIP",
            "ROLE_USER"
        ]
    }', '642103', 1, 1, '2023-05-01 04:33:20', null, '2023-05-02 15:47:04', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (6, 'rafsal', '$2a$13$L4zZVIkX6o1KQLjCpwwxQ.8/cVm9ICx44AxkU.bHqJ11xFO4jvmT6', null, null, 'rafsal@interview.com',
        null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'OXSmjb', 0, 0, '2023-05-19 23:30:33', '2023-05-19 23:31:54', '2023-05-19 23:31:54', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (8, 'XhuytkQZ', '$2a$13$VGdoX2zgpMc0mATOxJzU2uY0GfTbCSj5TWXDWT2T/OTALGvhGsZhy', null, null,
        'XhuytkQZ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'ZBLgyz', 0, 0, '2023-05-19 23:37:02', '2023-05-19 23:37:03', '2023-05-19 23:37:03', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (9, 'xfKcHzAQ', '$2a$13$sSPqPJMMpi/4UdxpahYNOeZ4hXqZY6rmBis5AnLJtGM23tHHBc28W', null, null,
        'xfKcHzAQ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'tEwKvO', 0, 0, '2023-05-19 23:38:05', '2023-05-19 23:38:05', '2023-05-19 23:38:05', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (10, 'tTGKsjQN', '$2a$13$Z/XrRCc4ZI08WScWaojNW.e8luA/cmRwDY2pkpTGc22wMLSmS70Zm', null, null,
        'tTGKsjQN@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'sizXxy', 0, 0, '2023-05-19 23:38:12', '2023-05-19 23:38:13', '2023-05-19 23:38:13', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (11, 'YeehRLeX', '$2a$13$yo8ifJmbMacJ0XyPbHhOCuI8a7MewaW.0L9olyMQ2Y1/isdusb1pe', null, null,
        'YeehRLeX@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'fntbgo', 0, 0, '2023-05-19 23:39:25', '2023-05-19 23:39:25', '2023-05-19 23:39:25', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (12, 'oWInjirT', '$2a$13$MJwB2NYkwNJHND/CD2o5Vun/DaBtKAzSyWMJ7EF2EU1QlCdnoAOsC', null, null,
        'oWInjirT@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'holEiY', 0, 0, '2023-05-22 11:31:49', '2023-05-22 11:31:49', '2023-05-22 11:31:49', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (13, 'dcgpNuQZ', '$2a$13$51D3D1fa3la9l9f8YyBtjeHG7k0Lc9Hozv0xMfQSxavDvga/riWmG', null, null,
        'dcgpNuQZ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'eRaQlh', 0, 0, '2023-05-22 11:34:10', '2023-05-22 11:34:10', '2023-05-22 11:34:10', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (14, 'PytyKNpU', '$2a$13$MUDeyizw36l7MxdFpRP0yOsqMPNHcRZoWLxOTCRARYpxjD7XOr2TO', null, null,
        'PytyKNpU@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'kpEogF', 0, 0, '2023-05-22 11:38:51', '2023-05-22 11:38:52', '2023-05-22 11:38:52', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (15, 'YEJBdARM', '$2a$13$BF43Vnnd3pP8/ncDHIHufOPHfVgwjVQXu81nQkHMHacS5q31QYDwG', null, null,
        'YEJBdARM@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'sPBKVK', 0, 0, '2023-05-22 11:40:55', '2023-05-22 11:40:55', '2023-05-22 11:40:55', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (16, 'cUWejYWv', '$2a$13$uin75bnGkvF0YmM778N7yO6BNDCUjpA.mFpwBVvgEpnd9..v6iGgC', null, null,
        'cUWejYWv@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'EPlqZK', 0, 0, '2023-05-22 11:45:21', '2023-05-22 11:45:21', '2023-05-22 11:45:21', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (17, 'IWBQcesa', '$2a$13$NGvChpYxY/ppDbnhqOdaQ.xOYqoMEtArG8RRd5WPtz4tB6Gb9IA2m', null, null,
        'IWBQcesa@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'MHxxBT', 0, 0, '2023-05-22 11:45:57', '2023-05-22 11:45:58', '2023-05-22 11:45:58', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (18, 'OrLNaDcm', '$2a$13$Pkwqymj/./WbWCG46WwEYeDH8JFWUCDwsZq0JoYobwrhxo8ysO94e', null, null,
        'OrLNaDcm@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'lucTNW', 0, 0, '2023-05-22 11:46:47', '2023-05-22 11:46:47', '2023-05-22 11:46:47', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (19, 'TrNazrGn', '$2a$13$1t27jWwQnCdzlDpsb/940u0KaOF54tr9u/Ts7OcZ38SqGxBNfFWjG', null, null,
        'TrNazrGn@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'CrJlHV', 0, 0, '2023-05-22 11:48:51', '2023-05-22 11:48:51', '2023-05-22 11:48:51', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (20, 'eDsxWQWv', '$2a$13$Bk2CTP2OThAwzNzuH4ftMO4onJlqAXCyPi328obT/C5fzhXDnYGlW', null, null,
        'eDsxWQWv@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'hkCOHS', 0, 0, '2023-05-22 11:50:08', '2023-05-22 11:50:08', '2023-05-22 11:50:08', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (21, 'IAnjTNao', '$2a$13$W6kViSvl4g32C/s/EIYVEOmlNhF06XfqZwV73o.4W8fNH9SacYbC2', null, null,
        'IAnjTNao@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'mNMTiZ', 0, 0, '2023-05-22 11:58:15', '2023-05-22 11:58:15', '2023-05-22 11:58:15', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (22, 'WimEjAZL', '$2a$13$30yKE4vqTAimTNoZ7ZFQA.En7sfkuCh0hK1BxQ/SOeqNSH8OmL.yW', null, null,
        'WimEjAZL@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'svZHwo', 0, 0, '2023-05-22 12:02:20', '2023-05-22 12:02:20', '2023-05-22 12:02:20', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (23, 'DEPEzJxl', '$2a$13$/vCeaqorr18zJDQzCNFKYOahMYKn4g6wp8l9H2Vy2A1gjEezELx5O', null, null,
        'DEPEzJxl@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'OeHjnb', 0, 0, '2023-05-22 12:03:11', '2023-05-22 12:03:12', '2023-05-22 12:03:12', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (24, 'fQwNLgmW', '$2a$13$zkkhlHRPQVq7R3V2DuLfk.4EXZzOt.B0lHggRLthVzaSfsZtFOyOa', null, null,
        'fQwNLgmW@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'XpcaFx', 0, 0, '2023-05-22 12:03:49', '2023-05-22 12:03:49', '2023-05-22 12:03:49', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (25, 'rTBlHytA', '$2a$13$0HV5EIRcGSmvpH3QsMGQHuv0sVsejZihDjjgXawpBbrZThngX2Nla', null, null,
        'rTBlHytA@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'CYmKmU', 0, 0, '2023-05-22 12:04:16', '2023-05-22 12:04:17', '2023-05-22 12:04:17', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (26, 'fkPclayr', '$2a$13$U6p0qVNOH0LCQTKz92sNYOX4hrH1WtLcjuJTKn.mF6IiO.CMK/8YC', null, null,
        'fkPclayr@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'SHeusU', 0, 0, '2023-05-22 12:27:09', '2023-05-22 12:27:10', '2023-05-22 12:27:10', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (27, 'CUPRTqtR', '$2a$13$/dWp1cMTRmOhn7wPaHXRmuj3Z5fjqBGdfD8P6LdCmyDed1XwUS9.G', null, null,
        'CUPRTqtR@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'YBQfXQ', 0, 0, '2023-05-22 12:28:55', '2023-05-22 12:28:55', '2023-05-22 12:28:55', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (28, 'QBaViwvv', '$2a$13$k7sDNK6VLbYf7hTFvvmKIOUyf4ge2JzjguhtMHrXlovCH4s0Q0MCC', null, null,
        'QBaViwvv@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'OxXdwz', 0, 0, '2023-05-22 12:29:48', '2023-05-22 12:29:48', '2023-05-22 12:29:48', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (29, 'dgIIEPea', '$2a$13$wgHwRoISSf6Gnf4OTAeg7eo7KajnWCR7qt4hSPgLOmivGKXFplyrm', null, null,
        'dgIIEPea@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'iyAHgM', 0, 0, '2023-05-22 12:30:30', '2023-05-22 12:30:31', '2023-05-22 12:30:31', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (30, 'rqohkFsD', '$2a$13$.F5TOol316l6O/ACPPdeGuD/qB9j6gihs3xXOUd5.sKfQlPeudOQi', null, null,
        'rqohkFsD@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'gQvhOk', 0, 0, '2023-05-22 12:31:47', '2023-05-22 12:31:47', '2023-05-22 12:31:47', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (31, 'ujqzEqLe', '$2a$13$Gz43pXVg.HPxlyvR9Ju7I.S0wJIBxfaTijtS2LbIyyexWHjLEid7W', null, null,
        'ujqzEqLe@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'UJKiqg', 0, 0, '2023-05-22 12:33:18', '2023-05-22 12:33:18', '2023-05-22 12:33:18', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (32, 'NTyemSFU', '$2a$13$NlkqtBzqQ/hLMLfPP6dhPe37ja1qP9Fpgk810AudSc1Xrr8N/3iRm', null, null,
        'NTyemSFU@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'GEAXKO', 1, 0, '2023-05-22 12:34:05', '2023-05-22 12:34:05', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (33, 'GJgAqXbo', '$2a$13$BXIq71glwSEFKJBMLxuKJerXWBIsbGa.CFfGnt/.34RvurWlWDRFe', null, null,
        'GJgAqXbo@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'WSknwp', 1, 0, '2023-05-22 12:38:55', '2023-05-22 12:38:55', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (34, 'qBeZvsvN', '$2a$13$Ff9fRvoN7QB37J9MXHynRuNH8f8IPrfzSkZGBI1WqamsBd0h61xJa', null, null,
        'qBeZvsvN@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'SCxzko', 1, 0, '2023-05-22 12:40:15', '2023-05-22 12:40:15', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (35, 'esTJCGEu', '$2a$13$9Ll9r4ZgHmSsZ/EVITArzOqcE0hmr82BS3LQQ4.xf8MMY9RXxledm', null, null,
        'esTJCGEu@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'kuSici', 1, 0, '2023-05-22 12:41:14', '2023-05-22 12:41:14', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (36, 'wKUuVGhM', '$2a$13$N.Dl2OIDCGeSL94bPVdYfuqdbv.Zg1CbR8cHDokfX/i4wKEAn1oDO', null, null,
        'wKUuVGhM@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'sbqxhv', 1, 0, '2023-05-22 12:41:53', '2023-05-22 12:41:53', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (37, 'mOoRfpLg', '$2a$13$FvdalH9OumJb9cROe9xIFevCwty7WbFVep68Age3Fhg8xBV6HXwfu', null, null,
        'mOoRfpLg@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'rEJjhL', 1, 0, '2023-05-22 12:42:15', '2023-05-22 12:42:15', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (38, 'xepCwgbN', '$2a$13$yI1pvrOQzK6TbpktdhjTI.iPxc4TXbhnIasaneERlb5/fAdcyo1rG', null, null,
        'xepCwgbN@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'MuJmhb', 1, 0, '2023-05-22 12:47:22', '2023-05-22 12:47:22', '2023-05-22 12:47:22', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (39, 'kfoibJtl', '$2a$13$mvl/KERxriFVsp97OeAL8.Ra3fAjh6UVnrMDKUUUJ2mLcTC6YoY5e', null, null,
        'kfoibJtl@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'bienSP', 1, 1, '2023-05-22 12:48:32', '2023-05-22 12:48:33', '2023-05-22 12:48:33', null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (40, 'vdIZlrZb', '$2a$13$Xs9ldX.qxSu2B5m6UiDO/u7/uECOxZWfQJ/faV5lD49llWYSVQXsm', null, null,
        'vdIZlrZb@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'NyacRv', 0, 0, '2023-05-22 16:27:26', '2023-05-22 16:27:26', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (41, 'CBrpLIBy', '$2a$13$4aML.dO35OlV6pl84JepuOgnKzTmHpxn93HXFU3b/uSWQUKFQoqSu', null, null,
        'CBrpLIBy@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'cXHPed', 0, 0, '2023-05-22 16:27:26', '2023-05-22 16:27:26', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (42, 'LRXJeSMp', '$2a$13$MHB0Etx2kw8OUxXBUno8s.bp1cf/NOpi318K/rtrKC0loA6CI7Jvu', null, null,
        'LRXJeSMp@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'rnVWdw', 0, 0, '2023-05-22 16:27:26', '2023-05-22 16:27:26', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (43, 'rWTcANbD', '$2a$13$o0eVoJxgQtpHy5nBFLIIj.1KolCItZqg53y5RhrN/VZBygzlUWqUy', null, null,
        'rWTcANbD@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'TdQpGM', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (44, 'HmBCtOHg', '$2a$13$XxXys9jKh0oXeDjzd5843e.oiYzdQEp7GGuOUPrIa0JnAv.8oJj4m', null, null,
        'HmBCtOHg@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'NFmrlf', 0, 0, '2023-05-22 16:27:26', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (45, 'qTqRQllr', '$2a$13$aU/6R0g89N1zfg0.vPXFkOcm.wC1LOkUKyCSbvgpmAH57//uKPcnm', null, null,
        'qTqRQllr@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'SBBNtf', 0, 0, '2023-05-22 16:27:26', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (46, 'bwecyjBu', '$2a$13$FnsGGxwf2pv65FXaZPC4u.sLk0B5/rmu2OPP46ahrrC9S1EDY8F4u', null, null,
        'bwecyjBu@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'leaRaU', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (47, 'KHsBoeqr', '$2a$13$9UgbpjtfZFxgCmU01wEliOF7ru5nDOUSIyrgojeuYV7SDyLAAK/LS', null, null,
        'KHsBoeqr@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'FyEBUb', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (48, 'IqofCsMN', '$2a$13$i8bu4fJzcNW7UQIo8aiNI.B8wjjzeWuRnvd/Y6bM0uCFszKA.CwyS', null, null,
        'IqofCsMN@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'cEnwya', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (49, 'trBrNRKz', '$2a$13$OHprwxZevnzl.kXJhGWxVe.JiUU1wGNYZwsWvSUfZ39klj.MNZRM2', null, null,
        'trBrNRKz@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'uJCCCz', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (50, 'tcOORZsV', '$2a$13$/KCGkpQ3hGpIgjBwd/S0K.A9PwPpScsAbVaB8.Wi8ALNO3mSsrtbC', null, null,
        'tcOORZsV@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'dbZPwe', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (51, 'TYsVSmbj', '$2a$13$R/Yqd9TpfZScPrIllKQofOKWluvN9p96xoKLwgSWqulHJY0.FdEBm', null, null,
        'TYsVSmbj@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'tyjyBT', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (52, 'ySCjtJFC', '$2a$13$8OgEXcU5SZQtt8TERkaXDOVEGnhJ8hwEF/96VIURVACBWuyBS5TgC', null, null,
        'ySCjtJFC@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'jQXGsE', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (53, 'kITcaZuA', '$2a$13$imKAjaWdFw2aI2Y71xjOF.Hx5kCpsLYTo.vfNjAVVBnCVu.rXtGvO', null, null,
        'kITcaZuA@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'kecJWI', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (54, 'vVyGUidi', '$2a$13$/ga/mG1n43bC3sz5UVqykOq.0MXEzmhxt7T45jT5Elw1Dob6cleoi', null, null,
        'vVyGUidi@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'fYvuop', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (55, 'HXVsTOJG', '$2a$13$p0QoZ7msrIJGMAcjvoyMVeL7h2NRAQn9R4rp6zsFB0zwOLbv4ab42', null, null,
        'HXVsTOJG@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'ldpcmU', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (56, 'CRpYdUhb', '$2a$13$gR2YoXJxouI7mCaNfboiA.70Xs4Q0LHK8PX/mIqIkDOkMNywuZwri', null, null,
        'CRpYdUhb@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'rNpcfB', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (57, 'nPTGQJwI', '$2a$13$pwc9sDiG1fSbXNvnQzH7qeUuQllzjLrsNF02cgbZERV2uK0kGbCDW', null, null,
        'nPTGQJwI@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'fEfCiO', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (58, 'aMCeLsxf', '$2a$13$P4Mq0F5ltIh57O9a0osGDu/AEoXg.O23G.AVm/ZidajkJhS1C3HLe', null, null,
        'aMCeLsxf@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'MxJmIN', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (59, 'gwhPdHBQ', '$2a$13$TJeW2MmbkbSHmnd29NS1OOodC6xIii7lCB3fijwadSQrdBPFL3CvS', null, null,
        'gwhPdHBQ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'zNleeT', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (60, 'KKAqTheB', '$2a$13$4zbfeFN3YK9p6vmlXMlvVezRDtPL/NgGVkC4Pr4Unq4eYuuNPkgiW', null, null,
        'KKAqTheB@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'fRDUGb', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (61, 'vajrVjGm', '$2a$13$eqXXCd2l.eESXBwjS68ZZepp4ZPCBZ6HiARVszDVKVa9p3V26VInS', null, null,
        'vajrVjGm@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'qyKMPm', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (62, 'IuATaFtM', '$2a$13$WH1mwTVAHRsBheebF64CHev6pPBKOgzSdzUBJL2iB1bRIiduSgWcK', null, null,
        'IuATaFtM@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'pVeukN', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (63, 'YwSIwZmU', '$2a$13$tT2RodoMRYkpVsiyZg.15utQrQZU2s6Dg.HAwOhd14SVtxDlSsUbS', null, null,
        'YwSIwZmU@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'OBAumH', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (64, 'omeqkQEM', '$2a$13$zX/BGHvUYvciaA0VtxWoNOpcIs4vwVvzgk68Qa1asIuCSMDA1GZDS', null, null,
        'omeqkQEM@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'djJQqm', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (65, 'ZVcifskV', '$2a$13$nAhCyPDiA.cvFtIVu37hDOk.bHN2XMn5sCJq.4qOrW.a1oLnv1WlK', null, null,
        'ZVcifskV@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'daYmaV', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (66, 'gROSzVVT', '$2a$13$TU3REMSuwSQHEV6ph6r6vufDxUQ0X5.buh00tePYyDyaB.wB5FFem', null, null,
        'gROSzVVT@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'Oewxxi', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (67, 'NptUgmRw', '$2a$13$rKbRp5Crnk3sd0N/DYtmwOqDdo6nOG7P/ErDNApcOTT/u3xQ93IRa', null, null,
        'NptUgmRw@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'paOPEN', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (68, 'LSNoSpwr', '$2a$13$a1vUMDZGvagt10FTqf6aG.z8EnNqi7j75Kc8o5rs.d.pcv1H2NTay', null, null,
        'LSNoSpwr@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'KZcZec', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (69, 'jxLrlXNs', '$2a$13$JM5ckoijwflTSZaWh9D0Zuti2GX0j3dOOJ10oQAFYKrXT1sIpagSu', null, null,
        'jxLrlXNs@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'gZkrIS', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (70, 'UlryzEth', '$2a$13$6ZhdBV5rWjwkZbwxL1zieuKd16ukEY96WAR3ZW2VihXnztMtZL9lW', null, null,
        'UlryzEth@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'gPqbAu', 0, 0, '2023-05-22 16:27:27', '2023-05-22 16:27:27', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (71, 'TodPophd', '$2a$13$T.HWVNLh1KI8VEj.lKJUyuKwb4H/rblkD79MkNyhpBTiEgaGGGpRy', null, null,
        'TodPophd@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'VEFvew', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (72, 'sUJrGdbU', '$2a$13$gJBWEACcL4j2lySk/ZWleeSlBn5oybikN7ekAkeTUbkjIZKSFZQQG', null, null,
        'sUJrGdbU@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'EoxREc', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (73, 'kUuHcBPh', '$2a$13$UDp1hBQosJzPHodZGmuJVOuHkqo1snA05iZDt0.7b69qSIbKJlfua', null, null,
        'kUuHcBPh@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'TxPQmD', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (74, 'HDjBbsQN', '$2a$13$AZRr/dThsInjmRN/MorGLuSKNiWj4YLqgRwxRVnfwvwBTDQyelF1u', null, null,
        'HDjBbsQN@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'PDtleE', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (75, 'rqabRZox', '$2a$13$a7mdT4vahndvnTudqIrUo.FR0ZTeLqg74IGEtipKMGyM0pSo42WDm', null, null,
        'rqabRZox@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'FXhsEd', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (76, 'ycTNnrsk', '$2a$13$pRaEbSs3zsn9VOt8LBOnXujiG/TFPDvX25RssoM8JP214BPNng0WW', null, null,
        'ycTNnrsk@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'lTlaJV', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (77, 'YCYoFOHx', '$2a$13$U2e1Vkz7aOzk3d7qwuMhnOb0psW7mDnFYrDuqxZGlN/6vwIHNroqC', null, null,
        'YCYoFOHx@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'bKTzBD', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (78, 'UoCYCaZa', '$2a$13$BWZxEkWMkj8zHFfCXACcHeSW0YmJh0.oFJJY9xLs.jhnYTO.5koz6', null, null,
        'UoCYCaZa@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'wWwyzp', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (79, 'ejKhNcoR', '$2a$13$R5S6kMj7R4WiyWgMyBykh.vLRlT1oumjmHTBLJvOMynY0GWyncovm', null, null,
        'ejKhNcoR@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'vxJiGa', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (80, 'sPEFZtyt', '$2a$13$yRQ26ry1sr6CjaTIihqxue7Bmp9CPjcNgfmAoQApezVOKSAzGBKee', null, null,
        'sPEFZtyt@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'LWzgiV', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (81, 'mxSeFZSA', '$2a$13$pvng5arEaEKBAJSHapJwZONoiCtE7FPduTWt9R5AssB.ED3efPyeO', null, null,
        'mxSeFZSA@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'UapzGQ', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (82, 'EqnhDoaA', '$2a$13$9sM.sTbg6rcmtccVVxmJDeyviaxgycvRDrufoGlurU1mb6zTF70F2', null, null,
        'EqnhDoaA@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'mPndQr', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (83, 'zsQEyTtS', '$2a$13$hmuJrWI6T/JXrYNYWUDKpeEso2VtxMsYemPSKcDp2rMyILCmeuzsW', null, null,
        'zsQEyTtS@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'ejZTjJ', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (84, 'TWByPZbm', '$2a$13$iCB5BYL8d2FoJCWUoG2RhuW2farlFk.hjpFx/9vwk6O1e9zHxzbFS', null, null,
        'TWByPZbm@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'FRoiSs', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (85, 'sTFtrtyX', '$2a$13$G3WBl8EdpPD04fuc2U7MpeAbnNn/nGnr031TdzUqP6iq3gzkn79tW', null, null,
        'sTFtrtyX@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'jISniv', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (86, 'OeIjKtxL', '$2a$13$jDvkxkWqLRyyd9FEiDTGO.U4IVoXBg/tpNlcCvXKh6p3ipmbejV26', null, null,
        'OeIjKtxL@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'frPiuK', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (87, 'itbhQTKu', '$2a$13$UdVPNMV6VEPMj6.NR4gWSe0yk7kgcitA7jkhvx4Tb6s7cHW48OHFm', null, null,
        'itbhQTKu@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'aeYVqc', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (88, 'yWUOOTMI', '$2a$13$Us.xv3neOE50QZmBQyWkkuUtZGh4makNAuhKLVkwuO9hWR/4IJNvG', null, null,
        'yWUOOTMI@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'okIcGw', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (89, 'ldJQLozY', '$2a$13$LAUuXgYeCE3rF5tVAM6nK.WrvcoDZLpCBjaEp/svgTXqbkQbX9gQi', null, null,
        'ldJQLozY@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'dSbAUq', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:57', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (90, 'uXjVBqQu', '$2a$13$F7kr2cLCCbxzDsCpadVAVeBwhGwSajd5FvEx1/xXZI4bRgRu8Hgaq', null, null,
        'uXjVBqQu@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'KOBfhj', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (91, 'VTcfYXKd', '$2a$13$46INoJiq4XwwXk2mJfvwOuQIxszskZ3rTEnVZqiKPglDpy2JeJnYu', null, null,
        'VTcfYXKd@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'zRybkP', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (92, 'GdodmAHD', '$2a$13$kmqI.y/.qRQThkD1xqG56O.WMxGcKmTZ/I8/IS0NNiW2qG5dcs01G', null, null,
        'GdodmAHD@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'jRiwDK', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (93, 'ryzJVhMo', '$2a$13$u6d6TfbVgqD8434fvHWgv.daWiFEssiu65Lhi2RIjsXjHf2.wqnU6', null, null,
        'ryzJVhMo@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'obxRRI', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (94, 'sgMZqoMZ', '$2a$13$tKNVi.jhKgcnjD.SrDL6RuXtoencSjdNrXrV/H6qeBLFSqyFQ2fh6', null, null,
        'sgMZqoMZ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'JNpvRz', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (95, 'LxhFxpeI', '$2a$13$upeW3GDBFn/OsuGq2JHkO.1tvZSlk/Q8qlTuBprJZ.hX3OiZ4gCK.', null, null,
        'LxhFxpeI@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'hzaiHG', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (96, 'cxWxZQDC', '$2a$13$OBD9n1kBnqHmPSQEnPO2S.jVr79BQ/I8.VNHYPlfWg4tiN5B0.fgi', null, null,
        'cxWxZQDC@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'SVtjnf', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (97, 'EPfnuzCd', '$2a$13$VBKD6gTYgSguBu.J4rxpq.TRtMH91bBHV7RhS/bWP6xHCwfY7iF..', null, null,
        'EPfnuzCd@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'cgugxY', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (98, 'MouftQcQ', '$2a$13$ZbMO1z.RJGS4P5pHInVT5eP8ZuXn2/4vlX.XB3rHgOlc.axAL2Bfy', null, null,
        'MouftQcQ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'ckdrZn', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (99, 'xxfeauAc', '$2a$13$2TcwfHUMsEJUf3DepZyjquEwGOlZpRoq6I8YEBeyIpi/6DMgVXw8e', null, null,
        'xxfeauAc@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'wUQfjj', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (100, 'aOrFkhdD', '$2a$13$A6VaqH0dbemDzuFQ6g3PlOXGGvdzo/arRHuuWHv/I5F5c/SLqYjjC', null, null,
        'aOrFkhdD@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'gkXgnV', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (101, 'zfpDlHUM', '$2a$13$OZXujOP6ZpD3X94py7JQDOw0yRCBPkA9scWfImQwp5yf.mKdoQnGW', null, null,
        'zfpDlHUM@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'iIOtJa', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (102, 'bXCjNIdH', '$2a$13$J8seIKfjCmxVdTkGQd51U.52HwC6H/GCwjS4uZ576FZE4yTKe4EJu', null, null,
        'bXCjNIdH@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'vEQyHY', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (103, 'fxFwvQHc', '$2a$13$vg6t6j9gMVEnQ7dAWFLZ/OZl6cwpnyfcIBgLLLxLX/nd0RnRJstxe', null, null,
        'fxFwvQHc@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'EJyeJJ', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (104, 'tGuMsopL', '$2a$13$UIHxhNJGOmw5.d9Ozzhm3e6IGUSZcCdBUkyrTe1u0jnEpNLelU5fq', null, null,
        'tGuMsopL@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'IskjOX', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (105, 'DgrwoCDS', '$2a$13$P1/xlone3ZipKcYg6yDExOnk4ph6ck8UoCMWkCLoMlXzLZpkhGS/u', null, null,
        'DgrwoCDS@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'NmxCQz', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (106, 'PAFooGsz', '$2a$13$9Me9unrcuA.py4c/JY4uw.WqSP4bGKitiVKe.xZQlekNwvGiGjR2e', null, null,
        'PAFooGsz@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'xXbKzd', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (107, 'GsOdcIeO', '$2a$13$lnecVlbt0nkbQfBVh/95YO//L6zTX1QEDpdw2rdeUGFBv0hJ.0/x6', null, null,
        'GsOdcIeO@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'OUBcxl', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (108, 'YgVwOIwM', '$2a$13$M3Xy6uukGrZcibuGEYHkeOn8aRKQ9OeP4OuAj3iPbeZLvB358e3XW', null, null,
        'YgVwOIwM@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'qAmjhL', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (109, 'Odgydzwy', '$2a$13$LWWCA5UXQgb2EnNHNDeLR./ajQslUBRaD006UODlEFW0HXXK1wbcu', null, null,
        'Odgydzwy@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'eEqZbf', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (110, 'wtQvEknd', '$2a$13$npJ2SMLa1EN1uv/sc02gF.XzBNYVDBSz3htvLeVg3txfRIC9pCmle', null, null,
        'wtQvEknd@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'MnPWNI', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (111, 'eERTlmwP', '$2a$13$Q6ynoHtjPT689N1HdMOIhu6nEAyHJFz9xD6ABd1uxMSk4IPLJNxMG', null, null,
        'eERTlmwP@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'KEUsqF', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (112, 'NXeApYbx', '$2a$13$hkrpGX9Qg/M6/4j0G37Z6eScc8gHq7yPQncJR8LfVadXlKOx//bk6', null, null,
        'NXeApYbx@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'GAPVcQ', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (113, 'EXsMhQRE', '$2a$13$xkalVXyP2duL0VAidqFZhejMkjRuh3FyRd5bYSQhcqIrpwl40YfjS', null, null,
        'EXsMhQRE@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'yYmITO', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (114, 'KImEFGAG', '$2a$13$hbI/4jthkvL4bdLr.qCs2uL5gdxxEvyHj3SPSw3YX6h45ZHTKfZWe', null, null,
        'KImEFGAG@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'VHDpjv', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (115, 'HBfwfBWu', '$2a$13$.F6.EA2cbTZ7.rhrq0etoejJ4CTCQiKRjELpTEfnyvM0i7ijhp.3e', null, null,
        'HBfwfBWu@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'FOWEli', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (116, 'FxvlMOAh', '$2a$13$tHhC7Kaot5hcnb0MsDAyG.zRhXFCdeJMFrHMp3hkRZE92hQZ5ffFe', null, null,
        'FxvlMOAh@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'pmYLsf', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (117, 'DJSkEHUD', '$2a$13$qFq.VRGiCedGGVOlrAMaLOM.V3QLt7LebTr6kxaLd.Um0CHlHrZwa', null, null,
        'DJSkEHUD@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'wfcIYQ', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (118, 'lVhemADL', '$2a$13$Al.tT1T9zQ8Rh8bEy6LfzeibnC08UoJq3OrXrYMs6IH1JRKo9GKDu', null, null,
        'lVhemADL@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'CIxrFs', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (119, 'ptsTvrih', '$2a$13$nwEYEOUTvVZfuKD24hLPlOjNzC3j5C7TFqMlR4BCb6EPSX4CjoFqe', null, null,
        'ptsTvrih@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'ywjhOM', 0, 0, '2023-05-22 16:28:57', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (120, 'ILbCcppA', '$2a$13$i1YKtLRT81mMk1aCz4cezu./z94XesKb0A9AnYyK525K73aiSexfG', null, null,
        'ILbCcppA@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'kjSsLk', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (121, 'gYdBAfoe', '$2a$13$mrBUVnI/6.xoZ8EvxIe2iuc73Zu4bc/QCpFHAb2sVsKbN3p01m7DC', null, null,
        'gYdBAfoe@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'HcZDgw', 0, 0, '2023-05-22 16:28:58', '2023-05-22 16:28:58', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (122, 'qhGGHXbX', '$2a$13$Gs3kZ4ci8MJpuy2Hsrv2luePGfsM6/biV5MkmzY7LpzrsYPHZ243O', null, null,
        'qhGGHXbX@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'WuvlFx', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (123, 'RrNsuYdT', '$2a$13$P6ZA/0B6Cyvyf3vnz0qBt.9FD7LAc5E6SmS9XxCkg9F2TYgerPiIe', null, null,
        'RrNsuYdT@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'CpSctL', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (124, 'qutYyCjI', '$2a$13$DJUpD4bMNUFAoSZFnjxbReqNK8cNnP62IWhB29hH51wcXJYNgrJAu', null, null,
        'qutYyCjI@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'rhKbjM', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (125, 'AaEDcYIl', '$2a$13$fCIsIqrCAC32fpq5qAJoZO/kaahqZ16foyYZt/.OvsQ4qnhxE1wk2', null, null,
        'AaEDcYIl@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'PIQjFJ', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (126, 'SDSRHcfI', '$2a$13$hJ/VJL6UIuZxsgPyx7k0y.et6AFv5io9U.4qPEmsIzmEZGUsnliX.', null, null,
        'SDSRHcfI@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'iWpHjB', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (127, 'vpDwUVks', '$2a$13$uVezmPk/LGnozoHE3qcUtuBcZeY7.RgB327GMG5vGv.VzcHQY1OZS', null, null,
        'vpDwUVks@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'cUYSXs', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (128, 'quhMXTTI', '$2a$13$5yBvrr1DwnT5O56PYkhFAujeN5sgjM8G0mnKvmGEaJeFD.iqbTTJK', null, null,
        'quhMXTTI@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'oCRuYm', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (129, 'FsffnlXi', '$2a$13$ZR9kllLGoHFUXyetBSPCEucrDCzr2x0a5182PBhBS4SQqmLVrHi/G', null, null,
        'FsffnlXi@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'UzjTiQ', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (130, 'sPSzAoNA', '$2a$13$BYZ8cDU0AUe5upQ/3UgHx.n7FZqs61HMHN6/TSVzspsXTmva.Opxa', null, null,
        'sPSzAoNA@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'UHIdUL', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (131, 'DHJlyQKe', '$2a$13$h16KIul79531w9whRUAo2.Aa4Nd/ic09Hvhmjc9Wk3iscCVkf0pMC', null, null,
        'DHJlyQKe@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'LiZyLt', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (132, 'hmRlJvVM', '$2a$13$.NWSs0Xr6mqu15HRWwK1KeHHa0MaksxLdSkowySU1EmIhwmhCqUKy', null, null,
        'hmRlJvVM@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'vRhjOO', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (133, 'UDaxlZMT', '$2a$13$9rkC67xWtUkw4tLd2NSQ6O40ftOsGQLoUkjF59WVf1oBRXSx6U49.', null, null,
        'UDaxlZMT@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'TrEgbt', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (134, 'NaRodxeN', '$2a$13$E.lRNfoRllklnlq6JQpr.OlXNc9FI7SpIUaC39MTzJ0tfKYmwAOla', null, null,
        'NaRodxeN@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'BLZewU', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (135, 'dRADUflX', '$2a$13$RqySeUQnI46cNYvig74PgeXWElRr4Wqn6.NI104SVDIR9gyuqGoT.', null, null,
        'dRADUflX@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'DUnLoB', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (136, 'QuqswVAQ', '$2a$13$MQoYVnA6g39GRh/blc/xROFCAGyNPkjDlE44k3YWkedZjiMHxrPli', null, null,
        'QuqswVAQ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'WMDjOT', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (137, 'zQgSTwAv', '$2a$13$sfvkVUeEOkPbqYxQvz0FveWyJ46RGXTplL2md0/cGNPHtan8.iw3y', null, null,
        'zQgSTwAv@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'jMMXQi', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (138, 'GvAzJQHu', '$2a$13$Gym4kPjdOo7VkfJlolDXCOImyrYduIJpwa90XOt9Ktmi1OXNxwuUK', null, null,
        'GvAzJQHu@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'EwrriL', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (139, 'tWmrkNvE', '$2a$13$kwSZfpLrdAw7XK/iUHDxHOLaw5xFaXk3E1MM9DoVpjxhaVLghRFUu', null, null,
        'tWmrkNvE@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'vHNiGZ', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (140, 'QLCqMOmv', '$2a$13$lc2S5bRYXV3I5N//1xVxj.9Gu3wGgGgGZQjhuOarTMmKQkjauE8xK', null, null,
        'QLCqMOmv@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'jJLbue', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (141, 'XvHNNDVB', '$2a$13$dFa5iKkVyXxAWOMse/EjcuH5IcPQ88Q3mvgsEdDx7lbf6/9T/cbW6', null, null,
        'XvHNNDVB@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'MkHKQd', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (142, 'OQiBKnkR', '$2a$13$T8PvXYrV2OAkWPsfhwkVSeOGbr4.al/SO72wUomWi1qd1fPkuexOm', null, null,
        'OQiBKnkR@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'kTtfRx', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (143, 'DpcZqyFz', '$2a$13$TDcMinNp64m9V2q9wyIb0e97K1TjAGIqaCr161pnt3c/oHaogmdGS', null, null,
        'DpcZqyFz@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'VTQgHf', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (144, 'WUmHpzuP', '$2a$13$XM3HMF70dtexKl.potY2meCydWOsFrAySHYjlrq15gyOPVqJu9pZC', null, null,
        'WUmHpzuP@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'dQpynz', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (145, 'XvkqgPUV', '$2a$13$otSNwHndOSSyvhY4pS3L2eiAQxHUy7V9bDk.CGPA1d4p15Lmyyx6G', null, null,
        'XvkqgPUV@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'TZhXYW', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (146, 'SUPlLGfK', '$2a$13$28Fc6Ac.PR0MfgHJZUBFdehvLGbo30o3YdfNCIQFglXqqCM9iLoNS', null, null,
        'SUPlLGfK@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'vxFinn', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (147, 'xqAuAdcB', '$2a$13$uVqHH84vIPEzD9qDiVeBrO0ap58eobZ6kg2pIEcKzv0mJx9WuU2t.', null, null,
        'xqAuAdcB@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'akQoly', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (148, 'CxtztLIw', '$2a$13$kWHEamO92QQ2Z5J8g/HmMuz5pak/aS/CivegZ.Ymxf1u0nfZXam/W', null, null,
        'CxtztLIw@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'DYhACt', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (149, 'KcUdtfRC', '$2a$13$Kez.XDmVgzAi1iTkiTk23Ouit6WF7qR0wahsaH59WVZp8jMgJl776', null, null,
        'KcUdtfRC@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'PsloZj', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (150, 'sAcadrBX', '$2a$13$L7LT9LdbFyjiWj/dtJgEOesn7HTGCEdEd73cDNs6ZM5fd8d4qbulu', null, null,
        'sAcadrBX@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'KYALXF', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (151, 'VPypozwF', '$2a$13$9T8sj3JWOK8BvbdodfWsGeYYqN52Gy.ZvJapXgMkmhyaYeu4IkRa6', null, null,
        'VPypozwF@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'XotzMz', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (152, 'gPlnTDPG', '$2a$13$EK6tmicVRVJeSBxAHmHzsuFncWzt37S8ICUZFTKWRzBOI4bhymp7C', null, null,
        'gPlnTDPG@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'hEwyvX', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (153, 'KgfJeQSz', '$2a$13$f3fud1TovwQBD0KThxqw0.6yu8wxLuMYnoyG4U5f8oouKJCVJkGrW', null, null,
        'KgfJeQSz@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'cwSZzP', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (154, 'ljffIayR', '$2a$13$2sJNFB6FCMcOD/l8kF8z.eBt85Z3UqCjlQAjEzsar4Pt8xJAj.vVG', null, null,
        'ljffIayR@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'TuYqQM', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (155, 'IGgZDfcJ', '$2a$13$HGgVt9mGRFlMj6jm/Y5IUeTeFLfeUMny5wqdGh/l60AsLP9LuneK2', null, null,
        'IGgZDfcJ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'VVeMQO', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (156, 'ICDXVdAc', '$2a$13$DEN2EAm6bix1Qp.lKIqjTOTU5zVuZslc.vhvlf8M0PBIO1suRmdiK', null, null,
        'ICDXVdAc@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'ZiVXyk', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (157, 'GcTwwSmv', '$2a$13$yNvp4gEWU.AwBt4KMOFheu71pkLsnUuAp562C36MR9jTVWfR/WkW6', null, null,
        'GcTwwSmv@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'yBUwPL', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (158, 'VewuKkHs', '$2a$13$iwCPEjbqUQgyw6UuyIvkUOV.8lwiw74uohUMPk/RYjHnzRpoeiiJe', null, null,
        'VewuKkHs@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'NffJUp', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (159, 'juxsMwRy', '$2a$13$5lZ1E0sgzSoYRxgnb6Sfe.PTYpTK56wqklOgYHJljTfwQC2MhrCF.', null, null,
        'juxsMwRy@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'BBmsAN', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (160, 'kAZbknYL', '$2a$13$8U.whleCvjIu61xFTO.FYuplxW/8OL3IsN/vs4xv25m7OlAulVpWO', null, null,
        'kAZbknYL@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'iiBiiY', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (161, 'xpKbYzHj', '$2a$13$WSxk1xhmVXNP/ZAklowN1OtaI6XUa60SYqvjtvna.CGYdRoBkm2Y6', null, null,
        'xpKbYzHj@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'eslMih', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (162, 'cncPaayO', '$2a$13$C6ujuaFFv88hqQIzKIKbCenB45pyTEZ9tpBeFUuw27dVHh4Wu8fVu', null, null,
        'cncPaayO@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'YfAoqZ', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (163, 'BcgmQeqE', '$2a$13$dkkPnRH6c2kjEmNA0F7QkuDpEukTedX/oHgjPQHA.6RymZ06v7g2i', null, null,
        'BcgmQeqE@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'oeXfhC', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (164, 'JpvPHlOJ', '$2a$13$2oOXvzPexC3RbnbZFuiWW.F0Er68Oa71og56WvA7Rv2SENSUuEy0m', null, null,
        'JpvPHlOJ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'tcnlWQ', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (165, 'JEqgBsoS', '$2a$13$ffvF.ceUIyMbhNXIujtBFeZdk4HI4hzXBBd1VyGg7OjlMk4z2xIj6', null, null,
        'JEqgBsoS@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'DIaMKQ', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (166, 'ctSunqyZ', '$2a$13$JYLWNLURhSM3lXZFbAS0s.HwrVlKwVK7eH10AcXGsdcW95jH04mg.', null, null,
        'ctSunqyZ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'lCzMxe', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (167, 'ZHivVxNs', '$2a$13$O3YuPsfaU8Z9e.xVcBEefO1hFn/5zxLzT36kNANrwnOh/wiW8iwGu', null, null,
        'ZHivVxNs@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'DJMuIA', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (168, 'NoXSBpUC', '$2a$13$.fm3pzTbXtaK233VXN.Hs.q9ev7FsRMIAgOpT3aqmaLb0WgqtRLlS', null, null,
        'NoXSBpUC@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'sGaKEE', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (169, 'CkTEUMKF', '$2a$13$0wZ8lZjAII7JmZZ78DkpCO5RuDVzhbBMKmO7bv9z8urUC0T0DqOqW', null, null,
        'CkTEUMKF@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'awVejl', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (170, 'rmmvMgEv', '$2a$13$rofXfOoJyQx8lzG6MLo3o.BYnBUjwJnFRn3T373.BX7Xka0gKnmL.', null, null,
        'rmmvMgEv@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'SjYIAz', 0, 0, '2023-05-22 16:28:59', '2023-05-22 16:28:59', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (171, 'qzNuHnue', '$2a$13$psmnLEFjqxKr6vZq5uPfU.n9la0vz9GUCFjzXsdEP3ASOWFaDUrRq', null, null,
        'qzNuHnue@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'fBEFzm', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (172, 'prKoyfDX', '$2a$13$ewggORjXjX3f0yTBe1uIi..RbUWj3whAQ05cHJE.24LJzszRt6nGe', null, null,
        'prKoyfDX@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'ngKmgu', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (173, 'oSUVnbYl', '$2a$13$uv3SQ1HO74lcqtciPHmKmOhYA/vMMwJHDfNLeSjAcGgk0p/.x8j3m', null, null,
        'oSUVnbYl@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'iPOops', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (174, 'RuWLliYT', '$2a$13$/3Q5G7fuSlRo8/Bl5UYQcOHYSzmRI5g4CYiTQvnCyGNM4OuigV9be', null, null,
        'RuWLliYT@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'iTSNhO', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (175, 'JrxEWZvP', '$2a$13$Kf1BVqN3N.YcovzNuqGQu.s7df8oi.CXkYYjXht6LXOU1SxToGlP6', null, null,
        'JrxEWZvP@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'MkURba', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (176, 'HexYRrej', '$2a$13$/0UFhfrG8afVG4PlCb8Lue2/TJi0.jWF/7q/xbwLh7vEBUVVJnAYK', null, null,
        'HexYRrej@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'lGeQVC', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (177, 'KzJwOJhk', '$2a$13$M46saK4PPmi20KRKZoz78eFqxHIdnxO7xZmQyQ0lHd5z0trDveK8W', null, null,
        'KzJwOJhk@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'bJFEdK', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (178, 'NGQGIIzQ', '$2a$13$qcbwjPF3WRByTSaJoUAq5OHNR3UedsAlxWewujcstJTfsqBR0jn6e', null, null,
        'NGQGIIzQ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'yEjSyV', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (179, 'SLXwwAXm', '$2a$13$32ouI3cHi7a4qncqwchv5uWlO354B32so.KGlFsSgDxPLagquleoG', null, null,
        'SLXwwAXm@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'TiOytL', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (180, 'SyUcyxiV', '$2a$13$BVzOi8oG1kdyDZoNZUYR8edVk1lpQNGVubxAMOoQgw.6nQAj8dqnm', null, null,
        'SyUcyxiV@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'yxwAiD', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (181, 'BdgxTxmc', '$2a$13$UcKzm7NjyA5.qhgiBT6cfODttjaoXTJ8OhErDfnmHQWSYmGdJJ5OC', null, null,
        'BdgxTxmc@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'CcZgdZ', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (182, 'mngAhWwQ', '$2a$13$90ywTHPps5fkXCvvmsRuM.vfop0CAYkc6bSJymZ5aX4.PwEyEd6ea', null, null,
        'mngAhWwQ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'ghmfdI', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (183, 'DoNVxmqc', '$2a$13$xCmjVisz3PNrXAcdej0pNOJBHX.GfdcuQ0Tju9ZRBQITegitOyCTW', null, null,
        'DoNVxmqc@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'YCCXvu', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (184, 'IjCYdpot', '$2a$13$NjBbxqSmYsheM4aTJZEBVOsmc8OrUWGLdILF8BkRosy/2ufdUruky', null, null,
        'IjCYdpot@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'NALluE', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (185, 'jbaHXwWN', '$2a$13$fvyf3WYM6L1Yy.XdIK0qieSCRWUFoeijqW5LDMpfuudrYYE7ZAVDK', null, null,
        'jbaHXwWN@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'lxYZWt', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (186, 'UnIpyTON', '$2a$13$YXter5oZV.6hG5W0HmF8IO1XNzsX0zYobRTTh0d0h3wC/cW3Ixzk6', null, null,
        'UnIpyTON@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'RSityb', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (187, 'kZJTKDBF', '$2a$13$8gMgdutY9YpCfUBKVZ.4tOFgKhEklkSl1zxMQ5kAREqLWtUVKwJmO', null, null,
        'kZJTKDBF@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'UbbLII', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (188, 'FZxqqGHh', '$2a$13$KIyXP1y1v5JwpBfXYLpz1.PHavH4lXKizcDaeJn6ksQBnDCTFxZIy', null, null,
        'FZxqqGHh@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'GjpCXx', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (189, 'SjCvxJFA', '$2a$13$.S66qw16fqDG4VEpOgqEXOgtTom7FZ29oVMwinEulTlFJTXd/THCG', null, null,
        'SjCvxJFA@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'KjVsVU', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (190, 'nHRIwqBp', '$2a$13$MT92DTb2zH6ahlTmufIfP.ptPGTKw0Hs0VdLxoR/gMTIwjWwy5CT2', null, null,
        'nHRIwqBp@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'cQaekf', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (191, 'eXFlCnFY', '$2a$13$X.zolHqElJzkudtiBreGMuMiHPCz9FFwEIxP.Yh1ltJA0dcBP7IQq', null, null,
        'eXFlCnFY@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'rjygbn', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (192, 'zdsUGiBQ', '$2a$13$MhDjfgrvP/V2QHQHe7H9p.8.URxY1ZMXpDrnbBbEDmzJUFdYW0NOG', null, null,
        'zdsUGiBQ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'hdxiGK', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (193, 'vjHXfQRZ', '$2a$13$NMqNTW.fQVCt5TjypyGN3.kcU1Jv1xx40RB0lrrv7or.DRb2sqnC2', null, null,
        'vjHXfQRZ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'JLXFhQ', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (194, 'ZcIsKfoZ', '$2a$13$IKcg/Ko5kwIWAdKFq6NXmezB3bBfvehQ3A3UxveU52iODuSmBd/1e', null, null,
        'ZcIsKfoZ@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'nVXXWR', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (195, 'xMoIdRCu', '$2a$13$zDYvr22FVsAPf3/8vmixduax6z6Il75HA6BWeOVandAKhIdl011u6', null, null,
        'xMoIdRCu@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'VJRDFd', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (196, 'ISBnfJVn', '$2a$13$vqrk8YW2sCS.MIJcKo43bu8B.RBcggjcusU.i3cTfJD8f7P6QJj.2', null, null,
        'ISBnfJVn@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'OJuQjf', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (197, 'BLWzLqGf', '$2a$13$jieu2vNm1eBUJ1pA6/PmnOtNPnjiE.vB.NiR3CcJUM7vtI8sMuwrO', null, null,
        'BLWzLqGf@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'hRlpmD', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (198, 'dBkxEXtX', '$2a$13$jkV.H38tUrCSvJi9jLw87eJTQDLERgk30K1FaIJyU3s1m1x7vXzdu', null, null,
        'dBkxEXtX@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'gGottS', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (199, 'WopDhppL', '$2a$13$ts6PiLGWxDk4xOjbXlaFEuOxAt5xRPjCVpl9fVPZ2l9XFt0ZnlJY2', null, null,
        'WopDhppL@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'nmdNVe', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (200, 'obdGlype', '$2a$13$N/Xz85imdbMclBMt92mWq.S3eRx5oVuoR.EFH2sOFltvOfFmv0TPu', null, null,
        'obdGlype@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'trfdMo', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (201, 'aTwypKJc', '$2a$13$YP8jFhQLrAm8UCkPhfOEOuU.PStYOHfabOeg1wk9yCPf7ctH7Puoi', null, null,
        'aTwypKJc@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'bRBsFp', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (202, 'dSELeSjj', '$2a$13$/VDTOqm69yGDQULZbXECZuDQOGGzOVSAz4a9WXIqIxSfTwEk9SAsW', null, null,
        'dSELeSjj@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'NbNuzi', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (203, 'kXnDswaK', '$2a$13$ERlbhsFk14KCaTvnkaKq6OHzrY5OuXKCEysd2c9ER85RyfdLfuILq', null, null,
        'kXnDswaK@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'GBcqNY', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (204, 'uFtKwhCP', '$2a$13$Zv54mOoOE0ictPzu1F5TvObfQPV2hU5x4VCCcsePUbBBXdidd8a4y', null, null,
        'uFtKwhCP@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'afCXMQ', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);
INSERT INTO user (id, username, password, first_name, last_name, email, phone_no, roles, verification_code, is_verified,
                  is_active, created_at, updated_at, last_login, deleted_at)
VALUES (205, 'wstPkGMm', '$2a$13$d7tupW2Rd13McmRv8/zBge3z78G7wvl0nZpe8uHi9/xXEnO.y2tUK', null, null,
        'wstPkGMm@interview.com', null, '{
        "Roles": [
            "ROLE_USER"
        ]
    }', 'QNPEwr', 0, 0, '2023-05-22 16:31:38', '2023-05-22 16:31:38', null, null);

insert into plan (name, description, is_active, created_at, updated_at, deleted_at)
values ('freemium', 'Basic plan for all users, that allows platform usage', 1, '2023-05-22 19:49:34', null, null),
       ('plus', 'Additional plan with more offers for user', 1, '2023-05-22 19:49:34', null, null),
       ('gold', 'Basic plan for all users, that allows platform usage', 1, '2023-05-22 19:49:34', '2023-05-22 19:49:34',
        null),
       ('vip', 'exclusive offers for VIP members', 0, '2023-05-22 19:49:34', null, '2023-05-22 19:49:34'),
       ('500+', 'Many many discounts, everywhere', 1, '2023-05-22 19:49:34', null, null),
       ('800+', 'They pay You for visiting them', 1, '2023-05-22 19:49:34', null, null);
