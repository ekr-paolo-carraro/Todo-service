CREATE TABLE IF NOT EXISTS todos (
    id serial NOT NULL,
    userid VARCHAR(255),
    title VARCHAR(1024),
    done BOOLEAN,
    CONSTRAINT pkey_todos PRIMARY KEY (id)
)