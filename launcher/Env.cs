// ────────────────────────────────────────────
// Cubenet Launcher — конфигурация
// Меняй под свой проект перед билдом
// ────────────────────────────────────────────

namespace CubenetLauncher;

public static class Env
{
    /// <summary>Название проекта (используется в путях).</summary>
    public const string AppName = "Cubenet";

    /// <summary>URL для проверки обновлений.</summary>
    public const string UpdateCheckUrl = "https://example.com/launcher/version";

    /// <summary>URL для скачивания новой версии.</summary>
    public const string UpdateDownloadUrl = "https://example.com/launcher/latest";

    /// <summary>Версия лаунчера.</summary>
    public const string AppVersion = "1.0.0";
}
