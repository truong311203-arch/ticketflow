import { useEffect, useState } from "react";
import { Navbar } from "../components/Navbar";
import { useAuth } from "../context/AuthContext";
import { MOCK_EVENTS } from "../data/mockEvents";
import type { EventItem } from "../types";

export function HomePage() {
  const { session } = useAuth();
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  useEffect(() => {
    if (!toastMsg) return;
    const timer = window.setTimeout(() => setToastMsg(null), 3000);
    return () => window.clearTimeout(timer);
  }, [toastMsg]);

  const handleSelectEvent = (evt: EventItem) => {
    setToastMsg(`[Mock] Đã chọn "${evt.title}". Sẽ gắn API thật ở bước sau!`);
  };

  return (
    <div>
      <Navbar />

      <main className="container page-content">
        <section className="hero">
          <div>
            <h2>Khám phá sự kiện nổi bật</h2>
            <p>
              Danh sách dưới đây đang dùng <strong>Mock Data</strong>. Sau khi
              hoàn thiện Auth &amp; JWT ở Go Backend, trang này sẽ gọi API với
              header <code>Authorization: Bearer &lt;token&gt;</code>.
            </p>
          </div>
          <div className="session-box">
            <div className="session-label">Mock Token hiện tại:</div>
            <code className="token-preview">{session?.token}</code>
          </div>
        </section>

        <section>
          <div className="section-header">
            <h3>Sự kiện đang mở bán</h3>
            <span className="count-text">{MOCK_EVENTS.length} sự kiện</span>
          </div>

          <div className="events-grid">
            {MOCK_EVENTS.map((evt) => {
              const isSoldOut = evt.status === "SOLD_OUT";
              return (
                <article key={evt.id} className="event-card">
                  <div>
                    <span
                      className={`event-status ${
                        isSoldOut ? "status-sold-out" : "status-on-sale"
                      }`}
                    >
                      {isSoldOut ? "Hết vé" : `Còn ${evt.availableSeats} ghế`}
                    </span>
                    <h4 className="event-title">{evt.title}</h4>
                    <ul className="event-meta">
                      <li>Địa điểm: {evt.venue}</li>
                      <li>Thời gian: {evt.date}</li>
                    </ul>
                  </div>

                  <div className="event-footer">
                    <span className="event-price">{evt.price}</span>
                    <button
                      type="button"
                      className="btn btn-primary"
                      disabled={isSoldOut}
                      onClick={() => handleSelectEvent(evt)}
                    >
                      {isSoldOut ? "Đã hết chỗ" : "Chọn ghế"}
                    </button>
                  </div>
                </article>
              );
            })}
          </div>
        </section>
      </main>

      {toastMsg && (
        <div className="toast" role="status" aria-live="polite">
          {toastMsg}
        </div>
      )}
    </div>
  );
}
