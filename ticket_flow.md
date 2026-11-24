# TicketFlow — MVP Product & Technical Specification

**Version:** 1.0  
**Project type:** Backend Engineering Portfolio Project  
**Status:** Initial Specification  
**Primary focus:** Backend Engineering, Database Design, Concurrency Control, Transaction Management

---

# 1. Project Overview

## 1.1. Product Description

TicketFlow là một nền tảng đặt vé sự kiện trực tuyến, cho phép người dùng khám phá sự kiện, xem sơ đồ ghế, lựa chọn ghế, giữ chỗ tạm thời, thanh toán và nhận vé điện tử.

Hệ thống cung cấp thêm chức năng quản trị để Admin tạo và quản lý địa điểm tổ chức, sơ đồ ghế và các sự kiện.

Đây là một dự án tập trung vào Backend Engineering. Frontend chỉ cần đủ để minh họa và kiểm thử các luồng nghiệp vụ.

## 1.2. Problem Statement

Một nền tảng booking cần giải quyết các vấn đề nghiệp vụ quan trọng:

- Nhiều người dùng có thể cùng muốn đặt một ghế tại cùng thời điểm.
- Hệ thống phải tránh tình trạng một ghế được bán cho nhiều người (double booking).
- Ghế được giữ tạm thời nhưng không thanh toán phải được giải phóng sau một khoảng thời gian.
- Trạng thái booking, payment và ticket phải nhất quán.
- Các thao tác liên quan đến nhiều ghế cần đảm bảo tính toàn vẹn dữ liệu.
- Hệ thống phải xử lý được các trường hợp request lặp lại, thanh toán thất bại và các thao tác đồng thời.

TicketFlow được xây dựng để nghiên cứu và giải quyết những bài toán này trong một ứng dụng Backend thực tế.

## 1.3. Project Goals

Mục tiêu chính không phải xây dựng một hệ thống có quy mô lớn ngay từ đầu, mà là hiểu và triển khai đúng các nền tảng Backend Engineering.

Các mục tiêu học tập:

1. Thiết kế RESTful API và tổ chức một Go Backend application.
2. Thiết kế relational database và hiểu cách PostgreSQL quản lý dữ liệu.
3. Hiểu và áp dụng database transaction, ACID, MVCC, row-level locking, isolation levels.
4. Giải quyết race conditions và concurrency trong nghiệp vụ booking.
5. Xây dựng background jobs để xử lý các tác vụ định kỳ.
6. Thiết kế authentication, authorization và validation.
7. Viết unit tests, integration tests và concurrency tests.
8. Sử dụng Docker để đóng gói và chạy hệ thống.
9. Xây dựng nền tảng để mở rộng về sau với Redis, Elasticsearch, Kafka và các kỹ thuật production engineering khác.

---

# 2. Actors & User Roles

Hệ thống có hai nhóm người dùng chính.

## 2.1. Customer

Customer là người sử dụng nền tảng để đặt vé.

Các quyền chính:

- Đăng ký và đăng nhập.
- Xem danh sách sự kiện.
- Tìm kiếm và lọc sự kiện.
- Xem thông tin chi tiết sự kiện.
- Xem sơ đồ ghế và trạng thái ghế.
- Giữ một hoặc nhiều ghế.
- Thanh toán booking.
- Xem lịch sử booking.
- Xem vé điện tử đã được phát hành.
- Hủy booking chưa thanh toán.

Customer chỉ được xem và thao tác trên những booking thuộc quyền sở hữu của mình.

## 2.2. Admin

Admin là người quản lý dữ liệu và hoạt động của nền tảng.

Các quyền chính:

- Quản lý địa điểm tổ chức (Venue).
- Quản lý danh sách ghế của từng Venue.
- Tạo và cập nhật Event.
- Thiết lập giá vé cho từng ghế trong một Event.
- Đưa Event vào trạng thái công khai hoặc ngừng bán.

Admin có quyền quản trị dữ liệu sự kiện, không mặc định có quyền can thiệp trực tiếp vào các booking hoặc payment đã phát sinh.

---

# 3. Functional Requirements

## 3.1. Authentication & Authorization

### FR-01: User Registration

- User có thể đăng ký bằng email và password.
- Email phải là duy nhất.
- Password phải được hash trước khi lưu vào database.
- User mới mặc định có role CUSTOMER.

