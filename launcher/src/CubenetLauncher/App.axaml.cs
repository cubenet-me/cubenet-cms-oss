using Avalonia;
using Avalonia.Controls.ApplicationLifetimes;
using Avalonia.Markup.Xaml;
using CubenetLauncher.Services;
using CubenetLauncher.ViewModels;
using CubenetLauncher.Views;

namespace CubenetLauncher;

public partial class App : Application
{
    public override void Initialize()
    {
        AvaloniaXamlLoader.Load(this);
    }

    public override void OnFrameworkInitializationCompleted()
    {
        Logger.Info("Launcher started");

        if (ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
        {
            var vm = new MainWindowViewModel();
            var updateService = new UpdateService();

            var loading = new LoadingWindow
            {
                DataContext = vm,
            };

            var main = new MainWindow
            {
                DataContext = vm,
            };

            desktop.MainWindow = loading;
            loading.Show();

            _ = StartAsync(vm, updateService, loading, main);
        }
        else
        {
            Logger.Warn("Not a desktop application lifetime");
        }

        base.OnFrameworkInitializationCompleted();
    }

    private static async Task StartAsync(
        MainWindowViewModel vm,
        UpdateService updateService,
        LoadingWindow loading,
        MainWindow main)
    {
        try
        {
            var updated = await updateService.CheckAndUpdateAsync(
                new Progress<(string status, double progress)>(state =>
                {
                    vm.StatusText = state.status;
                    vm.ProgressValue = state.progress;
                }));

            if (updated)
                return; // app will restart

            // No update — open main window
            vm.StatusText = "Запуск...";
            vm.ProgressValue = 100;
            await Task.Delay(200);
            vm.IsLoading = false;

            if (Current?.ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
            {
                Logger.Info("Opening main window");
                desktop.MainWindow = main;
                main.Show();
                loading.Close();
            }
        }
        catch (Exception ex)
        {
            Logger.Error($"Startup failed: {ex}");
            vm.StatusText = "Ошибка запуска";
            await Task.Delay(2000);
            loading.Close();
        }
    }
}
