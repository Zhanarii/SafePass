import React, { useState } from "react";

// Фото сотрудников
const employeesPhotos = {
  "Арайлым Жазыбаева": "https://randomuser.me/api/portraits/women/65.jpg",
  "Данияр Сулейменов": "https://randomuser.me/api/portraits/men/12.jpg",
  "Айшан Тлеуберген": "https://randomuser.me/api/portraits/women/22.jpg",
  "Ерболат Кайрат": "https://randomuser.me/api/portraits/men/33.jpg",
  "Мадина Нуржан": "https://randomuser.me/api/portraits/women/44.jpg",
  "Иван Петров": "https://randomuser.me/api/portraits/men/15.jpg",
  "Алия Жумабаева": "https://randomuser.me/api/portraits/women/22.jpg",
  "Дмитрий Соколов": "https://randomuser.me/api/portraits/men/20.jpg",
  "Рустам Алиев": "https://randomuser.me/api/portraits/men/25.jpg",
};

const exportToCSV = (data, filename) => {
  if (!data || !data.length) return;
  const headers = Object.keys(data[0]);
  const csvRows = [
    headers.join(","),
    ...data.map(row =>
      headers.map(field => JSON.stringify(row[field], (_, v) => v ?? "")).join(",")
    ),
  ];
  const blob = new Blob([csvRows.join("\n")], { type: "text/csv" });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  window.URL.revokeObjectURL(url);
};

const buttonStyle = `cursor-pointer rounded-lg border flex items-center justify-center px-4 py-2 text-sm font-medium transition-transform duration-200`;

const getButtonClasses = (darkMode, variant = "primary", active = false) => {
  if (darkMode) {
    if (variant === "primary") return `${buttonStyle} ${active ? "bg-[#2A3331] text-white" : "bg-[#1C2422] text-gray-200 border-[#2A3331] hover:bg-[#24312E]"}`;
    if (variant === "secondary") return `${buttonStyle} bg-[#131918] text-gray-200 border-[#2A3331] hover:bg-[#1C2422]"`;
  } else {
    if (variant === "primary") return `${buttonStyle} ${active ? "bg-gray-300 text-gray-900" : "bg-white text-gray-700 border border-gray-300 hover:bg-gray-100"}`;
    if (variant === "secondary") return `${buttonStyle} bg-white text-gray-700 border border-gray-200 hover:bg-gray-100`;
  }
};

