"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import ProtectedRoute from "@/components/auth/protected-route";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { apiClient, Car } from "@/lib/api-client";
import { toast } from "sonner";
import {
  Search,
  Car as CarIcon,
  Fuel,
  Settings,
  ArrowLeft,
} from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import Image from "next/image";

// Custom image component with error handling
const CarImage = ({
  src,
  alt,
  className,
}: {
  src: string;
  alt: string;
  className?: string;
}) => {
  const [imageError, setImageError] = useState(false);
  const [imageLoading, setImageLoading] = useState(true);

  return (
    <>
      {!imageError ? (
        <Image
          src={src}
          alt={alt}
          fill
          className={`${className} ${
            imageLoading ? "opacity-0" : "opacity-100"
          } transition-opacity duration-200`}
          sizes="(max-width: 768px) 100vw, (max-width: 1024px) 50vw, 25vw"
          onLoad={() => setImageLoading(false)}
          onError={() => {
            setImageError(true);
            setImageLoading(false);
          }}
        />
      ) : (
        <div className="h-full flex items-center justify-center bg-gray-100">
          <CarIcon className="h-12 w-12 text-gray-400" />
        </div>
      )}
      {imageLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-100">
          <CarIcon className="h-8 w-8 text-gray-300 animate-pulse" />
        </div>
      )}
    </>
  );
};

export default function CarsPage() {
  const [cars, setCars] = useState<Car[]>([]);
  const [filteredCars, setFilteredCars] = useState<Car[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedFuelType, setSelectedFuelType] = useState<string>("all");
  const router = useRouter();

  useEffect(() => {
    const fetchCars = async () => {
      try {
        const data = await apiClient.getCars();
        setCars(data || []); // Handle null/undefined response
        setFilteredCars(data || []); // Handle null/undefined response
      } catch (error) {
        toast.error("Failed to fetch cars");
        console.error("Error fetching cars:", error);
        setCars([]); // Set empty array on error
        setFilteredCars([]); // Set empty array on error
      } finally {
        setLoading(false);
      }
    };

    fetchCars();
  }, []);

  useEffect(() => {
    let filtered = cars;

    // Filter by search term
    if (searchTerm) {
      filtered = filtered.filter(
        (car) =>
          car.brand.toLowerCase().includes(searchTerm.toLowerCase()) ||
          car.model.toLowerCase().includes(searchTerm.toLowerCase()) ||
          car.name.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Filter by fuel type
    if (selectedFuelType !== "all") {
      filtered = filtered.filter(
        (car) => car.fuel_type.toLowerCase() === selectedFuelType.toLowerCase()
      );
    }

    setFilteredCars(filtered);
  }, [searchTerm, selectedFuelType, cars]);

  const fuelTypes = ["all", ...new Set(cars.map((car) => car.fuel_type))];

  const formatPrice = (
    price: { rental_price_daily: number; sale_price?: number } | number
  ) => {
    const priceValue =
      typeof price === "number" ? price : price.rental_price_daily;
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 0,
    }).format(priceValue);
  };

  if (loading) {
    return (
      <ProtectedRoute>
        <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-8">
          <div className="max-w-7xl mx-auto">
            <div className="flex items-center gap-4 mb-8">
              <Skeleton className="h-10 w-10" />
              <Skeleton className="h-8 w-48" />
            </div>
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {[...Array(8)].map((_, i) => (
                <Card key={i} className="overflow-hidden">
                  <Skeleton className="h-48 w-full" />
                  <CardHeader>
                    <Skeleton className="h-6 w-3/4" />
                    <Skeleton className="h-4 w-1/2" />
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      <Skeleton className="h-4 w-full" />
                      <Skeleton className="h-4 w-full" />
                      <Skeleton className="h-10 w-full" />
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-8">
        <div className="max-w-7xl mx-auto">
          <div className="flex items-center gap-4 mb-8">
            <Button
              variant="outline"
              size="icon"
              onClick={() => router.push("/dashboard")}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Browse Cars</h1>
              <p className="text-gray-600">Find your perfect ride</p>
            </div>
          </div>

          {/* Filters */}
          <div className="bg-white rounded-lg shadow-sm p-6 mb-8">
            <div className="flex flex-col sm:flex-row gap-4">
              <div className="flex-1">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                  <Input
                    placeholder="Search by brand, model, or color..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </div>
              <div className="sm:w-48">
                <select
                  title="Filter by fuel type"
                  value={selectedFuelType}
                  onChange={(e) => setSelectedFuelType(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  {fuelTypes.map((fuel) => (
                    <option key={fuel} value={fuel}>
                      {fuel === "all" ? "All Fuel Types" : fuel}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <div className="mt-4 text-sm text-gray-600">
              Showing {filteredCars.length} of {cars.length} cars
            </div>
          </div>

          {filteredCars.length === 0 ? (
            <Card className="text-center py-12">
              <CardContent>
                <CarIcon className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  No cars found
                </h3>
                <p className="text-gray-600">
                  Try adjusting your search filters
                </p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {filteredCars.map((car) => (
                <Card
                  key={car.id}
                  className="overflow-hidden hover:shadow-lg transition-shadow group"
                >
                  <div className="relative h-48 bg-gray-200">
                    {car.images && car.images.length > 0 ? (
                      <CarImage
                        src={car.images[0]}
                        alt={`${car.brand} ${car.model}`}
                        className="object-cover group-hover:scale-105 transition-transform duration-300"
                      />
                    ) : (
                      <div className="h-full flex items-center justify-center">
                        <CarIcon className="h-12 w-12 text-gray-400" />
                      </div>
                    )}
                    <div className="absolute top-3 left-3">
                      <Badge className="bg-blue-600 text-white">
                        {car.fuel_type}
                      </Badge>
                    </div>
                    <div className="absolute top-3 right-3 bg-white rounded-full p-2 shadow-md">
                      <span className="text-xs font-semibold text-blue-600">
                        {car.year}
                      </span>
                    </div>
                  </div>

                  <CardHeader>
                    <CardTitle className="text-lg">
                      {car.brand} {car.model}
                    </CardTitle>
                    <CardDescription>
                      {car.mileage.toLocaleString()} miles • {car.year}
                    </CardDescription>
                  </CardHeader>

                  <CardContent className="space-y-4">
                    <div className="grid grid-cols-2 gap-4 text-sm text-gray-600">
                      <div className="flex items-center gap-1">
                        <Fuel className="h-4 w-4" />
                        <span>{car.fuel_type}</span>
                      </div>
                      <div className="flex items-center gap-1">
                        <Settings className="h-4 w-4" />
                        <span>{car.engine?.transmission || "N/A"}</span>
                      </div>
                    </div>

                    <div className="flex items-center justify-between pt-2 border-t">
                      <div>
                        <div className="text-2xl font-bold text-green-600">
                          {formatPrice(car.price)}
                        </div>
                        <div className="text-xs text-gray-500">per day</div>
                      </div>
                      <Button
                        onClick={() => router.push(`/cars/${car.id}`)}
                        className="bg-blue-600 hover:bg-blue-700"
                      >
                        View Details
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </ProtectedRoute>
  );
}
