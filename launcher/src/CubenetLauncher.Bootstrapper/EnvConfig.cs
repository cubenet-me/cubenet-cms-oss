namespace CubenetLauncher.Bootstrapper;

public class EnvConfig
{
    public string AppName { get; private set; } = "Cubenet";
    public string LauncherDownloadUrl { get; private set; } = "https://example.com/launcher/latest";
    public string LauncherExeName { get; private set; } = "CubenetLauncher";
    public string AppVersion { get; private set; } = "1.0.0";

    public static EnvConfig Load(string path)
    {
        var config = new EnvConfig();

        if (!File.Exists(path))
            return config;

        foreach (var line in File.ReadAllLines(path))
        {
            var trimmed = line.Trim();
            if (trimmed.Length == 0 || trimmed.StartsWith('#'))
                continue;

            var eq = trimmed.IndexOf('=');
            if (eq < 0) continue;

            var key = trimmed[..eq].Trim();
            var val = trimmed[(eq + 1)..].Trim().Trim('"', '\'');

            switch (key)
            {
                case "APP_NAME":
                    config.AppName = val;
                    break;
                case "LAUNCHER_DOWNLOAD_URL":
                    config.LauncherDownloadUrl = val;
                    break;
                case "LAUNCHER_EXE_NAME":
                    config.LauncherExeName = val;
                    break;
                case "APP_VERSION":
                    config.AppVersion = val;
                    break;
            }
        }

        return config;
    }

    public string GetInstallDir()
    {
        var appDir = $".{AppName}";
        var baseDir = OperatingSystem.IsWindows()
            ? Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData), appDir)
            : Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), appDir);

        return Path.Combine(baseDir, "bin");
    }

    public string GetLauncherPath()
    {
        var dir = GetInstallDir();
        var name = OperatingSystem.IsWindows()
            ? $"{LauncherExeName}.exe"
            : LauncherExeName;

        return Path.Combine(dir, name);
    }
}
