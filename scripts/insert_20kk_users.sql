-- Убедитесь, что расширение pgcrypto установлено
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Генерация случайных записей
DO $$
DECLARE
    i INTEGER;
    v_username VARCHAR(50);
    v_hash_password TEXT;
BEGIN
    FOR i IN 1..20000000 LOOP
        -- Генерация случайного имени пользователя
        v_username := 'user_' || i;

        -- Генерация случайного пароля (например, хэш MD5)
        v_hash_password := encode(digest(random()::text, 'md5'), 'hex');

        INSERT INTO users (username, hash_password)
        VALUES (v_username, v_hash_password);
        
        -- Вывод количества вставленных записей каждые 10000
        IF i % 10000 = 0 THEN
            RAISE NOTICE 'Inserted % records', i;
        END IF;
    END LOOP;
END $$;
