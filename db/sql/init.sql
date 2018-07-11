USE challenge;
CREATE TABLE test
(
  col VARCHAR(10)
);

INSERT INTO test  (col) VALUES('ok');


CREATE TABLE IF NOT EXISTS user_
(
  user_id INT(10) NOT NULL AUTO_INCREMENT,
  username VARCHAR(45) not null DEFAULT "guest",
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  password VARCHAR(200) DEFAULT NULL,
  PRIMARY KEY (user_id),
  UNIQUE (username)
);



CREATE TABLE IF NOT EXISTS message
(
  message_id INT (10) NOT NULL AUTO_INCREMENT,
  message_uuid varchar(50),
  user_id int(10) DEFAULT 0,
  message_to_user_id int(10),
  message_type VARCHAR(30),
  message_content varchar(100),
  room_id int(10),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(message_id)
);

CREATE TABLE IF NOT EXISTS room
(
  room_id INT (10) NOT NULL AUTO_INCREMENT,
  room_uuid varchar(50),
  room_name varchar(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(room_id)
);

CREATE TABLE IF NOT EXISTS message_read
(
  read_id int(10) NOT NULL AUTO_INCREMENT,
  message_id INT (10) ,
  message_uuid VARCHAR (50) ,
  user_id INT (10),
  PRIMARY KEY(read_id)
);

CREATE TABLE IF NOT EXISTS room_user
(
  room_user_id INT (10) NOT NULL AUTO_INCREMENT,
  room_uuid varchar(50),
  user_id int(10),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(room_user_id)
);