### FR-02: User Login

- User đăng nhập bằng email và password.
- Backend xác thực thông tin và trả về JWT.
- Các API yêu cầu xác thực phải kiểm tra JWT.

### FR-03: Role-Based Authorization

Hệ thống hỗ trợ hai role:

- CUSTOMER
- ADMIN

Backend phải kiểm tra role khi truy cập các API quản trị.

Không được chỉ dựa vào việc ẩn hoặc hiện chức năng trên frontend để đảm bảo authorization.

## 3.2. Venue Management

### FR-04: Create Venue

Admin có thể tạo Venue với các thông tin:

- Name
- Location
- Description (optional)

### FR-05: Manage Venue Seats

Admin có thể tạo và quản lý danh sách ghế trong một Venue.

Thông tin cơ bản của một ghế:

- Section
- Row
- Seat number
- Position coordinates (optional, phục vụ hiển thị seat map)

Một Venue có thể chứa nhiều VenueSeat.

Lưu ý: VenueSeat mô tả ghế vật lý, không chứa trạng thái bán vé của một sự kiện cụ thể.

## 3.3. Event Management

### FR-06: Create Event

Admin có thể tạo Event với:

- Name
- Description
- Venue
- Category
- Start time
- End time (optional)
- Sales start time
- Sales end time
- Status

Một Event thuộc về một Venue.

### FR-07: Initialize Event Inventory

Khi tạo Event, hệ thống tạo danh sách EventSeat dựa trên các VenueSeat thuộc Venue tương ứng.

Mỗi EventSeat có:

- VenueSeat reference
- Event reference
- Price
- Availability status

Trạng thái ban đầu của tất cả EventSeat là AVAILABLE.

Mỗi ghế vật lý có một trạng thái bán vé riêng trong từng Event.

### FR-08: Update Event

Admin có thể cập nhật thông tin Event khi phù hợp với các quy tắc nghiệp vụ.

Đối với MVP, không hỗ trợ thay đổi cấu trúc venue hoặc inventory sau khi đã phát sinh booking, trừ khi có quy trình xử lý riêng.

### FR-09: Publish Event

Event chỉ được hiển thị cho Customer khi có trạng thái PUBLISHED và nằm trong thời gian được phép bán vé.

## 3.4. Event Discovery

### FR-10: List Events

Customer có thể xem danh sách Event đang được công khai.

Hệ thống hỗ trợ:

- Pagination
- Search theo tên
- Filter theo location
- Filter theo category
- Filter theo ngày diễn ra
- Sort theo thời gian diễn ra

MVP sử dụng PostgreSQL để thực hiện tìm kiếm và lọc.

Không yêu cầu Elasticsearch.

### FR-11: Event Detail

Customer có thể xem:

- Event information
- Venue information
- Performer information nếu có
- Event time
- Ticket availability overview

Performer có thể được giản lược hoặc lưu dưới dạng thông tin mô tả trong MVP. Chưa yêu cầu xây dựng một hệ thống quản lý Performer riêng.

### FR-12: Seat Map

Customer có thể xem sơ đồ ghế của Event.

Mỗi ghế có một trong các trạng thái:

- AVAILABLE
- HELD
- SOLD

Thông tin trả về bao gồm vị trí ghế, giá và trạng thái hiện tại.

Client có thể refresh seat map để nhận trạng thái mới. Real-time seat updates bằng WebSocket chưa thuộc MVP.

## 3.5. Booking & Seat Reservation

Đây là domain quan trọng nhất của hệ thống.

### FR-13: Create Seat Hold

Customer có thể chọn một hoặc nhiều ghế trong cùng một Event để giữ chỗ.

Request phải thỏa mãn:

- Event tồn tại và đang mở bán.
- Tất cả ghế thuộc Event tương ứng.
- Tất cả ghế đều đang AVAILABLE hoặc có hold đã hết hạn.
- Số lượng ghế không vượt quá giới hạn mỗi booking.

Các ghế phải được giữ thành công đồng thời. Nếu một ghế không hợp lệ hoặc không khả dụng, toàn bộ thao tác phải thất bại.

Khi giữ ghế thành công:

