import type { EventItem } from "../types";

// Dữ liệu giả lập cho Homepage.
// Sau này sẽ thay bằng lời gọi: GET /api/v1/events (kèm Authorization: Bearer <jwt>)
export const MOCK_EVENTS: EventItem[] = [
  {
    id: "evt-001",
    title: "TicketFlow Live Concert 2026",
    venue: "Nhà thi đấu Phú Thọ, TP.HCM",
    date: "20/12/2026 · 19:30",
    price: "500.000đ",
    availableSeats: 42,
    status: "ON_SALE",
  },
  {
    id: "evt-002",
    title: "Go Backend Engineering Summit",
    venue: "GEM Center, Quận 1, TP.HCM",
    date: "10/01/2027 · 08:30",
    price: "350.000đ",
    availableSeats: 18,
    status: "ON_SALE",
  },
  {
    id: "evt-003",
    title: "Indie Acoustic Night",
    venue: "Cung Văn Hóa Hữu Nghị, Hà Nội",
    date: "15/01/2027 · 20:00",
    price: "400.000đ",
    availableSeats: 0,
    status: "SOLD_OUT",
  },
];
