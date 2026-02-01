import { Routes, Route, Navigate } from "react-router-dom";
import { AttendancePage } from "@/pages/home/AttendancePage";
import { GroupPage } from "@/pages/groups/GroupPage";
import { LoginPage } from "@/pages/auth/LoginPage";
import Layout from "@/app/Layout";
import { useAuthStore } from "@/features/store";
import { useEffect } from "react";
import { NotFound } from "@/app/NotFound";
import AllGroupsPage from "@/pages/groups/AllGroupsPage";
import { HistoryPage } from "@/pages/history/HistoryPage";

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
            <Route path="/history" element={<HistoryPage />} />
            <Route path="/groups" element={<AllGroupsPage />} />
            <Route path="/group/:groupName" element={<GroupPage />} />
          </Route>
        </>
      )}
      <Route path="*" element={<NotFound isAuth={isAuth} />} />
    </Routes>
  );
}
