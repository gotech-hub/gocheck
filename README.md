# GoCheck

[![Go Version](https://img.shields.io/badge/Go-%3E=1.23.2-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
<!-- [![Build Status](https://img.shields.io/github/actions/workflow/status/gotech-hub/gocheck/ci.yml?branch=main)](https://github.com/gotech-hub/gocheck/actions) -->

![Screenshot](https://i.ibb.co/vCJ974SW/Screenshot-2025-07-08-at-16-20-32.png)

> **English below | Tiếng Việt bên dưới**

---

## GoCheck

GoCheck is an automatic static analysis tool and library for Go, helping you detect clean code, performance, and security issues in your project. Reports can be exported in HTML and JSON formats.

---

## Mục lục | Table of Contents
- [Giới thiệu](#giới-thiệu)
- [Yêu cầu hệ thống](#yêu-cầu-hệ-thống)
- [Cài đặt](#cài-đặt)
- [Sử dụng](#sử-dụng)
  - [Dùng như CLI](#dùng-như-cli)
  - [Dùng như thư viện](#dùng-như-thư-viện)
- [Tính năng](#tính-năng)
- [Ví dụ đầu ra](#ví-dụ-đầu-ra)
- [API chính](#api-chính)
- [Đóng góp](#đóng-góp)
- [License](#license)
- [Liên hệ](#liên-hệ)

---

## Giới thiệu
GoCheck là một thư viện và công cụ kiểm tra mã nguồn Go tự động, giúp phát hiện các vấn đề về clean code, hiệu năng và bảo mật trong dự án của bạn. Kết quả có thể xuất ra dưới dạng báo cáo HTML và JSON.

## Yêu cầu hệ thống
- Go >= 1.23.2

## Cài đặt
### Cài đặt vào project:
```bash
go get github.com/gotech-hub/gocheck
```
> Thay `github.com/gotech-hub/gocheck` bằng đường dẫn repo thực tế của bạn nếu cần.

### Cài đặt CLI tool toàn cục:
```bash
go install github.com/gotech-hub/gocheck@latest
GOPROXY=direct GOSUMDB=off go install github.com/gotech-hub/gocheck@v1.0.4
```
Sau khi cài đặt, bạn có thể chạy lệnh `gocheck` ở bất kỳ đâu (nếu `$GOPATH/bin` hoặc `$GOBIN` đã nằm trong biến môi trường `PATH`).

## Prerequisites

Before running `gocheck`, make sure you have the following tools installed:

- [gosec](https://github.com/securego/gosec):
  ```sh
  go install github.com/securego/gosec/v2/cmd/gosec@latest
  ```
- [staticcheck](https://staticcheck.io/):
  ```sh
  go install honnef.co/go/tools/cmd/staticcheck@latest
  ```

## Sử dụng
### Dùng như CLI
Quét mã nguồn và xuất báo cáo:
```bash
gocheck --path <thư_mục_cần_quét> --html=true --json=true
```
Ví dụ:
```bash
gocheck --path=. --html=true --json=true
```
- `--path`: Đường dẫn thư mục cần quét (mặc định là thư mục hiện tại)
- `--html`: Xuất báo cáo HTML (mặc định: true)
- `--json`: Xuất báo cáo JSON (mặc định: true)

Sau khi chạy, bạn sẽ nhận được các file `report.html` và/hoặc `report.json` trong thư mục hiện tại.

### Dùng như thư viện
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

## Tính năng
- **Quét toàn bộ thư mục**: Tự động tìm tất cả file `.go` trong thư mục chỉ định.
- **Phân tích Clean Code**: Phát hiện các hàm quá dài, gợi ý tách nhỏ để dễ bảo trì.
- **Phân tích Hiệu năng**: Cảnh báo các vòng lặp for có thể ảnh hưởng đến hiệu năng.
- **Phân tích Bảo mật**: Phát hiện hardcode mật khẩu, API key trong mã nguồn.
- **Báo cáo HTML & JSON**: Xuất kết quả ra file `report.html` và `report.json`.

## Ví dụ đầu ra
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

## API chính
- `scanner.ScanDir(path string) ([]string, error)`: Quét và trả về danh sách file Go trong thư mục.
- `analyzer.Analyze(files []string) []analyzer.Finding`: Phân tích các file và trả về danh sách findings.


## Đóng góp
Pull request và issue luôn được chào đón! Nếu bạn muốn đóng góp, hãy tạo PR hoặc issue mới. Đọc thêm ở [CONTRIBUTING.md](CONTRIBUTING.md) nếu có.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Liên hệ
- [GitHub Issues](https://github.com/gotech-hub/gocheck/issues)
---

> Made with by GoTech-Hub Team 