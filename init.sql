CREATE TABLE `deployment` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `project_id` INTEGER NOT NULL DEFAULT '0',
  `tag` VARCHAR(32) NOT NULL DEFAULT '' ,
  `description` VARCHAR(255) NOT NULL DEFAULT '',
  `status` VARCHAR(4) NOT NULL DEFAULT '0',
  `release_path` text NOT NULL DEFAULT '',
  `deployment_path` text NOT NULL DEFAULT '',
  `output` text NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE project (id INTEGER PRIMARY KEY AUTOINCREMENT, owner_id VARCHAR (32) NOT NULL DEFAULT '0', name VARCHAR (32) NOT NULL DEFAULT '', code VARCHAR (32) NOT NULL DEFAULT '', description VARCHAR (64) NOT NULL DEFAULT '', lang VARCHAR (16) NOT NULL DEFAULT '', run_type VARCHAR (4) NOT NULL DEFAULT '0', source_path text NOT NULL DEFAULT (''), release_path text NOT NULL DEFAULT (''), deployment_status VARCHAR (4) NOT NULL DEFAULT '0', deployment_path text NOT NULL DEFAULT (''), tag VARCHAR (32) NOT NULL DEFAULT '', server_list text NOT NULL DEFAULT (''), start_script text NOT NULL DEFAULT (''), stop_script text NOT NULL DEFAULT (''), restart_script text NOT NULL DEFAULT (''), before_deployment_script text NOT NULL DEFAULT (''), deployment_script text NOT NULL DEFAULT (''), after_deployment_script text NOT NULL DEFAULT (''), package_script text NOT NULL DEFAULT (''), status VARCHAR (4) NOT NULL DEFAULT '0', deployed_at DATETIME DEFAULT (datetime('now')), created_at DATETIME NOT NULL DEFAULT (datetime('now')), updated_at DATETIME NOT NULL DEFAULT (datetime('now')));

CREATE TABLE server (id INTEGER PRIMARY KEY AUTOINCREMENT, ip VARCHAR (64) NOT NULL DEFAULT '', core_num INTEGER NOT NULL DEFAULT '0', memory_size INTEGER NOT NULL DEFAULT '0', username VARCHAR (64) NOT NULL DEFAULT (''), password VARCHAR (64) NOT NULL DEFAULT (''), private_key text NOT NULL DEFAULT (''), port INTEGER NOT NULL DEFAULT (22), created_at datetime NOT NULL DEFAULT (datetime('now')), updated_at datetime NOT NULL DEFAULT (datetime('now')));

CREATE TABLE `userinfo` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `username` VARCHAR(16) NOT NULL DEFAULT '' ,
  `email` VARCHAR(64) NOT NULL DEFAULT '' ,
  `ticket` VARCHAR(32) NOT NULL DEFAULT '' ,
  `role` VARCHAR(255) NOT NULL DEFAULT '' ,
  `enable` VARCHAR(3) NOT NULL DEFAULT '1' ,
  `created_at`  DATETIME NOT NULL DEFAULT (datetime('now')),
  `updated_at`  DATETIME NOT NULL DEFAULT (datetime('now'))
);