- Tạo một Booking.
- Booking có trạng thái PENDING_PAYMENT.
- Các EventSeat tương ứng chuyển sang HELD.
- Lưu thông tin `hold_expires_at`.
- Tạo BookingItem cho từng ghế.
- Lưu snapshot giá vé tại thời điểm booking.

**Default hold duration: 10 minutes.**

Thời gian giữ ghế phải có thể cấu hình.

### FR-14: Concurrency Control

Hệ thống phải đảm bảo một EventSeat không thể được giữ hoặc bán thành công cho nhiều Booking cùng một thời điểm.

PostgreSQL là cơ chế bảo vệ tính toàn vẹn cuối cùng.

MVP ưu tiên sử dụng:

- Database transaction.
- Row-level locking (`SELECT ... FOR UPDATE`) hoặc atomic conditional updates.
- Database constraints.
- Kiểm tra kết quả cập nhật dữ liệu.

Redis distributed lock không thuộc MVP.

### FR-15: Booking Idempotency

Các request tạo booking có thể được gửi lại do timeout hoặc lỗi mạng.

Backend cần hỗ trợ idempotency cho thao tác tạo hold thông qua một idempotency key.

Các request lặp lại với cùng user và cùng idempotency key không được tạo ra nhiều booking khác nhau.

### FR-16: Booking Status

Booking hỗ trợ các trạng thái:

| Status | Meaning |
|---|---|
| PENDING_PAYMENT | Đang giữ ghế, chờ thanh toán |
| CONFIRMED | Thanh toán thành công, booking đã xác nhận |
| EXPIRED | Hold hết hạn |
| CANCELLED | User chủ động hủy booking chưa thanh toán |
| PAYMENT_FAILED | Thanh toán thất bại |

Các thay đổi trạng thái phải được kiểm tra theo state transition hợp lệ.

### FR-17: Cancel Booking

Customer có thể hủy booking ở trạng thái PENDING_PAYMENT.

Khi hủy:

- Booking chuyển sang CANCELLED.
- Các ghế đang được booking đó giữ chuyển về AVAILABLE.
- Những ghế khác không liên quan không được ảnh hưởng.

MVP không hỗ trợ hủy booking đã thanh toán hoặc hoàn tiền.

### FR-18: Expired Hold Handling

Nếu Customer không thanh toán trong 10 phút:

- Booking phải chuyển sang EXPIRED.
- Các EventSeat thuộc booking phải được giải phóng.
- Ghế trở về AVAILABLE.

Hệ thống sử dụng Background Worker để tìm và xử lý các hold hết hạn.

Việc giải phóng ghế phải được thực hiện an toàn bằng transaction và chỉ cập nhật các hold vẫn còn thuộc booking hết hạn.

Một hold đã hết hạn về mặt thời gian phải được xem là không còn hợp lệ, kể cả khi Background Worker chưa kịp chạy.

Worker không phải cơ chế duy nhất quyết định một ghế có thể được đặt hay không. Booking Module vẫn phải kiểm tra thời hạn hold trong lúc tạo booking mới.

## 3.6. Payment

### FR-19: Mock Payment

MVP sử dụng Mock Payment Provider thay vì tích hợp cổng thanh toán thật.

Payment Module phải được thiết kế thông qua một interface để có thể thay thế implementation về sau.

Các kết quả thanh toán được hỗ trợ:

- SUCCESS
- FAILED

Mock provider có thể mô phỏng các kịch bản thành công và thất bại phục vụ testing.

Client không được phép tự quyết định một giao dịch là SUCCESS chỉ bằng cách gửi trường `success=true`.

### FR-20: Payment Processing

Khi Customer yêu cầu thanh toán:

1. Backend xác thực quyền sở hữu booking.
2. Kiểm tra booking còn trong trạng thái PENDING_PAYMENT.
3. Kiểm tra hold chưa hết hạn.
4. Gọi Mock Payment Provider.
5. Xử lý kết quả thanh toán.
6. Cập nhật Payment và Booking trong database transaction.
7. Nếu thành công, chuyển các ghế từ HELD sang SOLD.
8. Phát hành vé điện tử.
9. Nếu thất bại, chuyển booking sang PAYMENT_FAILED và giải phóng các ghế.

Thao tác thanh toán phải có khả năng xử lý request lặp lại mà không xác nhận booking hoặc phát hành vé nhiều lần.

