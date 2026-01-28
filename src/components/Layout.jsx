import { Outlet } from "react-router-dom";
import { Sidebar } from "@/components/Sidebar";

export default function Layout() {
  return (
    <div className="flex h-screen bg-white">
      <Sidebar />

      <main className="flex-1 overflow-y-auto">
        <div className="border-b px-8 py-6 space-y-2">
          <h1 className="text-3xl font-semibold tracking-tight">
            Мониторинг посещаемости
          </h1>
          <p className="text-sm text-muted-foreground">
            Пропуски по неуважительной причине
          </p>
        </div>

        <div className="px-8 py-8">
          <Outlet />
        </div>

        <div className="border-t px-8 py-6 text-center text-sm text-muted-foreground">
          ГАПОУ ОКЭИ • {new Date().getFullYear()}
        </div>
      </main>
    </div>
  );
}
