# GoCheck

GoCheck là một công cụ kiểm tra mã nguồn Go tự động, giúp phát hiện các vấn đề về clean code, hiệu năng và bảo mật trong dự án của bạn. Kết quả được xuất ra dưới dạng báo cáo HTML và JSON.

## Tính năng
- **Quét toàn bộ thư mục**: Tự động tìm tất cả file `.go` trong thư mục chỉ định.
- **Phân tích Clean Code**: Phát hiện các hàm quá dài, gợi ý tách nhỏ để dễ bảo trì.
- **Phân tích Hiệu năng**: Cảnh báo các vòng lặp for có thể ảnh hưởng đến hiệu năng.
- **Phân tích Bảo mật**: Phát hiện hardcode mật khẩu, API key trong mã nguồn.
- **Báo cáo HTML & JSON**: Xuất kết quả ra file `report.html` và `report.json`.

## Cài đặt
Yêu cầu Go >= 1.23.2

```bash
git clone <repo-url>
cd gocheck
go mod tidy
```

## Sử dụng
```bash
go run main.go --path <thư_mục_cần_quét> [--html=true|false] [--json=true|false]
```
- `--path`: Đường dẫn thư mục cần quét (mặc định là thư mục hiện tại)
- `--html`: Xuất báo cáo HTML (mặc định: true)
- `--json`: Xuất báo cáo JSON (mặc định: true)

Sau khi chạy, bạn sẽ nhận được các file `report.html` và/hoặc `report.json` trong thư mục hiện tại.

## Cấu trúc dự án
```
main.go           // Điểm vào chương trình, xử lý tham số dòng lệnh
analyzer/         // Phân tích clean code, hiệu năng, bảo mật
  |- analyzer.go  // Hàm tổng hợp phân tích
  |- cleancode.go // Phân tích clean code
  |- performance.go // Phân tích hiệu năng
  |- security.go  // Phân tích bảo mật
  |- finding.go   // Định nghĩa Finding (output)
scanner/          // Quét file Go trong thư mục
report/           // Sinh báo cáo HTML, JSON
utils/            // Hàm tiện ích
```

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