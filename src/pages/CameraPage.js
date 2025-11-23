import React, { useState, useEffect } from "react";
import CountUp from 'react-countup';

// Mock данные для примера
const initialCameras = [
  { id: 1, name: "Камера №1", location: "Главный вход", enabled: true, fps_limit: 25, rtsp_url: "rtsp://...", src: null },
  { id: 2, name: "Камера №2", location: "Склад", enabled: false, fps_limit: 15, rtsp_url: "rtsp://...", src: null },
  { id: 3, name: "Камера №3", location: "Производственный цех", enabled: true, fps_limit: 20, rtsp_url: "rtsp://...", src: null },
  { id: 4, name: "Камера №4", location: "Территория", enabled: true, fps_limit: 30, rtsp_url: "rtsp://...", src: null },
];

// Компонент карточки камеры
function CameraCard({ cam, onSelect }) {
  return (
    <div
      onClick={() => onSelect(cam)}
      className={`cursor-pointer p-4 rounded-xl shadow-md transition-transform hover:scale-105 ${
        cam.enabled ? "bg-green-50" : "bg-red-50"
      }`}
    >
      <h3 className="font-semibold">{cam.name}</h3>
      <p className="text-sm text-gray-500">{cam.location}</p>
      <p className="text-xs">{cam.enabled ? "Активна" : "Выключена"}</p>
      <p className="text-xs">FPS лимит: {cam.fps_limit}</p>
    </div>
  );
}

