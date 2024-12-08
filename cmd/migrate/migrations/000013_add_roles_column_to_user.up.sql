ALTER TABLE users
   add column role_id INT references roles(id) NOT NULL default 1;
