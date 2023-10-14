CREATE TABLE IF NOT EXISTS tasks AUTHORIZATION pg_database_owner(
	id bigint NOT NULL,
	task_name varchar NOT NULL,
	is_done bool NOT NULL,
	CONSTRAINT tasks_pk PRIMARY KEY (id)
);
