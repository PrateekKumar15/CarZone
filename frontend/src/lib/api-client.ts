const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface Car {
  id: number;
  brand: string;
  model: string;
  year: number;
  color: string;
  price_per_day: number;
  available: boolean;
  owner_id: number;
  created_at: string;
  updated_at: string;
  images?: string[];
  features?: string[];
}

export interface User {
  id: number;
  name: string;
  email: string;
}

export interface Booking {
  id: number;
  car_id: number;
  customer_id: number;
  owner_id: number;
  start_date: string;
  end_date: string;
  total_cost: number;
  status: string;
  created_at: string;
  updated_at: string;
}

class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor() {
    this.baseURL = API_BASE_URL;
    if (typeof window !== "undefined") {
      this.token = localStorage.getItem("auth_token");
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const config: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        ...(this.token && { Authorization: `Bearer ${this.token}` }),
      },
      ...options,
    };

    const response = await fetch(url, config);

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    return response.json();
  }

  // Auth methods
  async login(email: string, password: string) {
    const response = await this.request<{ token: string; user: User }>(
      "/auth/login",
      {
        method: "POST",
        body: JSON.stringify({ email, password }),
      }
    );

    this.token = response.token;
    if (typeof window !== "undefined") {
      localStorage.setItem("auth_token", response.token);
    }

    return response;
  }

  async register(name: string, email: string, password: string) {
    return this.request<{ token: string; user: User }>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ name, email, password }),
    });
  }

  async logout() {
    if (typeof window !== "undefined") {
      localStorage.removeItem("auth_token");
    }
    this.token = null;

    return this.request("/auth/logout", {
      method: "GET",
    });
  }

  // Car methods
  async getCars(): Promise<Car[]> {
    return this.request<Car[]>("/cars");
  }

  async getCarById(id: number): Promise<Car> {
    return this.request<Car>(`/cars/${id}`);
  }

  async getCarsByBrand(brand: string): Promise<Car[]> {
    return this.request<Car[]>(
      `/cars/brand?brand=${encodeURIComponent(brand)}`
    );
  }

  async createCar(
    car: Omit<Car, "id" | "created_at" | "updated_at">
  ): Promise<Car> {
    return this.request<Car>("/cars", {
      method: "POST",
      body: JSON.stringify(car),
    });
  }

  async updateCar(id: number, car: Partial<Car>): Promise<Car> {
    return this.request<Car>(`/cars/${id}`, {
      method: "PUT",
      body: JSON.stringify(car),
    });
  }

  async deleteCar(id: number): Promise<void> {
    return this.request<void>(`/cars/${id}`, {
      method: "DELETE",
    });
  }

  // Booking methods
  async getBookings(): Promise<Booking[]> {
    return this.request<Booking[]>("/bookings");
  }

  async getBookingById(id: number): Promise<Booking> {
    return this.request<Booking>(`/bookings/${id}`);
  }

  async createBooking(
    booking: Omit<Booking, "id" | "created_at" | "updated_at">
  ): Promise<Booking> {
    return this.request<Booking>("/bookings", {
      method: "POST",
      body: JSON.stringify(booking),
    });
  }

  async updateBookingStatus(id: number, status: string): Promise<Booking> {
    return this.request<Booking>(`/bookings/${id}/status`, {
      method: "PUT",
      body: JSON.stringify({ status }),
    });
  }

  async deleteBooking(id: number): Promise<void> {
    return this.request<void>(`/bookings/${id}`, {
      method: "DELETE",
    });
  }

  async getBookingsByCustomer(customerId: number): Promise<Booking[]> {
    return this.request<Booking[]>(`/bookings/customer/${customerId}`);
  }

  async getBookingsByCar(carId: number): Promise<Booking[]> {
    return this.request<Booking[]>(`/bookings/car/${carId}`);
  }

  async getBookingsByOwner(ownerId: number): Promise<Booking[]> {
    return this.request<Booking[]>(`/bookings/owner/${ownerId}`);
  }
}

export const apiClient = new ApiClient();
