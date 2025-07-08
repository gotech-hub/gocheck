# GoCheck

GoCheck là một thư viện và công cụ kiểm tra mã nguồn Go tự động, giúp phát hiện các vấn đề về clean code, hiệu năng và bảo mật trong dự án của bạn. Kết quả có thể xuất ra dưới dạng báo cáo HTML và JSON.

## Yêu cầu hệ thống
- Go >= 1.23.2

## Cài đặt
Cài đặt GoCheck vào project của bạn bằng lệnh:

```bash
go get github.com/gotech-hub/gocheck
```

> Thay `github.com/gotech-hub/gocheck` bằng đường dẫn repo thực tế của bạn.

## Cài đặt CLI tool
Bạn có thể cài đặt GoCheck như một công cụ dòng lệnh toàn cục bằng lệnh:

```bash
go install github.com/gotech-hub/gocheck/cmd/gocheck@latest
```

Sau khi cài đặt, bạn có thể chạy lệnh `gocheck` ở bất kỳ đâu (nếu `$GOPATH/bin` hoặc `$GOBIN` đã nằm trong biến môi trường `PATH`).

Ví dụ:
```bash
gocheck --path <thư_mục_cần_quét> --html=true --json=true
```

## Sử dụng như một thư viện (Library)
Import GoCheck vào code của bạn và sử dụng API:

```go
import (
    "github.com/gotech-hub/gocheck/analyzer"
    "github.com/gotech-hub/gocheck/scanner"
)

func main() {
    files, _ := scanner.ScanDir("./path/to/your/code")
    findings := analyzer.Analyze(files)
    for _, f := range findings {
        fmt.Printf("%s:%d %s\n", f.File, f.Line, f.Message)
    }
}
```

## Sử dụng như một công cụ CLI
Bạn cũng có thể chạy GoCheck trực tiếp từ dòng lệnh:

```bash
go run cmd/gocheck/main.go --path <thư_mục_cần_quét> [--html=true|false] [--json=true|false]
```

```bash
go run cmd/gocheck/main.go --path=. [--html=true|false] [--json=true|false]
```

- `--path`: Đường dẫn thư mục cần quét (mặc định là thư mục hiện tại)
- `--html`: Xuất báo cáo HTML (mặc định: true)
- `--json`: Xuất báo cáo JSON (mặc định: true)

Sau khi chạy, bạn sẽ nhận được các file `report.html` và/hoặc `report.json` trong thư mục hiện tại.

## Tính năng
- **Quét toàn bộ thư mục**: Tự động tìm tất cả file `.go` trong thư mục chỉ định.
- **Phân tích Clean Code**: Phát hiện các hàm quá dài, gợi ý tách nhỏ để dễ bảo trì.
- **Phân tích Hiệu năng**: Cảnh báo các vòng lặp for có thể ảnh hưởng đến hiệu năng.
- **Phân tích Bảo mật**: Phát hiện hardcode mật khẩu, API key trong mã nguồn.
- **Báo cáo HTML & JSON**: Xuất kết quả ra file `report.html` và `report.json`.

## API chính
- `scanner.ScanDir(path string) ([]string, error)`: Quét và trả về danh sách file Go trong thư mục.
- `analyzer.Analyze(files []string) []analyzer.Finding`: Phân tích các file và trả về danh sách findings.

## Output mẫu
```json
[
  {
    "file": "main.go",
    "line": 12,
    "message": "Function main is too long (25 lines)",
    "severity": "Medium",
    "suggestion": "Tách hàm ra thành nhiều hàm nhỏ để dễ đọc và test."
  },
  {
    "file": "service.go",
    "line": 30,
    "message": "Hardcoded credential: \"myPassword\"",
    "severity": "High",
    "suggestion": "Không hardcode mật khẩu/API key. Dùng biến môi trường hoặc config file."
  }
]
```

## Đóng góp
Pull request và issue luôn được chào đón! 