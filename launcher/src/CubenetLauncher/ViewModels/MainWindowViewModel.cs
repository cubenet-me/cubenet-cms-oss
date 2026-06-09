using CommunityToolkit.Mvvm.ComponentModel;

namespace CubenetLauncher.ViewModels;

public partial class MainWindowViewModel : ViewModelBase
{
    [ObservableProperty]
    private string _statusText = "Инициализация...";

    [ObservableProperty]
    private double _progressValue;

    [ObservableProperty]
    private bool _isLoading = true;
}
