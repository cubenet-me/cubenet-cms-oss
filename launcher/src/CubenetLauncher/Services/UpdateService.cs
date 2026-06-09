using System.Diagnostics;
using System.Net.Http.Json;
using System.Security.Cryptography;
using System.Text.Json.Nodes;

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
        Logger.Info($"Fetching update info from {Env.UpdateCheckUrl}");

        UpdateInfo? info;
        try
        {
            var json = await _http.GetStringAsync(Env.UpdateCheckUrl);
            var node = JsonNode.Parse(json);
            info = new UpdateInfo(
                node?["url"]?.GetValue<string>(),
                node?["hash"]?.GetValue<string>()
            );
        }
        catch (Exception ex)
        {
            Logger.Warn($"Failed to check updates: {ex.Message}");
            return false;
        }

        if (info is null || string.IsNullOrEmpty(info.Url) || string.IsNullOrEmpty(info.Hash))
        {
            Logger.Warn("Invalid update info response");
            return false;
        }

        // Compare hash
        var currentHash = await Sha256OfFileAsync(currentPath);
        if (string.Equals(currentHash, info.Hash, StringComparison.OrdinalIgnoreCase))
        {
            Logger.Info("Launcher is up to date");
            progress.Report(("Обновлений нет", 100));
            await Task.Delay(300);
            return false;
        }

        Logger.Info($"Update available. Current: {currentHash}, remote: {info.Hash}");

        // Download
        progress.Report(("Загрузка обновления...", 10));
        Logger.Info($"Downloading from {info.Url}");

        var tmpPath = currentPath + ".tmp";

        try
        {
            using var response = await _http.GetAsync(info.Url, HttpCompletionOption.ResponseHeadersRead);
            response.EnsureSuccessStatusCode();

            var total = response.Content.Headers.ContentLength ?? -1;
            await using var src = await response.Content.ReadAsStreamAsync();
            await using var dst = File.Create(tmpPath);

            var buffer = new byte[81920];
            long read = 0;
            int bytes;
            while ((bytes = await src.ReadAsync(buffer)) > 0)
            {
                await dst.WriteAsync(buffer, 0, bytes);
                read += bytes;
                if (total > 0)
                    progress.Report(("Загрузка обновления...", 10 + 80.0 * read / total));
            }

            progress.Report(("Проверка хэша...", 90));
            var downloadedHash = await Sha256OfFileAsync(tmpPath);
            if (!string.Equals(downloadedHash, info.Hash, StringComparison.OrdinalIgnoreCase))
            {
                Logger.Error($"Hash mismatch: expected {info.Hash}, got {downloadedHash}");
                File.Delete(tmpPath);
                return false;
            }
        }
        catch (Exception ex)
        {
            Logger.Error($"Download failed: {ex}");
            if (File.Exists(tmpPath)) File.Delete(tmpPath);
            return false;
        }

        // Replace
        progress.Report(("Установка...", 95));
        Logger.Info("Replacing executable");

        var oldPath = currentPath + ".old";
        if (File.Exists(oldPath)) File.Delete(oldPath);
        File.Move(currentPath, oldPath);
        File.Move(tmpPath, currentPath);

        Logger.Info("Update applied, restarting");

        progress.Report(("Перезапуск...", 100));
        await Task.Delay(200);

        using var proc = new Process();
        proc.StartInfo.FileName = currentPath;
        proc.StartInfo.UseShellExecute = true;
        proc.Start();

        Environment.Exit(0);
        return true;
    }

    private static async Task<string> Sha256OfFileAsync(string path)
    {
        await using var stream = File.OpenRead(path);
        var hash = await SHA256.HashDataAsync(stream);
        return Convert.ToHexStringLower(hash);
    }

    private record UpdateInfo(string? Url, string? Hash);
}
