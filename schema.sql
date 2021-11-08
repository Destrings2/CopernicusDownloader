create table file
(
    uuid                char(36)      not null,
    title               varchar(255)  not null,
    beginPosition       datetime      not null,
    endPosition         datetime      not null,
    ingestionData       datetime      not null,
    orbitNumber         int           not null,
    relativeOrbitNumber int           not null,
    productLevel        varchar(255)  not null,
    footprint           varchar(2047) not null,
    platformName        varchar(255)  not null,
    productType         varchar(255)  not null,
    fileFormat          varchar(255)  not null,
    filename            varchar(255)  not null,
    size                bigint        not null,
    online              tinyint(1)    not null,
    requestedDate       datetime      null,
    downloaded          tinyint(1)    not null,
    lockedBy            varchar(255)  null,
    localPath           varchar(255)  null,
    checkmd5            char(32)      null,
    priority            int default 1 not null,
    constraint file_uuid_uindex
        unique (uuid)
);

create index file_user_id_fk
    on file (lockedBy);

alter table file
    add primary key (uuid);