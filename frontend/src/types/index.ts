export type UserRole = "CUSTOMER" | "ADMIN";

export interface AuthSession {
  email: string;
  role: UserRole;
  token: string;
}

export type EventStatus = "ON_SALE" | "SOLD_OUT";

export interface EventItem {
  id: string;
  title: string;
  venue: string;
  date: string;
  price: string;
  availableSeats: number;
  status: EventStatus;
}
