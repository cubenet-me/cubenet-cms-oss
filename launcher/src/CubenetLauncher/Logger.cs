namespace CubenetLauncher;

public static class Logger
{
    private static readonly string _logPath = Path.Combine(
        Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData),
        Env.AppName,
        "launcher.log");

    static Logger()
    {
        var dir = Path.GetDirectoryName(LogPath);
        if (dir is not null)
            Directory.CreateDirectory(dir);
    }

    public static string LogPath => _logPath;

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
