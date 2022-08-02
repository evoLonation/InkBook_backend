create table users
(
    user_id  varchar(20) primary key not null,
    nickname varchar(20)             not null,
    realname varchar(20)             not null,
    password varchar(20)             not null,
    gender   char(5),
    intro    varchar(255),
    email    varchar(50) unique      not null
);

create table teams
(
    team_id    int primary key auto_increment not null,
    name       varchar(20)                    not null,
    intro      varchar(255),
    captain_id varchar(20)                    not null,
    foreign key (captain_id) references users (user_id) on delete cascade
);

create table team_member
(
    team_id   int         not null,
    member_id varchar(20) not null,
    identity  int,
    foreign key (member_id) references users (user_id) on delete cascade,
    foreign key (team_id) references teams (team_id) on delete cascade,
    primary key (team_id, member_id)
);

create table project
(
    proj_id     int primary key auto_increment not null,
    team_id     int                            not null,
    name        varchar(50)                    not null,
    create_time datetime                       not null default now(),
    is_delete   bool                           not null default false,
    delete_time datetime                       null,
    foreign key (team_id) references teams (team_id) on delete cascade
);

create table document
(
    doc_id      int primary key auto_increment not null,
    name        varchar(20)                    not null,
    proj_id     int                        not null,
    create_time datetime                       not null default now(),
    modifier_id varchar(20),
    modify_time datetime,
    is_editing  bool                           not null default false,
    is_delete   bool                           not null default false,
    deleter_id  varchar(20),
    delete_time datetime,
    foreign key (proj_id) references project (proj_id) on delete cascade,
    foreign key (deleter_id) references users (user_id) on delete set null,
    foreign key (modifier_id) references users (user_id) on delete set null
);

create table prototype
(
    proto_id    int primary key auto_increment not null,
    name        varchar(20)                    not null,
    proj_id     int                        not null,
    length      int                            not null,
    width       int                            not null,
    create_time datetime                       not null default now(),
    modifier_id varchar(20),
    modify_time datetime,
    is_editing  bool                           not null default false,
    is_delete   bool                           not null default false,
    deleter_id  varchar(20),
    delete_time datetime,
    foreign key (proj_id) references project (proj_id) on delete cascade,
    foreign key (deleter_id) references users (user_id) on delete set null,
    foreign key (modifier_id) references users (user_id) on delete set null
);