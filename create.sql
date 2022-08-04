create table users
(
    user_id  varchar(20) primary key not null,
    nickname varchar(20)             not null,
    realname varchar(20),
    password varchar(20)             not null,
    gender   char(5),
    intro    varchar(255),
    email    varchar(50) unique      not null,
    url      varchar(255) default '_defaultavatar.webp'
);

create table teams
(
    team_id    int primary key auto_increment not null,
    name       varchar(20) not null,
    intro      varchar(255),
    captain_id varchar(20) not null,
    url        varchar(255) default '_defaultavatar.webp',
    foreign key (captain_id) references users (user_id) on delete cascade
);

create table team_members
(
    team_id   int         not null,
    member_id varchar(20) not null,
    identity  int,
    foreign key (member_id) references users (user_id) on delete cascade,
    foreign key (team_id) references teams (team_id) on delete cascade,
    primary key (team_id, member_id)
);

create table projects
(
    project_id  int primary key auto_increment not null,
    team_id     int         not null,
    name        varchar(50) not null,
    creator_id  varchar(20) not null,
    create_time datetime    not null default now(),
    is_deleted  bool        not null default false,
    delete_time datetime             default null,
    intro       varchar(255),
    img_url     varchar(255),
    foreign key (team_id) references teams (team_id) on delete cascade,
    foreign key (creator_id) references users (user_id) on delete cascade
);

create table documents
(
    doc_id      int primary key auto_increment not null,
    name        varchar(20) not null,
    project_id  int         not null,
    creator_id  varchar(20) not null,
    create_time datetime    not null default now(),
    modifier_id varchar(20),
    modify_time datetime             default null,
    is_editing  bool        not null default false,
    is_deleted  bool        not null default false,
    deleter_id  varchar(20),
    delete_time datetime             default null,
    content     json,
    editing_cnt int         not null default 0,
    foreign key (project_id) references projects (project_id) on delete cascade,
    foreign key (creator_id) references users (user_id) on delete cascade,
    foreign key (deleter_id) references users (user_id) on delete set null,
    foreign key (modifier_id) references users (user_id) on delete set null
);

create table prototypes
(
    proto_id    int primary key auto_increment not null,
    name        varchar(20) not null,
    project_id  int         not null,
    length      int         not null,
    width       int         not null,
    creator_id  varchar(20) not null,
    create_time datetime    not null default now(),
    modifier_id varchar(20),
    modify_time datetime             default null,
    is_editing  bool        not null default false,
    is_deleted  bool        not null default false,
    deleter_id  varchar(20),
    delete_time datetime,
    foreign key (project_id) references projects (project_id) on delete cascade,
    foreign key (creator_id) references users (user_id) on delete cascade,
    foreign key (deleter_id) references users (user_id) on delete set null,
    foreign key (modifier_id) references users (user_id) on delete set null
);

create table graphs
(
    graph_id    int primary key auto_increment not null,
    name        varchar(20) not null,
    project_id  int         not null,
    creator_id  varchar(20) not null,
    create_time datetime    not null default now(),
    modifier_id varchar(20),
    modify_time datetime             default null,
    is_editing  bool        not null default false,
    is_deleted  bool        not null default false,
    deleter_id  varchar(20),
    delete_time datetime             default null,
    content     json,
    editing_cnt int         not null default 0,
    foreign key (project_id) references projects (project_id) on delete cascade,
    foreign key (creator_id) references users (user_id) on delete cascade,
    foreign key (deleter_id) references users (user_id) on delete set null,
    foreign key (modifier_id) references users (user_id) on delete set null
);