Lưu ý: Mock Payment không thực hiện giao dịch tiền thật. Quy trình này chưa phải thiết kế hoàn chỉnh cho Stripe hoặc các cổng thanh toán thực tế.

## 3.7. E-ticket

### FR-21: Issue Ticket

Sau khi booking được xác nhận:

- Mỗi BookingItem tạo ra một Ticket.
- Mỗi Ticket có một ticket code duy nhất.
- Ticket chứa dữ liệu để tạo QR code.
- Ticket được liên kết với booking và event seat tương ứng.

Việc phát hành Ticket phải idempotent.

Đối với MVP, backend có thể lưu ticket code và trả về QR payload để frontend hiển thị QR code.

Chưa yêu cầu hệ thống check-in hoặc quét vé tại địa điểm tổ chức.

## 3.8. Booking History

### FR-22: View My Bookings

Customer có thể xem danh sách booking của mình.

Mỗi booking hiển thị:

- Booking ID
- Event
- Selected seats
- Total amount
- Booking status
- Created time
- Expiration time nếu có

### FR-23: View My Tickets

Customer có thể xem các vé điện tử đã được phát hành.

Chỉ Ticket hợp lệ thuộc booking đã CONFIRMED mới được hiển thị như vé có thể sử dụng.

---

# 4. Business Rules

Các quy tắc dưới đây là nền tảng để triển khai Domain Logic và Database Constraints.

## 4.1. Event & Inventory Rules

1. Một Event thuộc về đúng một Venue.
2. Một Venue có nhiều VenueSeat.
3. Một Event có nhiều EventSeat.
4. Một VenueSeat có thể xuất hiện trong nhiều Event khác nhau.
5. Một EventSeat chỉ tồn tại một lần trong một Event.
6. Trạng thái bán vé được quản lý ở EventSeat, không phải VenueSeat.

## 4.2. Booking Rules

1. Một Booking thuộc về một Customer và một Event.
2. Một Booking có thể chứa nhiều BookingItem.
3. Mỗi BookingItem tương ứng với một EventSeat.
4. Tất cả ghế trong một Booking phải được giữ thành công hoặc không ghế nào được giữ.
5. Một EventSeat chỉ có thể thuộc tối đa một hold còn hiệu lực tại một thời điểm.
6. Thời hạn hold mặc định là 10 phút.
7. Giá vé được snapshot tại thời điểm booking.
8. Một booking đã EXPIRED, CANCELLED hoặc PAYMENT_FAILED không thể tiếp tục thanh toán.
9. Một booking đã CONFIRMED không được xác nhận lần thứ hai.

## 4.3. Seat State Rules

Các trạng thái của EventSeat:

AVAILABLE → HELD → SOLD

Các trường hợp giải phóng:

- HELD → AVAILABLE khi hold hết hạn.
- HELD → AVAILABLE khi Customer hủy booking chưa thanh toán.
- HELD → AVAILABLE khi payment thất bại.

Không hỗ trợ SOLD → AVAILABLE trong MVP vì chưa có nghiệp vụ refund.

## 4.4. Payment Rules

1. Payment chỉ được xử lý khi booking còn hiệu lực.
2. Chỉ Payment SUCCESS mới cho phép xác nhận booking.
3. Một booking không được xác nhận nhiều lần.
4. Một booking không được phát hành nhiều bộ vé cho cùng các BookingItem.
5. Payment status và Booking status cần được cập nhật theo một quy trình có thể kiểm soát tính nhất quán.

---

# 5. Data Model — Initial Scope

Các entity dự kiến:

| Entity | Responsibility |
|---|---|
| User | Tài khoản và quyền truy cập |
| Venue | Địa điểm tổ chức |
| VenueSeat | Ghế vật lý thuộc Venue |
| Event | Thông tin sự kiện |
| EventSeat | Inventory và trạng thái ghế của một Event |
| Booking | Đơn đặt vé |
| BookingItem | Các ghế được đặt trong Booking |
| Payment | Thông tin giao dịch thanh toán |
| Ticket | Vé điện tử được phát hành |

## 5.1. Important Relationships

