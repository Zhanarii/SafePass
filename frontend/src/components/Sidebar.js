import React from "react";
import { Link, useLocation } from "react-router-dom";
import { Video, Users, AlertTriangle, BarChart2, Settings, Home, Moon, Sun } from "lucide-react";
import LogoLight from "../assets/img/logsafepass.png";
import LogoDark from "../assets/img/logoo2.png";

function Sidebar({ darkMode, toggleDarkMode }) {
  const location = useLocation();

  const menuItems = [
    { name: "Главная", path: "/dashboard", icon: <Home size={18} /> },
    { name: "Камеры", path: "/dashboard/camera", icon: <Video size={18} /> },
    { name: "Сотрудники", path: "/dashboard/employees", icon: <Users size={18} /> },
    { name: "Журнал", path: "/dashboard/logs", icon: <AlertTriangle size={18} /> },
    { name: "Аналитика", path: "/dashboard/analytics", icon: <BarChart2 size={18} /> },
    { name: "Настройки", path: "/dashboard/settings", icon: <Settings size={18} /> },
  ];

  return (
    <aside
      className={`w-64 flex flex-col justify-between p-6 transition-all duration-300 ${
        darkMode ? "bg-[#111B1A] text-white" : "bg-white text-gray-900 border-r border-gray-200"
      }`}
    >
        {/* Логотип */}
      <div className="sidebar-logo mb-10 text-center">
        <Link to="/dashboard">
          <img
            src={darkMode ? LogoDark : LogoLight}
            alt="Логотип SafePass"
            className="w-40 mx-auto mb-3 transition-all duration-300"
          />
        </Link>
      </div>

      {/* Навигация */}
      <nav className="flex-1 space-y-2">
        {menuItems.map((item) => (
          <Link
            key={item.path}
            to={item.path}
            className={`flex items-center gap-3 px-4 py-2 rounded-lg transition ${
              location.pathname === item.path
                ? darkMode
                  ? "bg-[#16AF87]/25 text-[#16AF87]"
                  : "bg-[#16AF87]/10 text-[#16AF87]"
                : "hover:text-[#16AF87]"
            }`}
          >
            {item.icon}
            <span>{item.name}</span>
          </Link>
        ))}
      </nav>

    
    </aside>
  );
}

export default Sidebar;