function LogsPage({ darkMode }) {
  const [selectedLog, setSelectedLog] = useState("violations");
  const [filters, setFilters] = useState({ employee: "", type: "", status: "" });

  const logs = {
    violations: [
      { id: 1, photo: employeesPhotos["Иван Петров"], name: "Иван Петров", time: "10:25", type: "Отсутствие каски", status: "new", camera: "Камера 1" },
      { id: 2, photo: employeesPhotos["Алия Жумабаева"], name: "Алия Жумабаева", time: "09:47", type: "Без жилета", status: "confirmed", camera: "Камера 2" },
    ],
    attendance: [
      { id: 1, photo: employeesPhotos["Арайлым Жазыбаева"], employee: "Арайлым Жазыбаева", time_in: "08:55", time_out: "17:30", badge_id: "B001", location: "Объект А" },
      { id: 2, photo: employeesPhotos["Данияр Сулейменов"], employee: "Данияр Сулейменов", time_in: "09:10", time_out: "18:00", badge_id: "B002", location: "Объект B" },
    ],
    access: [
      { id: 1, photo: employeesPhotos["Арайлым Жазыбаева"], employee: "Арайлым Жазыбаева", door: "Вход 1", time: "08:55", access_granted: true, badge_id: "B001" },
      { id: 2, photo: employeesPhotos["Данияр Сулейменов"], employee: "Данияр Сулейменов", door: "Вход 2", time: "09:10", access_granted: false, badge_id: "B002" },
    ],
  };

  const filteredLogs = {
    violations: logs.violations.filter(
      v =>
        (filters.employee === "" || v.name.toLowerCase().includes(filters.employee.toLowerCase())) &&
        (filters.type === "" || v.type.toLowerCase().includes(filters.type.toLowerCase())) &&
        (filters.status === "" || v.status === filters.status)
    ),
    attendance: logs.attendance.filter(
      v => filters.employee === "" || v.employee.toLowerCase().includes(filters.employee.toLowerCase())
    ),
    access: logs.access.filter(
      v => filters.employee === "" || v.employee.toLowerCase().includes(filters.employee.toLowerCase())
    ),
  };

  const currentData = filteredLogs[selectedLog] || [];
  const logNames = { violations: "Нарушения", attendance: "Посещаемость", access: "Доступ" };

  return (
    <div className={`${darkMode ? "bg-[#0F1514] text-white" : "bg-gray-100 text-gray-900"} min-h-screen p-8`}>
      <h2 className="text-2xl font-semibold mb-6">Журналы</h2>

      {/* Выбор журнала */}
      <div className="flex gap-4 mb-6">
        {["violations", "attendance", "access"].map(log => (
          <button
            key={log}
            onClick={() => setSelectedLog(log)}
            className={getButtonClasses(darkMode, "primary", selectedLog === log)}
          >
            {logNames[log]}
          </button>
        ))}
      </div>

      {/* Фильтры */}
      <div className="flex flex-wrap gap-4 mb-6">
        <input
          type="text"
          placeholder="Сотрудник"
          value={filters.employee}
          onChange={e => setFilters({ ...filters, employee: e.target.value })}
          className={`rounded-lg px-3 py-2 border focus:outline-none focus:ring-2 transition ${darkMode ? "bg-[#1C2422] text-gray-200 border-[#2A3331] focus:ring-gray-500" : "bg-white text-gray-900 border-gray-300 focus:ring-gray-400"}`}
        />
        {selectedLog === "violations" && (
          <>
            <input
              type="text"
              placeholder="Тип нарушения"
              value={filters.type}
              onChange={e => setFilters({ ...filters, type: e.target.value })}
              className={`rounded-lg px-3 py-2 border focus:outline-none focus:ring-2 transition ${darkMode ? "bg-[#1C2422] text-gray-200 border-[#2A3331] focus:ring-gray-500" : "bg-white text-gray-900 border-gray-300 focus:ring-gray-400"}`}
            />
            <select
              value={filters.status}
              onChange={e => setFilters({ ...filters, status: e.target.value })}
              className={`rounded-lg px-3 py-2 border focus:outline-none focus:ring-2 transition ${darkMode ? "bg-[#1C2422] text-gray-200 border-[#2A3331] focus:ring-gray-500" : "bg-white text-gray-900 border-gray-300 focus:ring-gray-400"}`}
            >
              <option value="">Все статусы</option>
              <option value="new">Новые</option>
              <option value="confirmed">Подтверждено</option>
              <option value="dismissed">Отклонено</option>
            </select>
          </>
        )}
      </div>

      {/* Таблица */}
      <div className={`rounded-2xl overflow-hidden border ${darkMode ? "bg-[#1C2422] border-[#2A3331]" : "bg-white border-gray-200 shadow-lg"}`}>
        <div className="flex justify-between items-center p-4 border-b">
          <h3 className="text-xl font-semibold">{logNames[selectedLog]}</h3>
          <button
            onClick={() => exportToCSV(currentData, `${logNames[selectedLog]}.csv`)}
            className={getButtonClasses(darkMode, "primary")}
          >
            Скачать CSV
          </button>
        </div>
        <table className="w-full text-left border-collapse">
          <thead className={`${darkMode ? "bg-[#232B29] text-gray-300" : "bg-gray-200 text-gray-700"}`}>
            <tr>
              {currentData.length > 0 ? (
                Object.keys(currentData[0]).map(key => <th key={key} className="p-4 font-medium capitalize">{key.replace(/_/g, " ")}</th>)
              ) : (
                <th className="p-4 text-gray-400">Нет данных</th>
              )}
            </tr>
          </thead>
          <tbody>
            {currentData.length > 0 ? (
              currentData.map(row => (
                <tr key={row.id} className={`border-t ${darkMode ? "border-[#2A3331] hover:bg-[#24312E]" : "border-gray-200 hover:bg-gray-50"} transition-colors`}>
                  {Object.keys(row).map(key => (
                    <td key={key} className="p-4">
                      {key === "photo" ? <img src={row[key]} alt={row.employee || row.name} className="w-24 h-16 object-cover rounded-lg" /> : String(row[key])}
                    </td>
                  ))}
                </tr>
              ))
            ) : (
              <tr>
                <td className="p-4 text-gray-400" colSpan={Object.keys(currentData[0] || {}).length || 1}>Нет данных</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default LogsPage;
