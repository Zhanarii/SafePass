import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

const employeesList = [
  { id: 1, name: "Арайлым Жазыбаева", position: "Инженер безопасности", photo_url: "https://randomuser.me/api/portraits/women/65.jpg", status: "active", last_seen_at: "2025-11-15 09:30", blocked: false },
  { id: 2, name: "Данияр Сулейменов", position: "Прораб", photo_url: "https://randomuser.me/api/portraits/men/12.jpg", status: "inactive", last_seen_at: "2025-11-14 17:45", blocked: false },
  { id: 3, name: "Айшан Тлеуберген", position: "Электрик", photo_url: "https://randomuser.me/api/portraits/women/22.jpg", status: "active", last_seen_at: "2025-11-15 08:15", blocked: false },
  { id: 4, name: "Ерболат Кайрат", position: "Сварщик", photo_url: "https://randomuser.me/api/portraits/men/33.jpg", status: "active", last_seen_at: "2025-11-15 10:00", blocked: true },
  { id: 5, name: "Мадина Нуржан", position: "Кладовщик", photo_url: "https://randomuser.me/api/portraits/women/44.jpg", status: "inactive", last_seen_at: "2025-11-14 16:30", blocked: false },
];

export default function EmployeesPage({ darkMode }) {
  const navigate = useNavigate();
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("");

  const filteredEmployees = employeesList.filter(emp => {
    return (
      emp.name.toLowerCase().includes(search.toLowerCase()) &&
      (status === "" || emp.status === status)
    );
  });

  return (
    <div className={`p-6 min-h-screen ${darkMode ? "bg-[#0F1514] text-white" : "bg-gray-50 text-gray-900"}`}>
      <h1 className="text-3xl font-bold mb-6">Сотрудники</h1>

      {/* Фильтры */}
      <div className="flex flex-col md:flex-row gap-4 mb-6">
        <input
          type="text"
          placeholder="Поиск по имени"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className={`flex-1 border rounded-lg px-4 py-2 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400 ${darkMode ? "bg-[#1C2422] border-[#2A3331] text-white placeholder-gray-400" : "bg-white border-gray-300 text-gray-900 placeholder-gray-500"}`}
        />
        <select
          value={status}
          onChange={(e) => setStatus(e.target.value)}
          className={`border rounded-lg px-4 py-2 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400 ${darkMode ? "bg-[#1C2422] border-[#2A3331] text-white" : "bg-white border-gray-300 text-gray-900"}`}
        >
          <option value="">Все статусы</option>
          <option value="active">На объекте</option>
          <option value="inactive">Не на объекте</option>
        </select>
      </div>

      {/* Таблица */}
      <div className={`overflow-x-auto rounded-xl shadow-lg ${darkMode ? "bg-[#1C2422]" : "bg-white"}`}>
        <table className="min-w-full text-left border-collapse">
          <thead className={`${darkMode ? "bg-[#131918]" : "bg-gray-100"}`}>
            <tr>
              <th className="p-3">Фото</th>
              <th className="p-3">Имя</th>
              <th className="p-3">Должность</th>
              <th className="p-3">Статус</th>
              <th className="p-3">Действия</th>
            </tr>
          </thead>
          <tbody>
            {filteredEmployees.map(emp => (
              <tr key={emp.id} className={`transition hover:scale-[1.01] ${darkMode ? "hover:bg-[#262E2B]" : "hover:bg-gray-50"} border-b`}>
                <td className="p-3">
                  <img src={emp.photo_url} alt={emp.name} className="w-12 h-12 rounded-full object-cover border-2 border-gray-300" />
                </td>
                <td className="p-3 font-medium">{emp.name}</td>
                <td className="p-3 text-gray-500">{emp.position}</td>
<td className="p-3">
  {emp.status === "active" ? "На объекте" : "Не на объекте"}
</td>


                <td className="p-3">
                  <button
                    onClick={() => navigate(`/dashboard/employees/${emp.id}`)}
                    className={`cursor-pointer rounded-lg border flex items-center justify-center h-16 text-xs font-medium transition hover:-translate-y-1 px-6 ${
                      darkMode
                        ? "border-[#2A3331] bg-[#131918] text-gray-200 hover:bg-[#1C2422]"
                        : "border-gray-200 bg-white text-gray-700 hover:bg-gray-100"
                    }`}
                  >
                    Просмотр
                  </button>
                </td>
              </tr>
            ))}
            {filteredEmployees.length === 0 && (
              <tr>
                <td colSpan={5} className="p-4 text-center text-gray-500">
                  Сотрудники не найдены
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