- User 1-N Booking.
- Venue 1-N VenueSeat.
- Venue 1-N Event.
- Event 1-N EventSeat.
- VenueSeat 1-N EventSeat.
- Event 1-N Booking.
- Booking 1-N BookingItem.
- BookingItem N-1 EventSeat.
- Booking 1-N Payment.
- BookingItem 1-0..1 Ticket.

## 5.2. Important Database Constraints

Database cần bảo vệ các invariant quan trọng, không chỉ dựa vào application-level validation.

Các constraint cần xem xét:

- Unique constraint trên `users.email`.
- Unique constraint trên `(event_id, venue_seat_id)` của EventSeat.
- Unique constraint trên `(booking_id, event_seat_id)` của BookingItem.
- Unique constraint trên Ticket code.
- Foreign keys giữa các entity.
- Check constraints cho các giá trị số và trạng thái phù hợp.
- Unique constraint cho idempotency key trong phạm vi user.

Database schema chi tiết, column types, indexes và transaction boundaries sẽ được thiết kế trong tài liệu ERD riêng.

---

# 6. High-Level System Architecture

## 6.1. Architectural Style

MVP sử dụng **Modular Monolith**.

Toàn bộ Backend chạy trong một Go application, nhưng chia thành các module có trách nhiệm nghiệp vụ riêng.

Các module chính:

### API Layer

Trách nhiệm:

- Routing.
- Middleware.
- JWT Authentication.
- RBAC Authorization.
- Request validation.
- Response formatting.
- Error handling.

API Layer không chứa business logic phức tạp.

### Event Module

Trách nhiệm:

- Quản lý Venue.
- Quản lý VenueSeat.
- Quản lý Event.
- Khởi tạo EventSeat inventory.
- Browse và search event.
- Truy xuất seat map.

### Booking Module

Trách nhiệm:

- Tạo hold.
- Kiểm tra availability.
- Concurrency control.
- Booking lifecycle.
- Cancel booking.
- Xử lý expiration.
- Điều phối quá trình booking và payment.

Đây là module có domain logic phức tạp nhất.

### Payment Module

Trách nhiệm:

- Định nghĩa Payment Provider interface.
- Tích hợp Mock Payment Provider.
- Theo dõi trạng thái thanh toán.
- Xử lý kết quả SUCCESS và FAILED.

### Ticket Module

Trách nhiệm:

- Phát hành Ticket.
- Tạo ticket code.
- Cung cấp QR payload.
- Truy xuất Ticket theo booking hoặc user.

### Background Worker

Trách nhiệm:

- Tìm các booking PENDING_PAYMENT đã hết hạn.
- Chuyển booking sang EXPIRED.
- Giải phóng các EventSeat đang được booking đó giữ.

Worker có thể chạy trong cùng codebase và deployment của Go application ở MVP.

## 6.2. Infrastructure

MVP chỉ yêu cầu:

- Go Backend.
- PostgreSQL.
- Docker Compose.

PostgreSQL là nguồn dữ liệu chính và là nơi đảm bảo tính toàn vẹn của booking.

Không sử dụng Redis, Elasticsearch, Kafka hoặc microservices trong MVP.

---

# 7. API Scope

Đây là danh sách endpoint dự kiến, chưa phải API contract cuối cùng.

## 7.1. Authentication

| Method | Endpoint | Purpose |
|---|---|---|
| POST | /auth/register | Register |
| POST | /auth/login | Login |

## 7.2. Event Discovery

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /events | Browse and search events |
| GET | /events/{eventId} | Event details |
| GET | /events/{eventId}/seats | Event seat map |

Query filters dự kiến:

- `term`
- `location`
- `category`
- `date`
- `page`
- `limit`
- `sort`

## 7.3. Admin Management

| Method | Endpoint | Purpose |
|---|---|---|
| POST | /admin/venues | Create venue |
| POST | /admin/venues/{venueId}/seats | Create venue seats |
| POST | /admin/events | Create event |
| PATCH | /admin/events/{eventId} | Update event |
| PATCH | /admin/events/{eventId}/publish | Publish event |

## 7.4. Booking

| Method | Endpoint | Purpose |
|---|---|---|
| POST | /bookings/holds | Hold selected seats |
| GET | /bookings/{bookingId} | Booking details |
| POST | /bookings/{bookingId}/payment | Process mock payment |
| POST | /bookings/{bookingId}/cancel | Cancel pending booking |
| GET | /me/bookings | User booking history |
| GET | /me/tickets | User tickets |

