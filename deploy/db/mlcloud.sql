drop database if exists mlcloud;
create database mlcloud charset = utf8;

use mlcloud;

create table user (
 id int NOT NULL AUTO_INCREMENT,
 username varchar(32),
 password varchar(40) NOT NULL,
 salt varchar(40) DEFAULT NULL,
 sys_admin tinyint(1) DEFAULT 0 NOT NULL,
 creation_time datetime NOT NULL,
 update_time datetime NOT NULL,
 primary key (id),
 UNIQUE (username)
);

create table tensorflow_job (
  id int NOT NULL AUTO_INCREMENT,
  num_ps int DEFAULT 0 NOT NULL,
  num_workers int DEFAULT 1 NOT NULL,
  image varchar(256),
  log_dir varchar(256) NOT NULL,
  data_dir varchar(256),
  command varchar(256) NOT NULL,
  arguments varchar(256) NOT NULL,
  num_gpu int DEFAULT 0 NOT NULL,
  tensorboard BOOL DEFAULT FALSE NOT NULL,
  tensorboard_host VARCHAR(256) DEFAULT NULL,
  has_master BOOL DEFAULT FALSE NOT NULL,
  PRIMARY KEY (id)
);

create table mxnet_job (
  id int NOT NULL AUTO_INCREMENT,
  mode varchar(32) NOT NULL,
  num_ps int DEFAULT 0 NOT NULL,
  num_workers int DEFAULT 1 NOT NULL,
  image varchar(256),
  log_dir varchar(256) NOT NULL,
  data_dir varchar(256),
  command varchar(256) NOT NULL,
  arguments varchar(256) NOT NULL,
  num_gpu int DEFAULT 0 NOT NULL,
  PRIMARY KEY (id)
);

create table job (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(32),
  type varchar(32), # 'tensorflow', 'mxnet'
  user_id int,
  user_name varchar(32),
  creation_time datetime DEFAULT CURRENT_TIMESTAMP NOT NULL,
  update_time datetime DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES user (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
