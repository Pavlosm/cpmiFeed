CREATE KEYSPACE cpmiEvents
WITH REPLICATION = {
  'class': 'SimpleStrategy',
  'replication_factor': 1
};

CREATE TABLE cpmiEvents.event (id int, data int, url text, description text, tags list<text>, timestamp timestamp,  PRIMARY KEY (id));