Các API tạo hold và xử lý thanh toán cần có cơ chế idempotency phù hợp.

## 7.5. API Design Requirements

- Sử dụng HTTP status codes phù hợp.
- Request và response sử dụng JSON.
- Có validation rõ ràng.
- API lỗi phải trả về error response nhất quán.
- Không trả về internal database errors cho client.
- Kiểm tra ownership và authorization ở phía Backend.

Một số status code dự kiến:

| Status | Meaning |
|---|---|
| 200 | Request successful |
| 201 | Resource created |
| 400 | Invalid request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Resource not found |
| 409 | Business conflict, e.g. seat unavailable |
| 500 | Internal server error |

---

# 8. Non-Functional Requirements

## 8.1. Data Consistency

Đây là yêu cầu quan trọng nhất.

Hệ thống phải đảm bảo:

- Không có double booking.
- Không có overselling.
- Một thao tác đặt nhiều ghế là atomic.
- Booking, Payment, EventSeat và Ticket có các trạng thái nhất quán theo quy tắc nghiệp vụ.
- Không phát hành nhiều vé trùng lặp cho một BookingItem.

## 8.2. Concurrency

Hệ thống cần hỗ trợ nhiều request booking được gửi đồng thời.

Acceptance test quan trọng:

- Tạo một EventSeat đang AVAILABLE.
- Gửi 100 request đồng thời từ 100 user khác nhau để giữ cùng ghế.
- Chỉ một request được giữ ghế thành công.
- Các request còn lại nhận conflict.
- Database không tồn tại nhiều hold hợp lệ cho cùng một EventSeat.

Ngoài ra cần test tình huống một booking chứa nhiều ghế nhưng chỉ một ghế còn available. Toàn bộ booking phải thất bại và không được giữ bất kỳ ghế nào.

Đây là mục tiêu correctness, không phải cam kết throughput production.

## 8.3. Security

- Password phải được hash bằng thuật toán phù hợp.
- JWT secret không được hard-code.
- Các API quản trị phải được bảo vệ bằng RBAC.
- Customer không được truy cập booking của user khác.
- Các trường input phải được validate.
- Database queries phải sử dụng parameterized queries.
- Không để lộ secret hoặc thông tin nhạy cảm trong logs.

## 8.4. Maintainability

- Phân chia module theo domain.
- Tránh đặt business logic trong HTTP handlers.
- Repository/data-access layer chịu trách nhiệm truy cập database.
- Có centralized error handling.
- Có cấu hình môi trường tách biệt.
- Code phải dễ viết unit tests và integration tests.

## 8.5. Deployment

- Application chạy được bằng Docker Compose.
- Có PostgreSQL container.
- Có database migrations.
- Có cấu hình môi trường phát triển.
- Có thể khởi tạo hệ thống từ trạng thái database rỗng.

---

# 9. Testing Strategy

## 9.1. Unit Testing

Ưu tiên kiểm thử:

- Booking state transitions.
- Hold expiration calculation.
- Payment result handling.
- Validation rules.
- Authorization rules.

## 9.2. Integration Testing

Kiểm thử với PostgreSQL thực tế:

- Create event và generate inventory.
- Hold seats.
- Cancel booking.
- Expire booking.
- Successful payment.
- Failed payment.
- Ticket issuance.

Các bài test cần kiểm tra trạng thái thực tế trong database, không chỉ kiểm tra HTTP response.

## 9.3. Concurrency Testing

Các trường hợp bắt buộc:

1. Nhiều user cùng hold một ghế.
2. Một user hold nhiều ghế, nhưng một ghế đã được user khác giữ.
3. Hai thao tác booking xảy ra gần như đồng thời.
4. Payment và expiration xảy ra gần cùng thời điểm.
5. Hai request payment cùng booking được gửi lặp lại.
6. Worker xử lý cùng một booking hết hạn nhiều lần.

Kết quả phải đảm bảo không vi phạm các business invariants.

---

# 10. MVP Acceptance Criteria

MVP được xem là hoàn thành khi:

