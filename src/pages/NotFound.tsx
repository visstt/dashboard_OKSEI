import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useNavigate } from "react-router-dom";
import { FileQuestion } from "lucide-react";

export function NotFound({ isAuth }: { isAuth: boolean }) {
  const navigate = useNavigate();

  return (
    <div className="flex min-h-screen items-center justify-center bg-neutral-50 text-neutral-900">
      <Card className="w-105 border-neutral-200 bg-white shadow-lg">
        <CardContent className="flex flex-col items-center gap-6 p-10 text-center">
          <FileQuestion className="h-16 w-16 text-neutral-400" />

          <div className="space-y-2">
            <h1 className="text-4xl font-bold tracking-tight">404</h1>
            <p className="text-sm text-neutral-500">
              Страница не найдена или была перемещена
            </p>
          </div>

          <div className="flex gap-3">
            <Button
              variant="outline"
              onClick={() => (isAuth ? navigate(-1) : navigate("/login"))}
            >
              Назад
            </Button>

            <Button
              onClick={() => (isAuth ? navigate("/") : navigate("/login"))}
            >
              На главную
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
