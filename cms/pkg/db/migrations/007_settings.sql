CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL DEFAULT ''
);

INSERT INTO settings (key, value) VALUES
    ('site_name', 'CubeNet CMS'),
    ('site_description', 'Мощная открытая CMS для управления Minecraft модовыми серверами')
ON CONFLICT (key) DO NOTHING;
