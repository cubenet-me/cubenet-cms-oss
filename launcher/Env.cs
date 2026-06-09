namespace CubenetLauncher;

public static class Env
{
    /// <summary>Название проекта (используется в путях).</summary>
    public const string AppName = "Cubenet";

    /// <summary>URL JSON-эндпоинта с информацией об обновлении.</summary>
    /// <remarks>Ответ: { "url": "...", "hash": "sha256hex" }</remarks>
    public const string UpdateCheckUrl = "https://example.com/launcher/update.json";
}
