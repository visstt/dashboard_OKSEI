import { Routes, Route } from "react-router-dom";
import { AttendancePage } from "@/pages/AttendancePage";
import { GroupPage } from "@/pages/GroupPage";
import { LoginPage } from "@/pages/LoginPage";
import Layout from "@/components/Layout";

export function AppRouter() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />

      <Route element={<Layout />}>
        <Route index element={<AttendancePage />} />
        <Route path="group/:groupName" element={<GroupPage />} />
      </Route>
    </Routes>
  );
}
