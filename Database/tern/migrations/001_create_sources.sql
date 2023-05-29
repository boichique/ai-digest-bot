CREATE TABLE 
sources (
  userid INT,
  source VARCHAR(255) NOT NULL,
  PRIMARY KEY (userid, source)
);

---- create above / drop below ----

DROP TABLE sources;