export default function CameraPage({ darkMode, toggleDarkMode }) {
  const [cameras, setCameras] = useState(initialCameras);
  const [selectedCam, setSelectedCam] = useState(null);
  const [filterStatus, setFilterStatus] = useState("all"); // all | enabled | disabled
  const [formData, setFormData] = useState({ name: "", location: "", rtsp_url: "", enabled: true, fps_limit: 25 });
  const [currentTime, setCurrentTime] = useState(new Date());

  // Обновление времени каждую секунду
  useEffect(() => {
    const interval = setInterval(() => setCurrentTime(new Date()), 1000);
    return () => clearInterval(interval);
  }, []);

  const formatTime = (date) => date.toLocaleTimeString();

  const filteredCams = cameras.filter(cam => {
    if (filterStatus === "all") return true;
    if (filterStatus === "enabled") return cam.enabled;
    if (filterStatus === "disabled") return !cam.enabled;
  });

  const handleSave = () => {
    if (selectedCam) {
      setCameras(cameras.map(cam => (cam.id === selectedCam.id ? { ...cam, ...formData } : cam)));
    } else {
      const newCam = { ...formData, id: Date.now() };
      setCameras([...cameras, newCam]);
    }
    setFormData({ name: "", location: "", rtsp_url: "", enabled: true, fps_limit: 25 });
    setSelectedCam(null);
  };

  return (
    <div className={darkMode ? "bg-[#0F1514] text-white min-h-screen p-6" : "bg-gray-50 text-gray-900 min-h-screen p-6"}>
      <h2 className="text-2xl font-semibold mb-4">Управление камерами</h2>

    
{/* Панель — карточки в цветовой гамме сайта */}
<div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 mb-6">
  
  {/* Карточка 1 — Сотрудники */}
  <div
    onClick={() => {}}
    className={`
      cursor-pointer flex p-4 rounded-xl text-sm transition border
      ${darkMode
        ? "bg-[#1C2422] border-[#2A3331] text-gray-200 hover:bg-[#161D1C]"
        : "bg-white border-gray-300 text-gray-700 hover:bg-gray-100"
      }
    `}
  >
    {/* 1/3 — число */}
   <div
  className={`w-1/3 flex flex-col items-center justify-center text-5xl font-bold
  ${darkMode ?  "text-gray-300" : "text-gray-500"}`}
>
      <CountUp end={32} duration={1} />
    </div>

    {/* 2/3 — текст */}
    <div className="flex-1 min-w-0">
      <div className={`text-xs ${darkMode ? "text-gray-300" : "text-gray-500"}`}>Сотрудники на объекте</div>
      <div className={`text-sm font-medium ${darkMode ? "text-gray-100" : "text-gray-500"}`}>Активность: высокая</div>
      <div className={`text-xs mt-1 ${darkMode ? "text-gray-400" : "text-gray-500"}`}>Последнее: 2 мин назад</div>
    </div>
  </div>

  {/* Карточка 2 — Нарушения */}
  <div
    onClick={() => {}}
    className={`
      cursor-pointer flex p-4 rounded-xl text-sm transition border
      ${darkMode
        ? "bg-[#1C2422] border-[#2A3331] text-gray-200 hover:bg-[#161D1C]"
        : "bg-white border-gray-300 text-gray-700 hover:bg-gray-100"
      }
    `}
  >
    <div
  className={`w-1/3 flex flex-col items-center justify-center text-5xl font-bold
  ${darkMode ?  "text-gray-300" : "text-gray-500"}`}
>
      <CountUp end={5} duration={1} />
    </div>

    <div className="flex-1 min-w-0">
      <div className={`text-xs ${darkMode ? "text-gray-300" : "text-gray-500"}`}>Нарушения сегодня</div>
      <div className={`text-sm font-medium ${darkMode ? "text-gray-100" : "text-gray-500"}`}>Требует проверки: 2</div>
      <div className={`text-xs mt-1 ${darkMode ? "text-red-300" : "text-red-600"}`}>+1 за час</div>
    </div>
  </div>

  {/* Карточка 3 — Соблюдение правил */}
  <div
    onClick={() => {}}
    className={`
      cursor-pointer flex p-4 rounded-xl text-sm transition border
      ${darkMode
        ? "bg-[#1C2422] border-[#2A3331] text-gray-200 hover:bg-[#161D1C]"
        : "bg-white border-gray-300 text-gray-700 hover:bg-gray-100"
      }
    `}
  >
   <div
  className={`w-1/3 flex flex-col items-center justify-center text-5xl font-bold
  ${darkMode ?  "text-gray-300" : "text-gray-500"}`}
>
      <CountUp end={87} duration={1} suffix="%" />
    </div>

    <div className="flex-1 min-w-0">
      <div className={`text-xs ${darkMode ? "text-gray-300" : "text-gray-500"}`}>Соблюдение правил</div>
      <div className={`text-sm font-medium ${darkMode ? "text-gray-300" : "text-gray-500"}`}>Цель: 95%</div>
      <div className={`text-xs mt-1 ${darkMode ? "text-green-300" : "text-green-600"}`}>+2% за неделю</div>
    </div>
  </div>
</div>



     {/* 🔥 Онлайн-трансляции — версия 1:1 как на скриншоте */}
<section className={`p-6 rounded-2xl ${darkMode ? "bg-[#0F1414]" : "bg-white shadow"}`}>
  <h2 className="text-xl font-semibold mb-4">Онлайн-трансляции с камер</h2>

  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
    {cameras.slice(0, 4).map(cam => (
      <div
        key={cam.id}
        className={`flex flex-col rounded-xl overflow-hidden border ${
          darkMode ? "border-[#2A3331]" : "border-gray-300"
        }`}
      >
        {/* Верхняя полоса */}
        <div
          className="flex justify-between items-center px-3 py-1 text-[11px]"
          style={{
            background: darkMode ? "rgba(0,0,0,0.55)" : "#f3f3f3",
            color: darkMode ? "#fff" : "#000",
          }}
        >
          <div className="flex items-center gap-1">
            <span className="w-2 h-2 rounded-full bg-red-500"></span>
            <span className="w-2 h-2 rounded-full bg-red-500"></span>
            Онлайн-трансляция
          </div>

          <div>{formatTime(currentTime)}</div>
        </div>

        {/* Видео */}
        <div className="relative w-full h-72 bg-black flex items-center justify-center">
          {cam.src ? (
            <video src={cam.src} controls className="w-full h-full object-cover" />
          ) : (
            <div className="text-gray-400 text-xs">RTSP (placeholder)</div>
          )}
        </div>

        {/* Нижняя полоса */}
        <div className={`${darkMode ? "bg-[#131918]" : "bg-gray-100"} text-sm px-3 py-2`}>
          Текущая камера: {cam.name}
        </div>
      </div>
    ))}
  </div>
</section>

      {/* Кнопки управления выбранной камерой */}
{selectedCam && (
  <div className={`mt-6 p-4 rounded-xl shadow-md ${darkMode ? "bg-[#1C2422]" : "bg-white"}`}>
    <h3 className="font-semibold text-lg mb-4">Детали камеры: {selectedCam.name}</h3>
    <video
      src={selectedCam.rtsp_url} // проксирование через HLS/WebRTC
      controls
      className="w-full h-64 bg-black rounded mb-4"
    />
    <div className="flex flex-wrap gap-3">
      <button className={`flex-1 min-w-[120px] px-4 py-3 rounded-lg font-medium transition transform hover:-translate-y-1
        ${darkMode ? "bg-green-600 text-white hover:bg-green-700" : "bg-green-500 text-white hover:bg-green-600"}`}>
        Старт
      </button>
      <button className={`flex-1 min-w-[120px] px-4 py-3 rounded-lg font-medium transition transform hover:-translate-y-1
        ${darkMode ? "bg-red-600 text-white hover:bg-red-700" : "bg-red-500 text-white hover:bg-red-600"}`}>
        Стоп
      </button>
      <button className={`flex-1 min-w-[120px] px-4 py-3 rounded-lg font-medium transition transform hover:-translate-y-1
        ${darkMode ? "bg-yellow-500 text-black hover:bg-yellow-600" : "bg-yellow-400 text-black hover:bg-yellow-500"}`}>
        Снимок
      </button>
      <button className={`flex-1 min-w-[120px] px-4 py-3 rounded-lg font-medium transition transform hover:-translate-y-1
        ${darkMode ? "bg-blue-600 text-white hover:bg-blue-700" : "bg-blue-500 text-white hover:bg-blue-600"}`}>
        Тест потока
      </button>
      {/* Кнопка для экспорта / сохранения видео */}
      <button
        className={`flex-1 min-w-[120px] px-4 py-3 rounded-lg font-medium transition transform hover:-translate-y-1
          ${darkMode ? "bg-gray-700 text-white hover:bg-gray-600" : "bg-gray-200 text-gray-900 hover:bg-gray-300"}`}
        onClick={() => {
          // Пример сохранения видео (скачивание)
          const link = document.createElement("a");
          link.href = selectedCam.rtsp_url;
          link.download = `${selectedCam.name}.mp4`;
          link.click();
        }}
      >
        Скачать видео
      </button>
    </div>
  </div>
)}

{/* Кнопки фильтров камер */}
<div className="flex flex-wrap gap-2 mb-6">
  {filteredCams.map(cam => (
    <button
      key={cam.id}
      onClick={() => setSelectedCam(cam)}
      className={`cursor-pointer rounded-lg border flex items-center justify-center h-16 text-sm font-medium transition transform hover:-translate-y-1 px-6
        ${selectedCam?.id === cam.id
          ? darkMode
            ? "bg-green-600 text-white border-green-700"
            : "bg-green-500 text-white border-green-600"
          : darkMode
            ? "bg-[#1C2422] border-[#2A3331] text-gray-200 hover:bg-[#161D1C]"
            : "bg-white border-gray-300 text-gray-700 hover:bg-gray-100"
        }`}
    >
      {cam.name}
    </button>
  ))}

 
</div>


    </div>

    
  );
}
