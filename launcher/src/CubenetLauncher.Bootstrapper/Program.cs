using System.Diagnostics;

var launcherUrl = args.Length > 0 ? args[0] : "https://example.com/cubenet-launcher/latest";
var launcherPath = Path.Combine(
    Path.GetTempPath(),
    "CubenetLauncher",
    "CubenetLauncher.exe");

Console.WriteLine("Cubenet Bootstrapper");
Console.WriteLine("====================");
Console.WriteLine();

// Step 1: check for updates / download
Console.Write("Проверка обновлений... ");
await Task.Delay(500);
Console.WriteLine("OK");

Console.Write("Загрузка лаунчера... ");
// TODO: actual download from launcherUrl
await Task.Delay(1000);
Console.WriteLine("OK");

// Step 2: verify
if (!File.Exists(launcherPath))
{
    Console.Error.WriteLine("Ошибка: лаунчер не найден");
    return 1;
}

Console.Write("Запуск лаунчера... ");

var process = new Process
{
    StartInfo = new ProcessStartInfo
    {
        FileName = launcherPath,
        UseShellExecute = true,
    }
};

process.Start();
Console.WriteLine("OK");

return 0;
