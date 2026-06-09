CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier TEXT NOT NULL UNIQUE,
    name JSONB NOT NULL DEFAULT '{}'::jsonb,
    color TEXT NOT NULL DEFAULT '#94a3b8',
    permissions JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO roles (identifier, name, color, permissions) VALUES
    ('admin', '{"ru":"Администратор","en":"Administrator"}'::jsonb, '#a78bfa', '["admin.*", "wallet.*"]'::jsonb),
    ('moderator', '{"ru":"Модератор","en":"Moderator"}'::jsonb, '#6ee7b7', '["admin.access", "admin.manage"]'::jsonb),
    ('staff', '{"ru":"Персонал","en":"Staff"}'::jsonb, '#6ee7b7', '["admin.access", "admin.manage"]'::jsonb),
    ('user', '{"ru":"Пользователь","en":"User"}'::jsonb, '#94a3b8', '[]'::jsonb)
ON CONFLICT (identifier) DO NOTHING;

ALTER TABLE users ADD COLUMN IF NOT EXISTS role_id UUID REFERENCES roles(id);

UPDATE users SET role_id = (SELECT id FROM roles WHERE identifier = users.role) WHERE role_id IS NULL;

ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;

ALTER TABLE users ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id);
