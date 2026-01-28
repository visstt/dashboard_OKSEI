# Система мониторинга посещаемости ГАПОУ ОКЭИ

Веб-интерфейс для контроля посещаемости студентов. Система представляет собой удобный дашборд, позволяющий отслеживать статистику посещаемости и получать информацию о пропусках.

## Технологии

### Frontend

- **React 19.2** - библиотека для построения пользовательских интерфейсов
- **Vite (Rolldown)** - современный инструмент сборки с быстрым HMR
- **React Router v7** - маршрутизация на стороне клиента
- **Tailwind CSS v4** - utility-first CSS фреймворк
- **shadcn/ui** - компоненты UI на основе Radix UI
- **Recharts** - библиотека для визуализации данных
- **Lucide React** - иконки
- **XLSX** - работа с Excel файлами

### Backend/Утилиты

- **Go** - конвертер данных посещаемости (в папке `converter/`)

### Инструменты разработки

- **Bun** - быстрый JavaScript runtime и пакетный менеджер
- **ESLint** - линтер для JavaScript/React
- **PostCSS** - обработка CSS

## Установка и запуск

### Требования

- [Node.js](https://nodejs.org/) v18+ или [Bun](https://bun.sh/) v1.0+
- [Go](https://golang.org/) 1.21+ (для работы с конвертером)

### Установка зависимостей

```bash
# Bun
bun install

# npm
npm install

# yarn
yarn install

# pnpm
pnpm install
```

### Запуск в режиме разработки

```bash
# Bun
bun run dev

# npm
npm run dev

# yarn
yarn dev

# pnpm
pnpm dev
```

Приложение будет доступно по адресу [http://localhost:5173](http://localhost:5173)

### Сборка для продакшена

```bash
# Bun
bun run build

# npm
npm run build

# yarn
yarn build

# pnpm
pnpm build
```

Собранные файлы будут находиться в папке `dist/`

### Предпросмотр продакшен-сборки

```bash
# Bun
bun run preview

# npm
npm run preview

# yarn
yarn preview

# pnpm
pnpm preview
```

### Линтинг

```bash
# Bun
bun run lint

# npm
npm run lint

# yarn
yarn lint

# pnpm
pnpm lint
```

## Структура проекта

```
├── src/
│   ├── components/    # React компоненты
│   │   ├── ui/       # UI компоненты (shadcn/ui)
│   │   └── ...
│   ├── pages/        # Страницы приложения
│   ├── lib/          # Утилиты и вспомогательные функции
│   └── utils/        # Конвертеры и обработчики данных
├── converter/        # Go-конвертер для данных посещаемости
├── public/           # Статические файлы
└── ...
```
