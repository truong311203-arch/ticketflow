import { useState, type FormEvent } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export function LoginPage() {
  const { isAuthenticated, loginMock } = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  // Nếu đã đăng nhập (có session), tự động chuyển vào Homepage
  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    loginMock(email);
    navigate("/", { replace: true });
  };

  return (
    <div className="auth-page">
      <main className="auth-card">
        <div className="brand">
          <span className="brand-badge">TF</span>
          <div>
            <h1>TicketFlow</h1>
            <p className="subtitle">Online Event Ticketing Platform</p>
          </div>
        </div>

        <div className="mock-banner">
          <strong>MOCK MODE:</strong> Chưa gọi API thật. Bạn có thể nhập email/mật khẩu bất kỳ (hoặc bỏ trống) để bypass vào Homepage.
        </div>

        <form onSubmit={handleSubmit} className="auth-form" noValidate>
          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="demo@ticketflow.dev"
              autoComplete="email"
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">Mật khẩu</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Nhập bất kỳ ký tự nào..."
              autoComplete="current-password"
            />
          </div>

          <button type="submit" className="btn btn-primary btn-full">
            Đăng nhập (Bypass Mock)
          </button>
        </form>

        <p className="auth-footer">
          Bước tiếp theo: Nối tới <code>POST /api/v1/auth/login</code> để nhận JWT thật.
        </p>
      </main>
    </div>
  );
}
