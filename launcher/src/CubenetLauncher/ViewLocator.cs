using Avalonia.Controls;
using Avalonia.Controls.Templates;
using CubenetLauncher.ViewModels;
using CubenetLauncher.Views;

namespace CubenetLauncher;

public class ViewLocator : IDataTemplate
{
    public Control? Build(object? param)
    {
        return param switch
        {
            MainWindowViewModel => new MainWindow(),
            _ => new TextBlock { Text = $"Not Found: {param?.GetType().Name}" },
        };
    }

    public bool Match(object? data)
    {
        return data is ViewModelBase;
    }
}
