import React, { useState, useRef, useEffect } from "react";
import { Bell, User, Search, Sun, Moon, Settings, LogOut } from "lucide-react";
import Logo from "../assets/img/logsafepass.png";

export default function Header({ darkMode, toggleDarkMode, onSearch }) {
  const [q, setQ] = useState("");
  const [notificationsOpen, setNotificationsOpen] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const notifRef = useRef(null);
  const userRef = useRef(null);

  const handleSearch = (e) => {
    e.preventDefault();
    if (onSearch) onSearch(q);
  };

  // Закрытие уведомлений и меню пользователя при клике вне блока
  useEffect(() => {
    const handleClickOutside = (e) => {
      if (notifRef.current && !notifRef.current.contains(e.target)) {
        setNotificationsOpen(false);
      }
      if (userRef.current && !userRef.current.contains(e.target)) {
        setUserMenuOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const notifications = [
    { id: 1, text: "Нарушение безопасности на Cam-3", time: "5 мин назад" },
    { id: 2, text: "Сотрудник вошел без каски", time: "20 мин назад" },
    { id: 3, text: "Системное сообщение", time: "1 час назад" },
  ];

  return (
    <header
      className={`flex items-center justify-between px-6 py-3 border-b transition-all 
        ${darkMode ? "bg-[#141919] border-white/10 text-white" : "bg-white border-gray-200 text-gray-900"}
        sticky top-0 z-50`}
    >
      {/* Центральная часть — поиск */}
      <form
        onSubmit={handleSearch}
        className="flex-1 max-w-xl mx-6"
        role="search"
        aria-label="Поиск по панели"
      >
        <div
          className={`flex items-center gap-2 px-4 py-2 rounded-full border transition-all
          ${darkMode ? "bg-[#0F1514] border-white/10" : "bg-gray-100 border-gray-200"}
          `}
        >
          <Search size={18} className="opacity-70" />
          <input
            type="text"
            placeholder="Поиск..."
            value={q}
            onChange={(e) => setQ(e.target.value)}
            className={`flex-1 bg-transparent outline-none text-sm ${
              darkMode ? "text-white placeholder-gray-400" : "text-gray-800 placeholder-gray-500"
            }`}
          />
        </div>
      </form>

      {/* Правая часть */}
      <div className="flex items-center gap-4 relative">
        {/* Кнопка уведомлений */}
        <div ref={notifRef} className="relative">
          <button
            aria-label="Уведомления"
            onClick={() => setNotificationsOpen(!notificationsOpen)}
            className={`p-2 rounded-lg transition ${
              darkMode ? "hover:bg-white/10" : "hover:bg-gray-200"
            }`}
          >
            <Bell size={20} />
          </button>
          {notificationsOpen && (
            <div
              className={`absolute right-0 mt-2 w-64 rounded-xl shadow-lg overflow-hidden z-50 transition-all
              ${darkMode ? "bg-[#1C2422] text-white border border-[#2A3331]" : "bg-white text-gray-900 border border-gray-200"}`}
            >
              <div className="p-3 text-sm font-medium border-b border-gray-200 dark:border-[#2A3331]">
                Уведомления
              </div>
              <div className="max-h-64 overflow-y-auto">
                {notifications.map((n) => (
                  <div key={n.id} className="px-3 py-2 hover:bg-gray-100 dark:hover:bg-[#232B29] cursor-pointer">
                    <div className="text-sm">{n.text}</div>
                    <div className="text-xs opacity-60">{n.time}</div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Переключатель темы */}
        <button
          aria-label="Сменить тему"
          onClick={toggleDarkMode}
          className={`p-2 rounded-lg transition ${
            darkMode ? "hover:bg-white/10" : "hover:bg-gray-200"
          }`}
        >
          {darkMode ? <Sun size={18} /> : <Moon size={18} />}
        </button>

        {/* Блок пользователя */}
        <div ref={userRef} className="relative">
          <button
            onClick={() => setUserMenuOpen(!userMenuOpen)}
            className={`flex items-center gap-2 px-3 py-1 rounded-full border transition ${
              darkMode ? "bg-[#0F1514] border-white/10" : "bg-gray-100 border-gray-200"
            }`}
          >
            <div className="w-8 h-8 rounded-full bg-gray-400 flex items-center justify-center text-white font-medium">
              W
            </div>
            <div className="text-left leading-tight hidden sm:block">
              <div className="text-sm font-medium">Павел Романов</div>
              <div className="text-xs opacity-60">Admin</div>
            </div>
          </button>

          {/* Выпадающее меню пользователя */}
          {userMenuOpen && (
            <div
              className={`absolute right-0 mt-2 w-48 rounded-xl shadow-lg overflow-hidden z-50 transition-all
              ${darkMode ? "bg-[#1C2422] text-white border border-[#2A3331]" : "bg-white text-gray-900 border border-gray-200"}`}
            >
              <div className="py-2">
                <button className="flex items-center gap-2 w-full px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-[#232B29] transition">
                  <User size={16} /> Профиль
                </button>
                <button className="flex items-center gap-2 w-full px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-[#232B29] transition">
                  <Settings size={16} /> Настройки
                </button>
                <button className="flex items-center gap-2 w-full px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-[#232B29] transition">
                  <LogOut size={16} /> Выйти
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
