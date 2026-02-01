export interface AuthState {
  isAuth: boolean;
  login: () => void;
  logout: () => void;
}
