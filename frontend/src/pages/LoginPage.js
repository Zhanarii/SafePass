import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

function LoginPage() {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();

    if (username === "admin" && password === "1234") {
      localStorage.setItem("auth", "true");
      navigate("/dashboard");
    } else {
      setError("Неверный логин или пароль");
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-br from-[#0f1f19] via-[#132b23] to-[#1a3c30]">
      <div className="bg-white/10 backdrop-blur-md rounded-2xl shadow-2xl p-10 w-full max-w-md border border-white/20">
        <h2 className="text-3xl font-bold text-white text-center mb-2">
          SafePass
        </h2>
        <p className="text-gray-300 text-center mb-8">
          Панель мониторинга средств индивидуальной защиты (PPE)
        </p>

        <form onSubmit={handleSubmit} className="flex flex-col gap-5">
          <div>
            <label className="block text-gray-200 mb-1">Логин</label>
            <input
              type="text"
              placeholder="Введите логин..."
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              className="w-full px-4 py-3 rounded-xl bg-white/10 border border-white/30 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[#16AF87] transition"
            />
          </div>

          <div>
            <label className="block text-gray-200 mb-1">Пароль</label>
            <input
              type="password"
              placeholder="Введите пароль..."
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full px-4 py-3 rounded-xl bg-white/10 border border-white/30 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[#16AF87] transition"
            />
          </div>

          {error && (
            <p className="text-red-400 text-sm text-center mt-1">{error}</p>
          )}

          <button
            type="submit"
            className="mt-4 bg-[#16AF87] hover:bg-[#149872] text-white font-semibold py-3 rounded-xl transition shadow-md shadow-[#16AF87]/25"
          >
            Войти
          </button>
        </form>

        <p className="text-gray-500 text-xs text-center mt-6">
          © 2025 SafePass — система мониторинга безопасности
        </p>
      </div>
    </div>
  );
}

export default LoginPage;
