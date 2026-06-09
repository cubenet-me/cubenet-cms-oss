using Avalonia.Controls;
using System.Diagnostics;

namespace CubenetLauncher.Bootstrapper;

public partial class BootstrapperWindow : Window
{
    private readonly EnvConfig _config;

    public BootstrapperWindow()
    {
        InitializeComponent();

        _config = LoadConfig();
        Logger.Info($"Bootstrapper started — app: {_config.AppName}, ver: {_config.AppVersion}");
        Logger.Info($"Install dir: {_config.GetInstallDir()}");
        Logger.Info($"Launcher path: {_config.GetLauncherPath()}");

        Title = $"{_config.AppName} Bootstrapper";
        _ = RunAsync();
    }

    private static EnvConfig LoadConfig()
    {
        var candidates = new[]
        {
            Path.Combine(AppContext.BaseDirectory, ".env"),
            Path.Combine(AppContext.BaseDirectory, "..", "..", "..", "..", "..", ".env"),
            Path.Combine(AppContext.BaseDirectory, "..", "..", "..", "..", ".env"),
        };

        foreach (var path in candidates)
        {
            var full = Path.GetFullPath(path);
            if (File.Exists(full))
            {
                Logger.Info($"Loading config from {full}");
                return EnvConfig.Load(full);
            }
        }

        Logger.Warn(".env not found, using defaults");
        return new EnvConfig();
    }

    private async Task RunAsync()
    {
        try
        {
            StatusText.Text = "Проверка обновлений...";
            Logger.Info("Checking for updates");
            await Task.Delay(1000);

            StatusText.Text = "Загрузка лаунчера...";
            Logger.Info("Downloading launcher");
            ProgressBar.IsIndeterminate = false;
            ProgressBar.Value = 30;
            await Task.Delay(1500);
            ProgressBar.Value = 60;
            await Task.Delay(1000);
            ProgressBar.Value = 90;
            await Task.Delay(500);
            ProgressBar.Value = 100;

            var launcherPath = _config.GetLauncherPath();
            Logger.Info($"Target launcher path: {launcherPath}");

            StatusText.Text = "Запуск...";

            if (File.Exists(launcherPath))
            {
                new Process
                {
                    StartInfo = new ProcessStartInfo
                    {
                        FileName = launcherPath,
                        UseShellExecute = true,
                        WorkingDirectory = Path.GetDirectoryName(launcherPath),
                    }
                }.Start();
                Logger.Info("Launcher process started");
            }
            else
            {
                Logger.Warn($"Launcher not found at {launcherPath}");
                StatusText.Text = "Лаунчер не найден";
                await Task.Delay(2000);
            }
        }
        catch (Exception ex)
        {
            Logger.Error($"Bootstrapper failed: {ex}");
            StatusText.Text = "Ошибка";
            await Task.Delay(2000);
        }

        Close();
    }
}
