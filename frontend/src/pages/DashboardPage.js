import React, { useState, useEffect, useMemo } from "react";
import Sidebar from "../components/Sidebar";
import Header from "../components/Header";
import {
  LineChart,
  Line,
  ResponsiveContainer,
  XAxis,
  YAxis,
  Tooltip,
} from "recharts";

import CountUp from "react-countup";

function DashboardPage({ darkMode, toggleDarkMode }) {
  const [period, setPeriod] = useState(30); // 7 | 14 | 30

  // Данные (пример)
  const employeeData = [
    { day: "01.11", value: 20 }, { day: "02.11", value: 22 }, { day: "03.11", value: 25 },
    { day: "04.11", value: 24 }, { day: "05.11", value: 27 }, { day: "06.11", value: 25 }, { day: "07.11", value: 26 },
  ];
  const violationsData = [
    { day: "01.11", value: 5 }, { day: "02.11", value: 4 }, { day: "03.11", value: 3 },
    { day: "04.11", value: 6 }, { day: "05.11", value: 3 }, { day: "06.11", value: 4 }, { day: "07.11", value: 2 },
  ];
  const complianceData = [
    { day: "01.11", value: 85 }, { day: "02.11", value: 88 }, { day: "03.11", value: 91 },
    { day: "04.11", value: 90 }, { day: "05.11", value: 93 }, { day: "06.11", value: 94 }, { day: "07.11", value: 92 },
  ];
const [employeeCount, setEmployeeCount] = useState(0);

useEffect(() => {
  // Генерируем более реальные данные по периодам
  const base = 40; // минимальное значение
  const multiplier = {
    7: 1.2,
    14: 1.8,
    30: 2.6
  };

  const randomBoost = Math.floor(Math.random() * 10) + 5; // случайный +5…15
  setEmployeeCount(Math.floor(base * multiplier[period] + randomBoost));
}, [period]);

  const complianceLevel = 65;

  // стили карточки (был ReferenceError — теперь определены)
  const cardBase = `rounded-2xl p-6 transition-all duration-300 h-72 flex flex-col justify-between`;
  const lightBg = "bg-white shadow-lg";
  const darkBg = "bg-[#1C2422] shadow-none";

  // --- камеры ---
  const cameras = [
    { id: 1, label: "Камера №1 — Главный вход", title: "Главный вход", src: null },
    { id: 2, label: "Камера №2 — Производственный цех", title: "Производственный цех", src: null },
    { id: 3, label: "Камера №3 — Склад", title: "Склад", src: null },
    { id: 4, label: "Камера №4 — Территория", title: "Территория", src: null },
  ];
  const [mainCameraId, setMainCameraId] = useState(1);
  const [currentTime, setCurrentTime] = useState(new Date());

  useEffect(() => {
    const t = setInterval(() => setCurrentTime(new Date()), 1000);
    return () => clearInterval(t);
  }, []);

  const formatTime = (d) => {
    const hh = String(d.getHours()).padStart(2, "0");
    const mm = String(d.getMinutes()).padStart(2, "0");
    const ss = String(d.getSeconds()).padStart(2, "0");
    return `${hh}:${mm}:${ss}`;
  };

  const getCam = (id) => cameras.find((c) => c.id === id) || cameras[0];

  function getDateRange(days) {
    const end = new Date();
    const start = new Date();
    start.setDate(end.getDate() - days + 1);
    const options = { day: "2-digit", month: "2-digit" };
    return `${start.toLocaleDateString("ru-RU", options)} - ${end.toLocaleDateString("ru-RU", options)}`;
  }

  // --- подготовка данных для графика (chartData) ---
  const chartData = useMemo(() => {
    // простой подход: если исходных дней меньше, циклим
    const base = violationsData.length ? violationsData : [{ day: "—", value: 0 }];
    const res = [];
    for (let i = 0; i < period; i++) {
      res.push(base[i % base.length]);
    }
    return res.slice(-period);
  }, [violationsData, period]);

  const totalViolations = chartData.reduce((s, x) => s + (x.value || 0), 0);
  const potentialReduction = Math.round(totalViolations * 0.35); // пример 35%

  {/* Универсальный стиль кнопок */}
const buttonStyle = `cursor-pointer rounded-lg border flex items-center justify-center 
h-16 text-xs font-medium transition hover:-translate-y-1 px-6`;

const getButtonClasses = (darkMode, active = false) => {
  if (darkMode) {
    return `
      ${buttonStyle}
      border-[#2A3331]
      ${active ? "bg-[#1C2422]" : "bg-[#131918]"}
      text-gray-200
      hover:bg-[#1C2422]
    `;
  } else {
    return `
      ${buttonStyle}
      border-gray-200
      ${active ? "bg-gray-100" : "bg-white"}
      text-gray-700
      hover:bg-gray-100
    `;
  }
};

  // --- render ---
  return (
    <div className={darkMode ? "bg-[#0F1514] text-white min-h-screen" : "bg-gray-50 text-gray-900 min-h-screen"}>
      <div className="flex min-h-screen">
        <div className="flex-1 flex flex-col">
          <main className="flex-1 p-8 space-y-6">
            <div className="flex items-center justify-between">
              <h2 className="text-2xl font-semibold">SafePass Dashboard</h2>
              <p className="text-sm opacity-70">Панель управления</p>
            </div>

            {/* VIDEO BLOCK */}
            <section className={`p-6 rounded-2xl ${darkMode ? "bg-[#1C2422]" : "bg-white shadow"}`}>
              <div className="flex flex-col md:flex-row gap-4">
                {/* Main video */}
                <div className="flex-1 rounded-xl overflow-hidden border flex flex-col" style={{ borderColor: darkMode ? "#2A3331" : "#E5E7EB" }}>
                  <div className="relative w-full h-64 md:h-96 bg-black/70 flex items-center justify-center">
                    <div className="text-center">
                      <div className="text-gray-300/80 text-xs mb-2">RTSP (placeholder)</div>
                    </div>

                    <div className="absolute top-3 left-3 flex items-center gap-2">
                      <div className="relative flex items-center">
                        <span className="absolute inline-flex h-3 w-3 rounded-full bg-red-500 opacity-75 animate-ping"></span>
                        <span className="relative inline-flex h-3 w-3 rounded-full bg-red-600"></span>
                      </div>
                      <div className="text-[12px] font-semibold text-white/90">Онлайн-трансляция</div>
                    </div>

                    <div className="absolute top-3 right-3 px-3 py-1 text-xs rounded-lg font-medium"
                         style={{ background: darkMode ? "rgba(0,0,0,0.4)" : "rgba(255,255,255,0.85)", color: darkMode ? "#fff" : "#111" }}>
                      {formatTime(currentTime)}
                    </div>
                  </div>

                  <div className={`p-3 ${darkMode ? "bg-[#131918]" : "bg-gray-50"} text-sm`}>
                    Текущая камера: <span className="font-semibold">{getCam(mainCameraId).label}</span>
                  </div>
                </div>

                {/* minis */}
                <div className="flex flex-col gap-4 w-32 md:w-40">
                  {cameras.filter(c => c.id !== mainCameraId).slice(0, 3).map(cam => (
                    <div key={cam.id}
                         onClick={() => setMainCameraId(cam.id)}
                         className={`cursor-pointer rounded-lg overflow-hidden border transition-transform transform hover:-translate-y-1 ${darkMode ? "border-[#2A3331]" : "border-gray-200"}`}>
                      <div className="relative w-full h-20 bg-black/60 flex items-center justify-center">
                        <div className="text-center">
                          <div className="mb-1 text-xs opacity-80">{cam.label}</div>
                        </div>
                        <div className="absolute top-1 right-1 px-1 py-0.5 text-[10px] rounded-md font-medium"
                             style={{ background: darkMode ? "rgba(0,0,0,0.45)" : "rgba(255,255,255,0.85)", color: darkMode ? "#fff" : "#111" }}>
                          {formatTime(currentTime)}
                        </div>
                      </div>
                      <div className={`px-2 py-1 text-xs ${darkMode ? "bg-[#131918] text-gray-100" : "bg-white text-gray-800"}`}>
                        <button
                          onClick={e => { e.stopPropagation(); setMainCameraId(cam.id); }}
                          className="w-full text-xs px-2 py-1 rounded bg-[#16AF87]/10 text-[#16AF87] hover:bg-[#16AF87]/20"
                        >
                          Сделать главной
                        </button>
                      </div>
                    </div>
                  ))}

                  <button
                    className={`cursor-pointer rounded-lg border flex items-center justify-center h-20 text-xs font-medium transition hover:-translate-y-1 ${
                      darkMode ? "border-[#2A3331] bg-[#131918] text-gray-200 hover:bg-[#1C2422]" : "border-gray-200 bg-white text-gray-700 hover:bg-gray-100"
                    }`}
                  >
                    Показать больше
                  </button>
                </div>
              </div>
            </section>

            {/* === STATS GRID: обязательно обернуть в section grid === */}
            <section className="grid grid-cols-3 gap-6">

 

{/* блок-1 */}
<div className={`${cardBase} ${darkMode ? darkBg : lightBg}`}>
  <div className="flex justify-between items-start">
    <div>
      <p className={`text-sm ${darkMode ? "text-gray-300" : "text-gray-500"}`}>Сотрудники на объекте</p>
    </div>
    <div>
      <select
        value={period}
        onChange={(e) => setPeriod(Number(e.target.value))}
        className={`text-xs rounded-lg px-3 py-1 font-medium transition-all duration-200 ${darkMode ? "bg-[#0C1311] text-white" : "bg-gray-100 text-gray-800"}`}
      >
        <option value={7}>7 дней</option>
        <option value={14}>14 дней</option>
        <option value={30}>30 дней</option>
      </select>
    </div>
  </div>

  <div className="flex-1 flex items-center justify-center">
    <div className="text-center">
     <div className="text-8xl font-extrabold leading-none" style={{ color: "#16AF87" }}>
  <CountUp end={employeeCount} duration={1.5} separator=" " />
</div>

    </div>
  </div>

  <div className={`mt-2 text-sm ${darkMode ? "text-gray-300" : "text-gray-500"}`}>
    {new Date(Date.now() - period * 24 * 60 * 60 * 1000).toLocaleDateString("ru-RU", {
      day: "2-digit",
      month: "2-digit"
    })} — {new Date().toLocaleDateString("ru-RU", {
      day: "2-digit",
      month: "2-digit"
    })}
  </div>
</div>


{/* блок-2 */}
<div className={`${cardBase} ${darkMode ? darkBg : lightBg} p-6 flex flex-col justify-between`}>
  {/* Заголовок и селект */}
  <div className="flex justify-between items-center mb-4">
    <p className={`text-sm font-medium ${darkMode ? "text-gray-300" : "text-gray-500"}`}>Нарушения</p>
    <select
      value={period}
      onChange={(e) => setPeriod(Number(e.target.value))}
      className={`text-xs rounded-lg px-3 py-1 font-medium transition-all duration-200 ${darkMode ? "bg-[#0C1311] text-white" : "bg-gray-100 text-gray-800"}`}
    >
      <option value={7}>7 дней</option>
      <option value={14}>14 дней</option>
      <option value={30}>30 дней</option>
    </select>
  </div>

  {/* График */}
  <div className="flex-1 w-full h-36 mb-4">
    <ResponsiveContainer width="100%" height="100%">
      <LineChart
        data={chartData.map((d, i) => ({ ...d, day: i + 1 }))}
        margin={{ top: 8, right: 8, left: -10, bottom: 6 }}
      >
        <XAxis dataKey="day" stroke={darkMode ? "#94A3B8" : "#374151"} tick={{ fontSize: 12 }} />
        <YAxis stroke={darkMode ? "#94A3B8" : "#374151"} allowDecimals={false} tick={{ fontSize: 12 }} width={36} />
        <Tooltip
          contentStyle={{ backgroundColor: darkMode ? "#0F1514" : "#fff", borderRadius: 6, border: "none" }}
          itemStyle={{ color: "#EF4444" }}
          formatter={(value) => [`${value} наруш.`]}
          labelFormatter={(label) => `День ${label}`}
        />
        <Line type="monotone" dataKey="value" stroke="#EF4444" strokeWidth={2} dot={{ r: 3 }} />
      </LineChart>
    </ResponsiveContainer>
  </div>

  {/* Статистика */}
  <div className="flex justify-between items-center mb-4">
    <div>
      <div className={`text-sm ${darkMode ? "text-gray-400" : "text-gray-500"}`}>Нарушения</div>
      <div className="text-2xl font-bold">{totalViolations}</div>
    </div>
    <div className="text-right">
      <div className={`text-sm ${darkMode ? "text-gray-400" : "text-gray-500"}`}>Сокращение</div>
      <div className="text-2xl font-semibold text-[#F59E0B]">≈ {potentialReduction}</div>
    </div>
  </div>

  {/* Период */}
   <div className={`mt-2 text-sm ${darkMode ? "text-gray-300" : "text-gray-500"}`}>
    {new Date(Date.now() - period * 24 * 60 * 60 * 1000).toLocaleDateString("ru-RU", {
      day: "2-digit",
      month: "2-digit"
    })} — {new Date().toLocaleDateString("ru-RU", {
      day: "2-digit",
      month: "2-digit"
    })}
  </div>
</div>

{/* блок-3 */}
  <div className={`${cardBase} ${darkMode ? darkBg : lightBg} p-6`}>
  <div className="flex justify-between items-center mb-4">
    <p className={`text-sm font-medium ${darkMode ? "text-gray-300" : "text-gray-500"}`}>
      Уровень соблюдения ПБ
    </p>
    <select
      value={period}
      onChange={(e) => setPeriod(Number(e.target.value))}
      className={`text-xs rounded-lg px-3 py-1 font-medium transition-all duration-200 ${darkMode ? "bg-[#0C1311] text-white" : "bg-gray-100 text-gray-800"}`}
    >
      <option value={7}>7 дней</option>
      <option value={14}>14 дней</option>
      <option value={30}>30 дней</option>
    </select>
  </div>

  <div className="flex-1 flex flex-col items-center justify-center">
  <div style={{ width: 220, height: 220 }} className="relative">
    <svg viewBox="0 0 140 140" width="220" height="220">
      {/* Нижний круг */}
      <circle
        cx="70"
        cy="70"
        r="60"
        stroke={darkMode ? "#0F1514" : "#F3F4F6"}
        strokeWidth="10"
        fill="none"
        strokeLinecap="round"
      />

      {/* Прогресс-бар */}
      {(() => {
        const radius = 60;
        const circumference = 2 * Math.PI * radius;
        const offset = circumference * (1 - Math.min(100, Math.max(0, complianceLevel)) / 100);

        return (
          <circle
            cx="70"
            cy="70"
            r={radius}
            stroke="#16AF87"  // один цвет прогресса
            strokeWidth="10"
            fill="none"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={offset}
            transform="rotate(-90 70 70)"
            style={{ transition: "stroke-dashoffset 700ms ease" }}
          />
        );
      })()}
    </svg>


      <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
        <div className="text-3xl font-bold" style={{ color: darkMode ? "#E6FBF6" : "#063741" }}>
          {complianceLevel}%
        </div>
        <div className="text-sm mt-1 text-gray-400 text-center">
          Соблюдение правил безопасности
        </div>
      </div>
    </div>
  </div>
</div>
</section>













            {/* TABLE */}
          
<section className={`p-6 rounded-2xl ${darkMode ? "bg-[#1C2422]" : "bg-white shadow"}`}>
  <div className="flex justify-between items-center mb-4">
    <h3 className="text-lg font-semibold">Сотрудники с нарушениями</h3>
    <button
      className={`cursor-pointer rounded-lg border flex items-center justify-center h-16 text-xs font-medium transition hover:-translate-y-1 px-6 ${
        darkMode
          ? "border-[#2A3331] bg-[#131918] text-gray-200 hover:bg-[#1C2422]"
          : "border-gray-200 bg-white text-gray-700 hover:bg-gray-100"
      }`}
      onClick={() => window.location.href = "/all-violations"}
    >
      Показать больше
    </button>
  </div>
  <div className={`overflow-x-auto rounded-xl border ${darkMode ? "border-[#2A3331]" : "border-gray-200"}`}>
    <table className="min-w-full text-left border-collapse">
      <thead className={darkMode ? "bg-[#131918] text-gray-300" : "bg-gray-100 text-gray-700"}>
        <tr>
          <th className="p-3 font-medium">Сотрудник</th>
          <th className="p-3 font-medium">Должность</th>
          <th className="p-3 font-medium">Время прихода</th>
          <th className="p-3 font-medium">Нарушение</th>
        </tr>
      </thead>
      <tbody>
        {[
          { id: 1, name: "Иван Петров", position: "Инженер", time: "08:15", violation: "Без каски" },
          { id: 2, name: "Дмитрий Соколов", position: "Охранник", time: "09:05", violation: "Без жилета" },
          { id: 3, name: "Рустам Алиев", position: "Монтажник", time: "09:10", violation: "Нарушение зоны доступа" },
        ].map(emp => (
          <tr
            key={emp.id}
            className={`border-t transition-colors ${
              darkMode
                ? "border-[#2A3331] hover:bg-[#24312E]"
                : "border-gray-200 hover:bg-gray-50"
            }`}
          >
            <td className="p-3 font-medium">{emp.name}</td>
            <td className="p-3 text-gray-400">{emp.position}</td>
            <td className="p-3 text-gray-400">{emp.time}</td>
            <td className="p-3 font-semibold text-[#EF4444]">{emp.violation}</td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
</section>


          </main>
        </div>
      </div>
    </div>
  );
}

export default DashboardPage;
