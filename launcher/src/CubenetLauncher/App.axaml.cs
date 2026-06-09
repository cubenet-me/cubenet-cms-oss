using Avalonia;
using Avalonia.Controls.ApplicationLifetimes;
using Avalonia.Markup.Xaml;
using CubenetLauncher.ViewModels;
using CubenetLauncher.Views;
using System.Threading.Tasks;

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

            _ = SimulateLoadingAsync(vm, loading, main);
        }
        else
        {
            Logger.Warn("Not a desktop application lifetime");
        }

        base.OnFrameworkInitializationCompleted();
    }

    private static async Task SimulateLoadingAsync(
        MainWindowViewModel vm,
        LoadingWindow loading,
        MainWindow main)
    {
        try
        {
            Logger.Info("Loading phase: initialization");
            vm.StatusText = "Инициализация...";
            await Task.Delay(1500);
            vm.ProgressValue = 30;

            Logger.Info("Loading phase: resources");
            vm.StatusText = "Загрузка ресурсов...";
            await Task.Delay(1000);
            vm.ProgressValue = 60;

            Logger.Info("Loading phase: preparation");
            vm.StatusText = "Подготовка...";
            await Task.Delay(1000);
            vm.ProgressValue = 90;

            await Task.Delay(500);
            vm.ProgressValue = 100;
            vm.IsLoading = false;

            if (Current?.ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
            {
                Logger.Info("Switching to main window");
                desktop.MainWindow = main;
                main.Show();
                loading.Close();
            }
        }
        catch (Exception ex)
        {
            Logger.Error($"Loading failed: {ex}");
        }
    }
}
