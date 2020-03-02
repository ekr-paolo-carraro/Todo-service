UPDATE public.todos 
    SET userid = $2, title = $3, done = $4
    WHERE id = $1