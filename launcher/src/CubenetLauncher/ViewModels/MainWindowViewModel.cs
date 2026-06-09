using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;

namespace CubenetLauncher.ViewModels;

public partial class MainWindowViewModel : ViewModelBase
{
    [ObservableProperty]
    private string _statusText = "Инициализация...";

    [ObservableProperty]
    private double _progressValue;

    [ObservableProperty]
    private bool _isLoading = true;

    [ObservableProperty]
    private ViewModelBase _currentPage = null!;

    public HomeViewModel HomeVM { get; } = new();
    public SettingsViewModel SettingsVM { get; } = new();

    public string HomeButtonBg => CurrentPage is HomeViewModel ? "#008C45" : "Transparent";
    public string SettingsButtonBg => CurrentPage is SettingsViewModel ? "#008C45" : "Transparent";

    public MainWindowViewModel()
    {
        _currentPage = HomeVM;
    }

    partial void OnCurrentPageChanged(ViewModelBase value)
    {
        OnPropertyChanged(nameof(HomeButtonBg));
        OnPropertyChanged(nameof(SettingsButtonBg));
    }

    [RelayCommand]
    private void GoToHome()
    {
        CurrentPage = HomeVM;
    }

    [RelayCommand]
    private void GoToSettings()
    {
        CurrentPage = SettingsVM;
    }
}
