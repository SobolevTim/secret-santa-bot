CREATE TABLE IF NOT EXISTS games (
    game_id SERIAL PRIMARY KEY,
    token VARCHAR(50) UNIQUE NOT NULL,              -- Уникальный токен игры
    organizer_id BIGINT NOT NULL,                   -- Telegram ID организатора
    name VARCHAR(100) NOT NULL,                     -- Название игры
    description TEXT,                               -- Описание игры
    is_drawn BOOLEAN DEFAULT FALSE,                 -- Статус жеребьевки (FALSE - еще не проведена)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Дата создания игры
);

CREATE TABLE IF NOT EXISTS participants (
    participant_id SERIAL PRIMARY KEY,
    game_id INT REFERENCES games(game_id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,                                    -- ID пользователя Telegram
    username VARCHAR(50),                                       -- Telegram username
    first_name VARCHAR(100),                                    -- Имя пользователя Telegram
    last_name VARCHAR(100),                                     -- Фамилия пользователя Telegram
    name VARCHAR(100),                                          -- Имя, введенное пользователем для игры
    gift_preferences TEXT,                                      -- Пожелания к подаркам
    assigned_to INT REFERENCES participants(participant_id) ON DELETE SET NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
