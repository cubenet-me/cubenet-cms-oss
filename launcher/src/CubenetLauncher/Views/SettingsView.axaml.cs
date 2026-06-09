using Avalonia.Controls;
using Avalonia.Platform.Storage;
using CubenetLauncher.ViewModels;

namespace CubenetLauncher.Views;

public partial class SettingsView : UserControl
{
    public SettingsView()
    {
        InitializeComponent();
    }

    private async void BrowseJavaPath(object? sender, Avalonia.Interactivity.RoutedEventArgs e)
    {
        var top = TopLevel.GetTopLevel(this);
        if (top is null) return;

        var files = await top.StorageProvider.OpenFilePickerAsync(new FilePickerOpenOptions
        {
            Title = "Выберите Java",
            AllowMultiple = false,
        });

        if (files.Count > 0 && DataContext is SettingsViewModel vm)
            vm.JavaPath = files[0].Path.LocalPath;
    }

    private async void BrowseGameDir(object? sender, Avalonia.Interactivity.RoutedEventArgs e)
    {
        var top = TopLevel.GetTopLevel(this);
        if (top is null) return;

        var dirs = await top.StorageProvider.OpenFolderPickerAsync(new FolderPickerOpenOptions
        {
            Title = "Выберите игровую директорию",
            AllowMultiple = false,
        });

        if (dirs.Count > 0 && DataContext is SettingsViewModel vm)
            vm.GameDirectory = dirs[0].Path.LocalPath;
    }
}
