DROP DATABASE IF EXISTS DBaaS;
CREATE DATABASE DBaaS; 
USE DBaaS;

SET NAMES utf8mb4 ;
SET character_set_client = utf8mb4 ;


CREATE TABLE uploads (
  id 				int NOT NULL AUTO_INCREMENT,
  email 			varchar(320) NOT NULL,
  inputs 			varchar(512),
  program_language  ENUM('py', 'cpp', 'c', 'java') NOT NULL,
  is_enable 	 	bool NOT NULL DEFAULT FALSE, 
  PRIMARY KEY (id)
) AUTO_INCREMENT=1;

CREATE TABLE jobs (
  id 				int NOT NULL AUTO_INCREMENT,
  upload 			int,
  job_query 		varchar(512),
  job_status  		ENUM('executed', 'none') NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (upload) REFERENCES uploads(id)
) AUTO_INCREMENT=1;

CREATE TABLE results (
  id 				int NOT NULL AUTO_INCREMENT,
  job 				int,
  output 			varchar(1024),
  execute_status  	ENUM('in progress', 'done') NOT NULL,
  execute_date 	 	datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (job) REFERENCES jobs(id)
) AUTO_INCREMENT=1;
