using Avalonia.Controls;
using System.Diagnostics;

namespace CubenetLauncher.Bootstrapper;

public partial class BootstrapperWindow : Window
{
    public BootstrapperWindow()
    {
        InitializeComponent();

        Logger.Info("Bootstrapper started");
        _ = RunAsync();
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

            StatusText.Text = "Запуск...";
            Logger.Info("Launching CubenetLauncher");

            var launcherPath = Path.Combine(
                AppContext.BaseDirectory,
                "CubenetLauncher");

            if (OperatingSystem.IsWindows())
                launcherPath += ".exe";

            if (File.Exists(launcherPath))
            {
                new Process
                {
                    StartInfo = new ProcessStartInfo
                    {
                        FileName = launcherPath,
                        UseShellExecute = true,
                    }
                }.Start();
                Logger.Info("Launcher process started");
            }
            else
            {
                Logger.Warn($"Launcher not found at {launcherPath}");
            }
        }
        catch (Exception ex)
        {
            Logger.Error($"Bootstrapper failed: {ex}");
        }

        Close();
    }
}
