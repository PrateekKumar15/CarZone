const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface User {
  id: string;
  username: string;
  email: string;
  phone: string;
  role: string;
}

export interface Car {
  id: string;
  owner_id?: string;
  name: string;
  brand: string;
  model: string;
  year: number;
  color?: string;
  vin?: string;
  fuel_type: string;
  engine: Engine;
  location_city: string;
  location_state: string;
  location_country: string;
  price: {
    rental_price_daily: number;
    rental_price_weekly: number;
    rental_price_monthly: number;
    sale_price?: number;
  };
  status: string;
  availability_type: string;
  is_available: boolean;
  features: Record<string, unknown>;
  description: string;
  images: string[];
  mileage: number;
  created_at: string;
  updated_at: string;
  // Populated fields
  owner?: User;
}

export interface Engine {
  engine_size: number;
  cylinders: number;
  horsepower: number;
  torque: number;
  transmission: string;
}

export interface Booking {
  id: string;
  customer_id: string;
  car_id: string;
  owner_id: string;
  start_date: string;
  end_date: string;
  total_cost: number;
  status: string;
  created_at: string;
  updated_at: string;
  // Populated fields
  car?: Car;
  customer?: User;
  owner?: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  phone: string;
  role?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
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
      const errorText = await response.text();
      throw new Error(
        `API Error: ${response.status} ${response.statusText} - ${errorText}`
      );
    }

    return response.json();
  }

  // Authentication methods
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify(data),
    });

    this.token = response.token;
    if (typeof window !== "undefined") {
      localStorage.setItem("auth_token", response.token);
    }

    return response;
  }

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify({
        ...data,
        role: data.role || "customer",
      }),
    });

    this.token = response.token;
    if (typeof window !== "undefined") {
      localStorage.setItem("auth_token", response.token);
    }

    return response;
  }

  async logout(): Promise<void> {
    if (typeof window !== "undefined") {
      localStorage.removeItem("auth_token");
    }
    this.token = null;

    await this.request("/auth/logout", {
      method: "GET",
    });
  }

  // Car methods
  async getCars(): Promise<Car[]> {
    return this.request<Car[]>("/cars");
  }

  async getCarById(id: string): Promise<Car> {
    return this.request<Car>(`/cars/${id}`);
  }

  async createCar(
    carData: Omit<Car, "id" | "created_at" | "updated_at">
  ): Promise<Car> {
    return this.request<Car>("/cars", {
      method: "POST",
      body: JSON.stringify(carData),
    });
  }

  async updateCar(id: string, carData: Partial<Car>): Promise<Car> {
    return this.request<Car>(`/cars/${id}`, {
      method: "PUT",
      body: JSON.stringify(carData),
    });
  }

  async deleteCar(id: string): Promise<void> {
    await this.request(`/cars/${id}`, {
      method: "DELETE",
    });
  }

  // Booking methods
  async getBookings(): Promise<Booking[]> {
    return this.request<Booking[]>("/bookings");
  }

  async getBookingById(id: string): Promise<Booking> {
    return this.request<Booking>(`/bookings/${id}`);
  }

  async getBookingsByCustomer(customerId: string): Promise<Booking[]> {
    return this.request<Booking[]>(`/bookings/customer/${customerId}`);
  }

  async getBookingsByCar(carId: string): Promise<Booking[]> {
    return this.request<Booking[]>(`/bookings/car/${carId}`);
  }

  async getBookingsByOwner(ownerId: string): Promise<Booking[]> {
    return this.request<Booking[]>(`/bookings/owner/${ownerId}`);
  }

  async createBooking(bookingData: {
    customer_id: string;
    car_id: string;
    owner_id: string;
    booking_type: "rental" | "purchase";
    start_date?: string;
    end_date?: string;
    notes?: string;
  }): Promise<Booking> {
    return this.request<Booking>("/bookings", {
      method: "POST",
      body: JSON.stringify(bookingData),
    });
  }

  async updateBookingStatus(id: string, status: string): Promise<Booking> {
    return this.request<Booking>(`/bookings/${id}/status`, {
      method: "PUT",
      body: JSON.stringify({ status }),
    });
  }

  async deleteBooking(id: string): Promise<void> {
    await this.request(`/bookings/${id}`, {
      method: "DELETE",
    });
  }

  // Helper method to set token manually
  setToken(token: string): void {
    this.token = token;
    if (typeof window !== "undefined") {
      localStorage.setItem("auth_token", token);
    }
  }

  // Helper method to get current token
  getToken(): string | null {
    return this.token;
  }
}

export const apiClient = new ApiClient();
