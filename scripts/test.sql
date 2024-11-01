INSERT INTO users (
    username,
    hash_password
) VALUES (
    'vladimir', '$2a$10$rsEwRf1jqI/RGXfKPaC.QuPBA/ayaMhESKxM2HAj9USxfuzJzAwbm'
)
ON CONFLICT (id) DO NOTHING
RETURNING 
    id,
    username,
    hash_password



UPDATE users 
SET 
    username = 'vladimir',
    hash_password = '$2a$10$rsEwRf1jqI/RGXfKPaC.QuPBA/ayaMhESKxM2HAj9USxfuzJzAwbm'
WHERE 
    id = '60d1c75c-8786-480c-9adf-58f7b6da0b93'
ON CONFLICT (id) DO NOTHING

INSERT INTO users (
    username,
    hash_password
) VALUES (
    'vladimir', 
    '$2a$10$rsEwRf1jqI/RGXfKPaC.QuPBA/ayaMhESKxM2HAj9USxfuzJzAwbm'
)
RETURNING 
    id,
    username,
    hash_password;




INSERT INTO test_users
(
    username,
    hash_password
)
SELECT 
    'user1'::VARCHAR(50), 
    'pass1'::TEXT
WHERE
    NOT EXISTS (
        SELECT id ,username,hash_password
        FROM test_users 
        WHERE username = 'user1'::VARCHAR(50) AND hash_password = 'pass1'
    )
RETURNING 
    id,
    username,
    hash_password;