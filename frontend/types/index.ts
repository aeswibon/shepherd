export interface User {
  username: string;
  email: string;
  password: string;
}

export interface App {
  id: string;
  name: string;
  namespace: string;
  deployedAt: number;
  status: string;
}
