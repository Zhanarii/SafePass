
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Navbar from "./components/Navbar";
import Sidebar from "./components/Sidebar";
import LoginPage from "./pages/LoginPage";
import DashboardPage from "./pages/DashboardPage";
import CameraPage from "./pages/CameraPage";
import LogsPage from "./pages/LogsPage";
import EmployeesPage from "./pages/EmployeesPage";
import EmployeeDetailPage from "./pages/EmployeeDetailPage";
import AnalyticsPage from "./pages/AnalyticsPage";
import SettingsPage from "./pages/SettingsPage";
import { useState } from "react";
import Header from "./components/Header";

function App() {

  const [darkMode, setDarkMode] = useState(true);
 const toggleDarkMode = () => setDarkMode(!darkMode);

  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />

        <Route
          path="/dashboard/*"
          element={
            <div className="grid grid-cols-[250px_1fr] h-screen bg-gray-50">
             
              <Sidebar darkMode={darkMode} toggleDarkMode={toggleDarkMode}/>

              
              <div className="flex flex-col overflow-hidden">
              
               
                <main className="flex-1 overflow-y-auto">
                  <Header darkMode={darkMode} toggleDarkMode={toggleDarkMode}/>
                  <Routes>
                    <Route path="/" element={<DashboardPage darkMode={darkMode}  toggleDarkMode={toggleDarkMode} />} />
                    <Route path="camera" element={<CameraPage darkMode={darkMode}  toggleDarkMode={toggleDarkMode}/>} />
                    <Route path="logs" element={<LogsPage darkMode={darkMode}  toggleDarkMode={toggleDarkMode} />} />
                    <Route path="employees" element={<EmployeesPage darkMode={darkMode}  toggleDarkMode={toggleDarkMode}/>} />
                    <Route path="employees/:id" element={<EmployeeDetailPage darkMode={darkMode}  toggleDarkMode={toggleDarkMode}/>} />
                    <Route path="analytics" element={<AnalyticsPage darkMode={darkMode}  toggleDarkMode={toggleDarkMode}/>} />
                    <Route path="settings" element={<SettingsPage darkMode={darkMode}  toggleDarkMode={toggleDarkMode}/>} />
                  </Routes>
                </main>
              </div>
            </div>
          }
        />

        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
