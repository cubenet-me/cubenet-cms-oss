using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using CubenetLauncher.Models;
using CubenetLauncher.Services;

namespace CubenetLauncher.ViewModels;

public partial class SettingsViewModel : ViewModelBase
{
    [ObservableProperty]
    private int _ramMb;

    [ObservableProperty]
    private string _javaPath = string.Empty;

    [ObservableProperty]
    private string _gameDirectory = string.Empty;

    public SettingsViewModel()
    {
        var settings = SettingsService.Load();
        _ramMb = settings.RamMb;
        _javaPath = settings.JavaPath;
        _gameDirectory = settings.GameDirectory;
    }

    [RelayCommand]
    private void Save()
    {
        var settings = new Settings
        {
            RamMb = RamMb,
            JavaPath = JavaPath,
            GameDirectory = GameDirectory,
        };
        SettingsService.Save(settings);
        Logger.Info("Settings saved");
    }
}