- [ ] User có thể đăng ký và đăng nhập.
- [ ] JWT authentication và RBAC hoạt động.
- [ ] Admin có thể tạo Venue và VenueSeat.
- [ ] Admin có thể tạo Event và EventSeat inventory.
- [ ] Customer có thể browse và search Event.
- [ ] Customer có thể xem seat map.
- [ ] Customer có thể giữ một hoặc nhiều ghế.
- [ ] Booking được tạo trong một transaction.
- [ ] Hệ thống không cho phép double booking.
- [ ] Hold có thời hạn 10 phút.
- [ ] Worker có thể giải phóng hold hết hạn.
- [ ] Customer có thể hủy booking chưa thanh toán.
- [ ] Mock Payment hỗ trợ cả SUCCESS và FAILED.
- [ ] Booking thành công chuyển ghế sang SOLD.
- [ ] Ticket được phát hành idempotently.
- [ ] Customer xem được booking history và e-ticket.
- [ ] Có unit tests, integration tests và concurrency tests.
- [ ] Dự án chạy được bằng Docker Compose.

---

# 11. Explicitly Out of Scope

Các tính năng sau không thuộc MVP:

| Feature | Reason for Deferral |
|---|---|
| Redis | Chưa cần cho correctness ban đầu |
| Elasticsearch | PostgreSQL đủ cho search scope nhỏ |
| Kafka | Chưa cần event streaming |
| CDC / Debezium | Chưa có nhu cầu đồng bộ downstream |
| Microservices | Modular Monolith phù hợp hơn |
| Dedicated API Gateway | Chưa có nhiều backend services |
| CDN | Chưa có yêu cầu phân phối static assets quy mô lớn |
| Stripe integration | Mock Payment trước |
| WebSocket | Chưa cần cập nhật seat map real-time |
| Virtual Waiting Room | Chưa cần xử lý flash-sale quy mô lớn |
| Distributed Lock | PostgreSQL đảm bảo concurrency cho MVP |
| Refund after successful payment | Chưa triển khai nghiệp vụ hoàn tiền |
| Ticket transfer | Chưa triển khai chuyển nhượng vé |
| Promo codes | Chưa nằm trong nghiệp vụ cốt lõi |
| Dynamic pricing | Giá vé cố định theo EventSeat |
| Check-in scanner | Chưa có nghiệp vụ soát vé |
| Notification service | Chưa cần email/SMS delivery |

Không tự động đưa các thành phần trên vào MVP chỉ vì chúng xuất hiện trong các ví dụ System Design.

Chỉ mở rộng khi có một vấn đề kỹ thuật hoặc nghiệp vụ cụ thể cần giải quyết.

---

# 12. Suggested Development Roadmap

Thứ tự dưới đây ưu tiên xây dựng theo từng luồng nghiệp vụ end-to-end.

## Phase 1 — Foundation

- Initialize Go project.
- Set up Docker Compose.
- Configure PostgreSQL.
- Set up migrations.
- Implement application configuration.
- Establish API structure and error handling.
- Implement authentication and authorization.

**Deliverable:** Backend skeleton with working authentication.

## Phase 2 — Event & Inventory

- Design Venue and VenueSeat schema.
- Implement Venue management.
- Design Event and EventSeat schema.
- Implement Event CRUD.
- Generate EventSeat inventory.
- Implement Event listing and seat map APIs.

**Deliverable:** Users can browse events and view seats.

## Phase 3 — Core Booking

- Design Booking and BookingItem.
- Implement seat hold.
- Implement PostgreSQL transaction.
- Implement row-level locking or conditional updates.
- Implement booking expiration fields.
- Implement cancel booking.
- Implement idempotency.

**Deliverable:** Users can safely hold seats.

## Phase 4 — Payment & Ticket

- Define Payment Provider interface.
- Implement Mock Payment Provider.
- Implement payment state transitions.
- Implement booking confirmation.
- Implement Ticket issuance.
- Implement QR payload.
- Implement booking history and ticket APIs.

**Deliverable:** Complete booking-to-ticket flow.

## Phase 5 — Reliability & Testing

- Implement Background Worker.
- Test expiration behavior.
- Write integration tests.
- Write concurrency tests.
- Test duplicate requests.
- Test failure and recovery scenarios.

**Deliverable:** MVP with verified core business invariants.

## Phase 6 — Production Foundations

