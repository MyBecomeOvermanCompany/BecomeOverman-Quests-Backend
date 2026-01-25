-- Обновить существующих пользователей: дать им 100 монет если у них 0
UPDATE users 
SET coin_balance = 100 
WHERE coin_balance = 0 OR coin_balance IS NULL;
