export interface User {
  id: string;
  email: string;
  nickname: string;
  staff: boolean;
  created_at: string;
}

export interface Schedule {
  id: string;
  title: string;
  capacity: number;
  created_at: string;
  reserved: number;
  reservations: Reservation[];
}

export interface Reservation {
  id: string;
  user: User
}
