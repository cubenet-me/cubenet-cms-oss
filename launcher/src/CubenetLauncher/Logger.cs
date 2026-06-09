namespace CubenetLauncher;

public static class Logger
{
    private static readonly string LogPath = Path.Combine(
        Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData),
        "Cubenet",
        "launcher.log");

    static Logger()
    {
        var dir = Path.GetDirectoryName(LogPath);
        if (dir is not null)
            Directory.CreateDirectory(dir);
    }

    public static void Info(string message) => Write("INFO", message);
    public static void Warn(string message) => Write("WARN", message);
    public static void Error(string message) => Write("ERROR", message);

    private static void Write(string level, string message)
    {
        var line = $"[{DateTime.UtcNow:yyyy-MM-dd HH:mm:ss}] [{level}] {message}";
        try
        {
            File.AppendAllText(LogPath, line + Environment.NewLine);
        }
        catch
        {
            // ignore write failures
        }
    }
}
