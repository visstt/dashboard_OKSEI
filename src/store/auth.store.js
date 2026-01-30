import { create } from "zustand";
import { persist } from "zustand/middleware";

export const useAuthStore = create(
  persist(
    (set) => ({
      isAuth: false,

      // Симуляция входа
      login: () =>
        set({
          isAuth: true,
        }),

      // Симуляция выхода
      logout: () =>
        set({
          isAuth: false,
        }),
    }),
    {
      name: "auth-store",
    }
  )
);
