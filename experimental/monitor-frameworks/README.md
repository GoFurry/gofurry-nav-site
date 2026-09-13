# Historical cross-framework monitor examples

These Gin, Echo, and Fiber examples intentionally retain `github.com/gofurry/monitor`
through the local `third-party/monitor` submodule. They demonstrate the original
framework-neutral `net/http` integration and remain useful for comparison.
The submodule and its `.gitmodules` entry are retained for this purpose.

Production Nav Backend uses the published `github.com/gofiber/contrib/v3/monitor`
middleware with `Config.Next` to count downstream traffic. That Fiber-specific
API cannot replace the Gin/Echo adapters. These examples are historical references,
not production templates; they remain excluded from active builds, CI, scans,
and deployment. No production module uses their local replacement.
