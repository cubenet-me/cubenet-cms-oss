using System.Text.Json.Serialization;

namespace CubenetLauncher.Services;

public record UpdateInfo(
    [property: JsonPropertyName("url")] string? Url,
    [property: JsonPropertyName("hash")] string? Hash
);