- Improve logging.
- Add health checks.
- Improve configuration management.
- Add GitHub Actions CI.
- Add basic load testing with k6.
- Document deployment and local development.

**Deliverable:** Reproducible and demonstrable backend project.

---

# 13. Technology Stack

| Layer | Technology |
|---|---|
| Programming Language | Go |
| HTTP Framework | TBD: Gin or Chi |
| Database | PostgreSQL |
| Database Driver | pgx |
| Database Migration | TBD: golang-migrate or another migration tool |
| Query Access | SQL / pgx, SQLC optional |
| Authentication | JWT |
| Password Hashing | bcrypt or Argon2id |
| Containerization | Docker |
| Local Orchestration | Docker Compose |
| API Testing | Postman / Bruno / curl |
| Automated Testing | Go testing |
| Load Testing | k6 |
| CI | GitHub Actions |

Technology choices that are still TBD should be decided based on project needs rather than added automatically.

---

# 14. Future Evolution (Post-MVP)

After the MVP is complete and tested, the project can evolve in several directions.

## Phase 2: Performance & Caching

- Introduce Redis.
- Cache event details and search results.
- Implement rate limiting.
- Evaluate Redis for temporary reservation coordination.
- Benchmark PostgreSQL-only versus cached reads.

## Phase 3: Search Infrastructure

- Introduce Elasticsearch.
- Implement search indexing.
- Synchronize PostgreSQL changes to Elasticsearch.
- Learn CDC and eventual consistency.
- Consider transactional outbox or Debezium/Kafka.

## Phase 4: Asynchronous Processing

- Introduce Kafka.
- Move background processing toward event-driven architecture.
- Implement consumers and retry handling.
- Study at-least-once delivery and idempotent consumers.
- Study dead-letter queues and transactional outbox.

## Phase 5: Real Payment & Production Reliability

- Integrate Stripe.
- Implement payment webhook verification.
- Handle duplicate webhook events.
- Handle payment/booking reconciliation.
- Implement compensation for payment and reservation failure scenarios.
- Improve observability and load testing.

These are future learning milestones, not MVP requirements.

---

# 15. Engineering Principles & Project Constraints

All future design and implementation decisions should follow these principles:

1. **Correctness before performance.** Avoid double booking before optimizing throughput.
2. **PostgreSQL is the source of truth.** Cache or search indexes must not determine whether a seat can be sold.
3. **Start simple.** Prefer a Modular Monolith over premature microservices.
4. **Understand the problem before choosing technology.** Every infrastructure component must solve an identified problem.
5. **Make state transitions explicit.** Booking, Payment and EventSeat must have clearly defined lifecycle states.
6. **Use database transactions for multi-row business operations.**
7. **Design for retries and concurrency.** Network requests can be duplicated and operations can overlap.
8. **Test business invariants, not only successful API responses.**
9. **Do not add unnecessary complexity to the MVP.**
10. **Keep future extensibility in mind without implementing future requirements prematurely.**

---

# 16. Development Collaboration Context

This document is the project context for future discussions and implementation work on TicketFlow.

When assisting with this project, follow these guidelines:

- Treat this specification as the primary source of truth for MVP scope and business requirements.
- Do not silently expand MVP scope.
- Clearly distinguish between current MVP requirements and post-MVP suggestions.
- Help design the system incrementally, rather than generating the entire application at once.
- Explain the underlying engineering concepts, trade-offs and reasons behind recommendations.
- Prioritize practical Backend Engineering knowledge, especially Go, PostgreSQL, transaction management, concurrency, testing and reliability.
- Challenge incorrect assumptions and point out potential race conditions, inconsistent states and data integrity problems.
- When designing database schemas, explain relationships, constraints, indexes and transaction boundaries.
- When designing APIs, explain endpoint responsibilities, validation, authorization, status codes and idempotency.
- Prefer code that is production-conscious but appropriate for a fresher-level portfolio project.
- Avoid introducing Redis, Kafka, Elasticsearch, microservices or other infrastructure unless a concrete requirement justifies it.
- Maintain consistency with existing architecture and decisions. When a decision needs to change, explain the reason and its impact on the rest of the design.

**The overall goal is not merely to finish a ticket booking application. The goal is to use TicketFlow as a practical learning environment for understanding how a reliable Backend system is designed, implemented, tested and evolved.**