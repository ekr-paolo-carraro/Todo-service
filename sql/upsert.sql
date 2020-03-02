INSERT INTO todos (id, userid, title, done) 
VALUES ($1,$2,$3,$4) 
ON CONFLICT ON CONSTRAINT pkey_todos
DO UPDATE SET title = $3, done = $4
WHERE id = $1