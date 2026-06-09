using Avalonia.Controls;

namespace CubenetLauncher.Views;

public partial class MainWindow : Window
{
    public MainWindow()
    {
        InitializeComponent();
    }

    public static string? ResolveAssetsPath()
    {
        var dir = AppContext.BaseDirectory;

        for (var i = 0; i < 10; i++)
        {
            var candidate = Path.GetFullPath(Path.Combine(dir, "assets"));
            if (Directory.Exists(candidate) && File.Exists(Path.Combine(candidate, "logo.webp")))
                return candidate;

            var parent = Path.GetDirectoryName(dir);
            if (parent is null || parent == dir)
                break;
            dir = parent;
        }

        return null;
    }
}
