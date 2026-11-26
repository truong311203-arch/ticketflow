import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export function Navbar() {
  const { session, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login", { replace: true });
  };

  return (
    <header className="navbar">
      <div className="container navbar-inner">
        <div className="brand-inline">
          <span className="brand-badge">TF</span>
          <span className="brand-name">TicketFlow</span>
          <span className="env-pill">MOCK DATA</span>
        </div>

        <div className="user-menu">
          <div className="user-info">
            <span className="user-email">{session?.email}</span>
            <span className="role-badge">{session?.role}</span>
          </div>
          <button type="button" onClick={handleLogout} className="btn btn-outline">
            Đăng xuất
          </button>
        </div>
      </div>
    </header>
  );
}
