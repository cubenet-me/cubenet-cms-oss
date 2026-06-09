using System.Diagnostics;
using System.Net.Http;

namespace CubenetLauncher.Services;

public class UpdateService
{
    private readonly HttpClient _http = new();

    public async Task<bool> CheckAndUpdateAsync(IProgress<(string status, double progress)> progress)
    {
        var currentPath = Environment.ProcessPath;
        if (string.IsNullOrEmpty(currentPath))
        {
            Logger.Warn("Cannot determine process path");
            return false;
        }

        progress.Report(("Проверка обновлений...", 0));
        Logger.Info($"Checking for updates from {Env.UpdateCheckUrl}");
        await Task.Delay(800);

        // TODO: actual version check
        var hasUpdate = false;

        if (!hasUpdate)
        {
            Logger.Info("No update available");
            progress.Report(("Обновлений нет", 100));
            await Task.Delay(300);
            return false;
        }

        progress.Report(("Загрузка обновления...", 30));
        Logger.Info($"Downloading update from {Env.UpdateDownloadUrl}");

        var newPath = currentPath + ".new";
        var oldPath = currentPath + ".old";

        // TODO: actual download
        await Task.Delay(1500);
        progress.Report(("Обновление загружено", 90));

        progress.Report(("Установка...", 95));

        // Replace executable
        if (File.Exists(oldPath))
            File.Delete(oldPath);

        File.Move(currentPath, oldPath);
        File.Move(newPath, currentPath);

        Logger.Info("Update applied, restarting");

        progress.Report(("Перезапуск...", 100));
        await Task.Delay(200);

        // Restart
        using var proc = new Process();
        proc.StartInfo.FileName = currentPath;
        proc.StartInfo.UseShellExecute = true;
        proc.Start();

        Environment.Exit(0);
        return true;
    }
}
