USE challenge;

CREATE TABLE test(col VARCHAR(10));

INSERT INTO test(col) VALUES('ok');


CREATE TABLE IF NOT EXISTS user_ (
  user_id INT(10) NOT NULL AUTO_INCREMENT,
  username VARCHAR(45) not null DEFAULT "guest",
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  password VARCHAR(200) DEFAULT NULL,
  PRIMARY KEY (user_id),
  UNIQUE(username)
);

CREATE TABLE IF NOT EXISTS user_activity_history(
  activity_id INT(10) NOT NULL AUTO_INCREMENT,
  user_id int(10) DEFAULT 0,
  active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (activity_id)
);

CREATE TABLE IF NOT EXISTS message (
  message_id INT(10) NOT NULL AUTO_INCREMENT,
  user_id int(10) DEFAULT 0,
  message_to_user_id int(10),
  message_content varchar(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (message_id)
);
