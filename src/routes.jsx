import { Routes, Route, Navigate } from "react-router-dom";
import { AttendancePage } from "@/pages/AttendancePage";
import { GroupPage } from "@/pages/GroupPage";
import { LoginPage } from "@/pages/LoginPage";
import Layout from "@/components/Layout";
import { useAuthStore } from "./store/auth.store";
import { useEffect } from "react";
import { NotFound } from "./pages/NotFound";

function Logout() {
  const logout = useAuthStore((state) => state.logout);

  useEffect(() => {
    logout();
  }, [logout]);

  return <Navigate to="/login" replace />;
}

export function AppRouter() {
  const isAuth = useAuthStore((state) => state.isAuth);

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      {isAuth && (
        <>
          <Route path="/logout" element={<Logout />} />
          <Route element={<Layout />}>
            <Route index element={<AttendancePage />} />
            <Route path="group/:groupName" element={<GroupPage />} />
          </Route>
        </>
      )}
      <Route path="*" element={<NotFound isAuth={isAuth} />} />
    </Routes>
  );
}
