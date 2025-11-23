import React from "react";
import { useParams, useNavigate } from "react-router-dom";

const employeesList = [
  { id: 1, name: "Арайлым Жазыбаева", position: "Инженер безопасности", photo_url: "https://randomuser.me/api/portraits/women/65.jpg", badge_id: "B001", status: "active", last_seen_at: "2025-11-15 09:30", blocked: false },
  { id: 2, name: "Данияр Сулейменов", position: "Прораб", photo_url: "https://randomuser.me/api/portraits/men/12.jpg", badge_id: "B002", status: "inactive", last_seen_at: "2025-11-14 17:45", blocked: false },
  { id: 3, name: "Айшан Тлеуберген", position: "Электрик", photo_url: "https://randomuser.me/api/portraits/women/22.jpg", badge_id: "B003", status: "active", last_seen_at: "2025-11-15 08:15", blocked: false },
  { id: 4, name: "Ерболат Кайрат", position: "Сварщик", photo_url: "https://randomuser.me/api/portraits/men/33.jpg", badge_id: "B004", status: "active", last_seen_at: "2025-11-15 10:00", blocked: true },
  { id: 5, name: "Мадина Нуржан", position: "Кладовщик", photo_url: "https://randomuser.me/api/portraits/women/44.jpg", badge_id: "B005", status: "inactive", last_seen_at: "2025-11-14 16:30", blocked: false },
];

export default function EmployeeDetailPage({ darkMode }) {
  const { id } = useParams();
  const navigate = useNavigate();
  const employee = employeesList.find(emp => emp.id === parseInt(id));

  if (!employee) return <div className={`p-6 ${darkMode ? "bg-[#0F1514] text-white" : "bg-gray-50 text-gray-900"}`}>Сотрудник не найден</div>;

  const cardBg = darkMode ? "bg-[#1C2422]" : "bg-gray-100";
  const cardText = darkMode ? "text-white" : "text-gray-900";
  const tableBg = darkMode ? "bg-[#1C2422] text-white" : "bg-white text-gray-900";
  const tableHeaderBg = darkMode ? "bg-[#131918]" : "bg-gray-100";

  return (
    <div className={`p-6 min-h-screen ${darkMode ? "bg-[#0F1514] text-white" : "bg-gray-50 text-gray-900"}`}>
      <button
        onClick={() => navigate(-1)}
        className={`mb-4 px-3 py-1 rounded ${darkMode ? "bg-[#262E2B] text-white" : "bg-gray-300 text-gray-900"}`}
      >
        Назад
      </button>

      <h1 className="text-2xl font-semibold mb-4">{employee.name}</h1>

      <div className="flex flex-col md:flex-row gap-6 mb-6">
        <img src={employee.photo_url} alt={employee.name} className="w-32 h-32 rounded-full" />
        <div className="space-y-1">
          <p><strong>Должность:</strong> {employee.position}</p>
          <p><strong>ID бейджа:</strong> {employee.badge_id}</p>
          <p><strong>Статус:</strong> {employee.status === "active" ? "На объекте" : "Не на объекте"}</p>
          <p><strong>Последнее посещение:</strong> {employee.last_seen_at}</p>
          <p><strong>Заблокирован:</strong> <span className={employee.blocked ? "text-red-500 font-semibold" : ""}>{employee.blocked ? "Да" : "Нет"}</span></p>
        </div>
      </div>

      <div className={`grid grid-cols-1 md:grid-cols-3 gap-4 mb-6`}>
        <div className={`${cardBg} p-4 rounded shadow text-center`}>
          <p className="text-sm text-gray-400">Количество нарушений</p>
          <p className="text-xl font-bold">5</p>
        </div>
        <div className={`${cardBg} p-4 rounded shadow text-center`}>
          <p className="text-sm text-gray-400">Последнее нарушение</p>
          <p className="text-xl font-bold">2025-11-14</p>
        </div>
        <div className={`${cardBg} p-4 rounded shadow text-center`}>
          <p className="text-sm text-gray-400">Активность сегодня</p>
          <p className="text-xl font-bold">3 часа</p>
        </div>
      </div>

     <h2 className="text-xl font-semibold mb-4">История нарушений</h2>
<div className="overflow-x-auto">
  <table className={`min-w-full border-separate border-spacing-0 ${darkMode ? "bg-[#1C2422] text-white" : "bg-white text-gray-900"} rounded-lg shadow`}>
    <thead className={`${darkMode ? "bg-[#131918]" : "bg-gray-100"} rounded-t-lg`}>
      <tr>
        <th className="p-3 border-b border-gray-400 text-left">Дата</th>
        <th className="p-3 border-b border-gray-400 text-left">Тип нарушения</th>
        <th className="p-3 border-b border-gray-400 text-left">Описание</th>
      </tr>
    </thead>
    <tbody>
      <tr className={`${darkMode ? "bg-[#1C2422]" : "bg-white"} hover:${darkMode ? "bg-[#262E2B]" : "bg-gray-50"} transition`}>
        <td className="p-3 border-b border-gray-400">2025-11-14</td>
        <td className="p-3 border-b border-gray-400">Без каски</td>
        <td className="p-3 border-b border-gray-400">На объекте без каски</td>
      </tr>
      <tr className={`${darkMode ? "bg-[#1C2422]" : "bg-gray-50"} hover:${darkMode ? "bg-[#262E2B]" : "bg-gray-100"} transition`}>
        <td className="p-3 border-b border-gray-400">2025-11-10</td>
        <td className="p-3 border-b border-gray-400">Без допуска</td>
        <td className="p-3 border-b border-gray-400">Попытка входа на объект без допуска</td>
      </tr>
    </tbody>
  </table>
</div>

    </div>
  );
}
