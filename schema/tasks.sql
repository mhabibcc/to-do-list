CREATE DATABASE "to-do-list";
\c "to-do-list"
CREATE TABLE IF NOT EXISTS tasks(
	id serial,
	task_name varchar NOT NULL,
	is_done bool NOT NULL,
	CONSTRAINT tasks_pk PRIMARY KEY (id)
);