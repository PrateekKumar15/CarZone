// Featured cars data for the landing page
export interface FeaturedCar {
  id: number;
  brand: string;
  model: string;
  year: number;
  price: number;
  image: string;
  fuel_type: string;
  transmission: string;
  mileage: number;
  rating: number;
  reviews: number;
  features: string[];
}

export const featuredCars: FeaturedCar[] = [
  {
    id: 1,
    brand: "BMW",
    model: "X5",
    year: 2023,
    price: 65000,
    image: "/images/cars/bmw-x5.jpg",
    fuel_type: "Petrol",
    transmission: "Automatic",
    mileage: 15000,
    rating: 4.8,
    reviews: 127,
    features: [
      "Leather Seats",
      "Sunroof",
      "Navigation",
      "Bluetooth",
      "Backup Camera",
    ],
  },
  {
    id: 2,
    brand: "Mercedes",
    model: "C-Class",
    year: 2023,
    price: 45000,
    image: "/images/cars/mercedes-c-class.jpg",
    fuel_type: "Petrol",
    transmission: "Automatic",
    mileage: 12000,
    rating: 4.7,
    reviews: 89,
    features: [
      "Premium Audio",
      "Heated Seats",
      "Keyless Entry",
      "LED Headlights",
    ],
  },
];

// Car categories for filtering
export const carCategories = [
  { id: "all", name: "All Cars", icon: "" },
  { id: "luxury", name: "Luxury", icon: "" },
  { id: "sports", name: "Sports", icon: "" },
];
