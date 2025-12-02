import React, { useState } from "react";

/**
 * SettingsPage
 * props:
 *  - darkMode (boolean)  — текущее состояние темы
 *  - setDarkMode (fn)    — функция переключения темы (из App)
 *
 * Все подписи на русском, все переключатели контролируемые (работают).
 */
export default function SettingsPage({ darkMode, setDarkMode }) {
  // локальные переключатели (пример)
  const [pushNotifications, setPushNotifications] = useState(true);
  const [emailNotifications, setEmailNotifications] = useState(false);
  const [appUpdates, setAppUpdates] = useState(true);
  const [largeText, setLargeText] = useState(false);

  // Для демонстрации — имя и email на русском
  const user = { name: "Павел Романов", email: "pavel@gmail.com" };

  return (
    <div className={`p-6 min-h-screen ${darkMode ? "bg-[#111312] text-white" : "bg-gray-100 text-black"}`}>
      {/* Заголовок */}
      <h1 className="text-2xl font-semibold mb-6">Настройки</h1>

      {/* Профиль */}
      <section className={`${darkMode ? "bg-[#1C2422] border-[#1C2422]" : "bg-white border-gray-200"} rounded-xl p-5 shadow mb-6 border`}>
        <div className="flex items-center gap-4">
          <img
            src="/avatar.png"
            alt="аватар"
            className="w-16 h-16 rounded-full border"
          />
          <div>
            <p className="font-semibold text-lg">{user.name}</p>
            <p className={`${darkMode ? "text-gray-300" : "text-gray-500"} text-sm`}>{user.email}</p>
          </div>
        </div>

        <div className="mt-4 flex gap-3">
          <button
            className={`px-3 py-2 rounded-lg border transition ${
              darkMode ? "bg-[#0F1514] border-[#1C2422] text-gray-200 hover:bg-[#1C2422] hover:text-white" : "bg-white border-gray-300 text-gray-700 hover:bg-gray-100"
            }`}
          >
            Редактировать профиль
          </button>

        
        </div>
      </section>

      {/* Уведомления */}
      <section className={`${darkMode ? "bg-[#1C2422]" : "bg-white"} rounded-xl p-5 shadow mb-6 border ${darkMode ? "border-[#1C2422]" : "border-gray-200"}`}>
        <h2 className="text-lg font-medium mb-3">Уведомления</h2>

        <SettingToggle
          label="Push-уведомления"
          checked={pushNotifications}
          onChange={() => setPushNotifications((v) => !v)}
          darkMode={darkMode}
        />

        <SettingToggle
          label="Email-уведомления"
          checked={emailNotifications}
          onChange={() => setEmailNotifications((v) => !v)}
          darkMode={darkMode}
        />

        <SettingToggle
          label="Обновления приложения"
          checked={appUpdates}
          onChange={() => setAppUpdates((v) => !v)}
          darkMode={darkMode}
        />
      </section>

      {/* Внешний вид */}
      <section className={`${darkMode ? "bg-[#1C2422]" : "bg-white"} rounded-xl p-5 shadow mb-6 border ${darkMode ? "border-[#1C2422]" : "border-gray-200"}`}>
        <h2 className="text-lg font-medium mb-3">Внешний вид</h2>

        {/* Переключатель темы — использует setDarkMode из props */}
        <div className="flex items-center justify-between py-2">
          <div>
            <div className="font-medium">Тёмная тема</div>
            <div className="text-sm text-gray-400">Включить/выключить тёмную тему интерфейса</div>
          </div>

          <Toggle
            checked={darkMode}
            onChange={() => {
              if (typeof setDarkMode === "function") setDarkMode(!darkMode);
            }}
            darkMode={darkMode}
            ariaLabel="Переключить тему"
          />
        </div>

        <SettingToggle
          label="Крупный шрифт"
          checked={largeText}
          onChange={() => setLargeText((v) => !v)}
          darkMode={darkMode}
        />
      </section>

      {/* Конфиденциальность */}
      <section className={`${darkMode ? "bg-[#1C2422]" : "bg-white"} rounded-xl p-5 shadow mb-6 border ${darkMode ? "border-[#1C2422]" : "border-gray-200"}`}>
        <h2 className="text-lg font-medium mb-3">Конфиденциальность</h2>

        <SettingLink label="Права доступа" darkMode={darkMode} />
        <SettingLink label="Скачать мои данные" darkMode={darkMode} />
        <SettingLink label="Удалить аккаунт" darkMode={darkMode} danger />
      </section>

      {/* Поддержка */}
      <section className={`${darkMode ? "bg-[#1C2422]" : "bg-white"} rounded-xl p-5 shadow mb-6 border ${darkMode ? "border-[#1C2422]" : "border-gray-200"}`}>
        <h2 className="text-lg font-medium mb-3">Поддержка</h2>
        <SettingLink label="FAQ" darkMode={darkMode} />
        <SettingLink label="Чат с поддержкой" darkMode={darkMode} />
        <SettingLink label="Отправить отзыв" darkMode={darkMode} />
      </section>

      {/* Выход */}
      <button
        className="w-full py-3 rounded-xl bg-red-500 text-white font-semibold mt-2 hover:opacity-95 transition"
        onClick={() => alert("Выход (заглушка)")}
      >
        Выйти из аккаунта
      </button>
    </div>
  );
}

/* ---------------------------------------------
   Toggle — универсальный контрол для ползунка
   принимает: checked (bool), onChange (fn), darkMode (bool)
   доступен по aria и клавиатуре
   --------------------------------------------- */
function Toggle({ checked, onChange, darkMode, ariaLabel = "toggle" }) {
  return (
    <button
      role="switch"
      aria-checked={checked}
      aria-label={ariaLabel}
      onClick={onChange}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          onChange && onChange();
        }
      }}
      className={`relative inline-flex items-center h-7 w-14 rounded-full transition-colors focus:outline-none ${
        checked
          ? darkMode
            ? "bg-[#1C2422]"
            : "bg-green-500"
          : darkMode
          ? "bg-[#0F1514] border border-[#1C2422]"
          : "bg-gray-300"
      }`}
    >
      <span
        className={`inline-block w-6 h-6 bg-white rounded-full shadow transform transition-transform ${
          checked ? "translate-x-7" : "translate-x-1"
        }`}
      />
    </button>
  );
}

/* ---------------------------------------------
   SettingToggle — строка с подписью и контролом
   controlled: принимает checked + onChange
   --------------------------------------------- */
function SettingToggle({ label, checked, onChange, darkMode }) {
  return (
    <div className="flex items-center justify-between py-3">
      <div>
        <div className="font-medium">{label}</div>
      </div>
      <Toggle checked={checked} onChange={onChange} darkMode={darkMode} ariaLabel={label} />
    </div>
  );
}

/* ---------------------------------------------
   SettingLink — кнопка-ссылка для настроек
   --------------------------------------------- */
function SettingLink({ label, danger = false, darkMode }) {
  return (
function SettingLink({ label, danger = false, darkMode }) {
  return (
    <button
      className={`w-full text-left py-3 text-sm rounded-md px-3 transition font-medium ${
        danger
          
          ? "text-[#16AF87] hover:bg-[#1C2422] hover:text-white"
          : "text-[#16AF87] hover:bg-gray-100"
      }`}
      onClick={() => alert(`${label} (заглушка)`)}
    >
      {label}
    </button>
  );
}

  );
}
