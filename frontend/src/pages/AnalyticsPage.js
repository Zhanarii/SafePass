import React, { useState, useMemo } from "react";
import {
  PieChart,
  Pie,
  Cell,
  Tooltip as ReTooltip,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  ResponsiveContainer,
  BarChart,
  Bar,
} from "recharts";

const PIE_COLORS = ["#16AF87", "#F59E0B", "#EF4444", "#6B7280"];

// ---------------- Heatmap ----------------
function Heatmap({ matrix, xLabels, yLabels, darkMode }) {
  const maxVal = Math.max(...matrix.flat(), 1);
  const getColor = (val) => {
    const t = val / maxVal;
    const r = Math.round(240 * t + 40 * (1 - t));
    const g = Math.round(120 * (1 - t) + 200 * t);
    const b = Math.round(80 * (1 - t));
    return `rgb(${r},${g},${b})`;
  };

  return (
    <div className="w-full overflow-auto">
      <div className="flex items-center gap-3 mb-2">
        <div className="text-sm font-medium text-gray-500">Heatmap — Камеры × Время</div>
        <div className="text-xs text-gray-400">цвет — интенсивность событий</div>
      </div>

      <div className={`${darkMode ? "bg-[#1C2422] border-[#2A3331]" : "bg-white border-gray-100"} rounded-xl border p-4`}>
        <div className="overflow-auto">
          <table className="border-collapse table-auto w-full">
            <thead>
              <tr>
                <th className="p-2 w-24 text-left text-xs text-gray-400">Time\Cam</th>
                {xLabels.map((x, i) => (
                  <th key={i} className="p-2 text-xs text-center text-gray-400 min-w-[64px]">{x}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {yLabels.map((y, row) => (
                <tr key={row}>
                  <td className="p-2 text-xs text-gray-400">{y}</td>
                  {xLabels.map((x, col) => {
                    const val = matrix[row][col];
                    return (
                      <td key={col} className="p-1">
                        <div
                          className="w-12 h-6 rounded-md flex items-center justify-center text-[11px] font-medium text-white"
                          style={{ backgroundColor: getColor(val) }}
                          title={`${val} events`}
                        >
                          {val}
                        </div>
                      </td>
                    );
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

// ---------------- AnalyticsPanel ----------------
export default function AnalyticsPanel({ darkMode }) {
  const [period, setPeriod] = useState(14);

  // ---------------- данные ----------------
  const pieData = [
    { name: "Без каски", value: 40 },
    { name: "Без жилета", value: 30 },
    { name: "Без допуска", value: 20 },
    { name: "Прочее", value: 10 },
  ];

  const lineData = Array.from({ length: 24 }).map((_, hour) => ({
    hour: String(hour).padStart(2, "0"),
    violations: Math.floor(2 + 4 * Math.abs(Math.sin((hour / 24) * Math.PI * 2)) + Math.random() * 2),
  }));

  const barData = Array.from({ length: 24 }).map((_, hour) => ({
    hour: String(hour).padStart(2, "0"),
    staff: Math.max(0, Math.floor(10 + 15 * Math.abs(Math.cos((hour / 24) * Math.PI * 2)) + Math.random() * 4)),
  }));

  const cameras = ["Cam-1", "Cam-2", "Cam-3", "Cam-4", "Cam-5", "Cam-6"];
  const hours = Array.from({ length: 8 }).map((_, i) => {
    const h = new Date();
    h.setHours(h.getHours() - (7 - i));
    return h.getHours().toString().padStart(2, "0") + ":00";
  });
  const heatMatrix = hours.map(() => cameras.map(() => Math.floor(Math.random() * 8)));

  const cardBase = "rounded-2xl p-6 transition-all duration-300 flex flex-col justify-between";
  const lightBg = "bg-white shadow-lg";
  const darkBg = "bg-[#1C2422] shadow-none";

  const pieInner = useMemo(() => ({ nameKey: "name", dataKey: "value" }), []);

  // ---------------- export заглушка ----------------
  const handleExport = () => {
    alert("Экспорт в таблицу (заглушка)");
  };
  const getDashboardButtonClasses = (darkMode) =>
  `cursor-pointer rounded-lg border flex items-center justify-center h-16 text-xs font-medium transition hover:-translate-y-1 px-6 ${
    darkMode
      ? "border-[#2A3331] bg-[#131918] text-gray-200 hover:bg-[#1C2422]"
      : "border-gray-200 bg-white text-gray-700 hover:bg-gray-100"
  }`;


  return (
    <div className={`${darkMode ? "bg-[#0F1514] text-white" : "bg-gray-50 text-gray-900"} p-6 min-h-screen space-y-6`}>
      {/* Заголовок + период */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
        <h2 className="text-2xl font-semibold">Аналитика — безопасность</h2>
        <div className="flex items-center gap-3">
          <select
            value={period}
            onChange={(e) => setPeriod(Number(e.target.value))}
            className={`rounded-md px-3 py-2 text-sm border transition ${
              darkMode
                ? "bg-[#0C1311] text-white border-[#2A3331]"
                : "bg-white text-gray-800 border-gray-200"
            }`}
          >
            <option value={7}>7 дней</option>
            <option value={14}>14 дней</option>
            <option value={30}>30 дней</option>
          </select>

<button
  className={getDashboardButtonClasses(darkMode)}
  onClick={handleExport}
>
  Экспортировать
</button>


        </div>
      </div>

      {/* Карточки аналитики */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Pie */}
        <div className={`${cardBase} ${darkMode ? darkBg : lightBg}`}>
          <h3 className="text-lg font-medium mb-2">Нарушения по типам</h3>
          <ResponsiveContainer width="100%" height={200}>
            <PieChart>
              <Pie data={pieData} cx="50%" cy="50%" innerRadius={50} outerRadius={80} paddingAngle={4} {...pieInner}>
                {pieData.map((entry, idx) => (
                  <Cell key={`cell-${idx}`} fill={PIE_COLORS[idx % PIE_COLORS.length]} />
                ))}
              </Pie>
              <ReTooltip />
            </PieChart>
          </ResponsiveContainer>
          <div className="mt-3 flex gap-3 flex-wrap">
            {pieData.map((d, i) => (
              <div key={i} className="flex items-center gap-2 text-sm">
                <span className="w-3 h-3 rounded-sm" style={{ background: PIE_COLORS[i % PIE_COLORS.length] }} />
                <span className="text-gray-400">{d.name}</span>
                <span className="text-gray-200 font-medium ml-1">{d.value}%</span>
              </div>
            ))}
          </div>
        </div>

        {/* Line */}
        <div className={`${cardBase} ${darkMode ? darkBg : lightBg}`}>
          <h3 className="text-lg font-medium mb-2">Нарушения по времени</h3>
          <ResponsiveContainer width="100%" height={200}>
            <LineChart data={lineData}>
              <CartesianGrid strokeDasharray="3 3" stroke={darkMode ? "#232B29" : "#f3f4f6"} />
              <XAxis dataKey="hour" stroke={darkMode ? "#94A3B8" : "#374151"} tick={{ fontSize: 12 }} />
              <YAxis stroke={darkMode ? "#94A3B8" : "#374151"} />
              <ReTooltip />
              <Line type="monotone" dataKey="violations" stroke="#EF4444" strokeWidth={2} dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>

        {/* Bar */}
        <div className={`${cardBase} ${darkMode ? darkBg : lightBg}`}>
          <h3 className="text-lg font-medium mb-2">Активность сотрудников</h3>
          <ResponsiveContainer width="100%" height={200}>
            <BarChart data={barData}>
              <CartesianGrid strokeDasharray="3 3" stroke={darkMode ? "#232B29" : "#f3f4f6"} />
              <XAxis dataKey="hour" stroke={darkMode ? "#94A3B8" : "#374151"} tick={{ fontSize: 12 }} />
              <YAxis stroke={darkMode ? "#94A3B8" : "#374151"} />
              <ReTooltip />
              <Bar dataKey="staff" fill="#16AF87" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Heatmap */}
        <div className={`${cardBase} ${darkMode ? darkBg : lightBg}`}>
          <Heatmap matrix={heatMatrix} xLabels={cameras} yLabels={hours} darkMode={darkMode} />
        </div>
      </div>
    </div>
  );
}
