"use client";

import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import { useAppSelector } from "@/hooks/redux";
import ProtectedRoute from "@/components/auth/protected-route";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import { Badge } from "@/components/ui/badge";
import { apiClient, Car } from "@/lib/api-client";
import { toast } from "sonner";
import {
  ArrowLeft,
  Car as CarIcon,
  Calendar,
  User,
  Phone,
  Mail,
  Info,
  MapPin,
  Star,
  Fuel,
  Settings,
  Gauge,
  Zap,
} from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import Image from "next/image";

// Custom Image Component with Error Handling
interface CarImageProps {
  src: string;
  alt: string;
  className?: string;
  fill?: boolean;
  sizes?: string;
}

const CarImage = ({
  src,
  alt,
  className = "",
  fill = false,
  sizes,
  ...props
}: CarImageProps) => {
  const [isLoading, setIsLoading] = useState(true);
  const [hasError, setHasError] = useState(false);

  const handleError = () => {
    setHasError(true);
    setIsLoading(false);
  };

  const handleLoad = () => {
    setIsLoading(false);
  };

  if (hasError || !src || src.includes("example.com")) {
    return (
      <div
        className={`flex items-center justify-center bg-gradient-to-br from-gray-100 to-gray-200 ${className}`}
      >
        <CarIcon className="h-16 w-16 text-gray-400" />
      </div>
    );
  }

  return (
    <div className={`relative overflow-hidden ${className}`}>
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-gradient-to-br from-gray-100 to-gray-200">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      )}
      <Image
        src={src}
        alt={alt}
        fill={fill}
        sizes={sizes}
        onError={handleError}
        onLoad={handleLoad}
        className={`object-cover transition-all duration-300 hover:scale-105 ${
          isLoading ? "opacity-0" : "opacity-100"
        }`}
        {...props}
      />
    </div>
  );
};

export default function CarDetailsPage() {
  const [car, setCar] = useState<Car | null>(null);
  const [loading, setLoading] = useState(true);
  const { user } = useAppSelector((state) => state.auth);
  const router = useRouter();
  const params = useParams();
  const carId = params.id as string;

  useEffect(() => {
    const fetchCar = async () => {
      try {
        const data = await apiClient.getCarById(carId);
        setCar(data);
      } catch (error) {
        toast.error("Failed to fetch car details");
        console.error("Error fetching car:", error);
        router.push("/cars");
      } finally {
        setLoading(false);
      }
    };

    if (carId) {
      fetchCar();
    }
  }, [carId, router]);

  const handleBooking = () => {
    if (!user?.id) {
      toast.error("Please login to book a car");
      return;
    }
    router.push(`/bookings/${carId}`);
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 0,
    }).format(price);
  };

  if (loading) {
    return (
      <ProtectedRoute>
        <div className="min-h-screen bg-gradient-to-br from-gray-50 to-white py-8">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <Skeleton className="h-96 w-full rounded-xl" />
              <div className="space-y-4">
                <Skeleton className="h-8 w-3/4" />
                <Skeleton className="h-4 w-1/2" />
                <Skeleton className="h-32 w-full" />
              </div>
            </div>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  if (!car) {
    return (
      <ProtectedRoute>
        <div className="min-h-screen bg-gradient-to-br from-gray-50 to-white py-8">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
            <h1 className="text-2xl font-bold text-gray-900">Car not found</h1>
            <Button onClick={() => router.push("/cars")} className="mt-4">
              Back to Cars
            </Button>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-gradient-to-br from-gray-50 to-white py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <div className="flex items-center justify-between mb-8">
            <Button
              variant="outline"
              size="sm"
              onClick={() => router.push("/cars")}
              className="flex items-center gap-2 hover:bg-gray-100 transition-colors"
            >
              <ArrowLeft className="h-4 w-4" />
              Back to Cars
            </Button>
            <div className="flex items-center gap-2">
              <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
              <span className="text-sm font-medium text-gray-600">4.8</span>
            </div>
          </div>

          <div className="grid grid-cols-1 xl:grid-cols-3 gap-8">
            {/* Image Carousel - Takes 2 columns on XL screens */}
            <div className="xl:col-span-2 space-y-4">
              {car.images && car.images.length > 0 ? (
                <Carousel className="w-full">
                  <CarouselContent>
                    {car.images.map((image, index) => (
                      <CarouselItem key={index}>
                        <div className="p-1">
                          <Card className="border-0 shadow-lg overflow-hidden">
                            <CardContent className="p-0">
                              <CarImage
                                src={image}
                                alt={`${car.brand} ${car.model} - Image ${
                                  index + 1
                                }`}
                                className="h-96 w-full rounded-xl"
                                fill
                                sizes="(max-width: 768px) 100vw, 66vw"
                              />
                            </CardContent>
                          </Card>
                        </div>
                      </CarouselItem>
                    ))}
                  </CarouselContent>
                  {car.images.length > 1 && (
                    <>
                      <CarouselPrevious className="left-4" />
                      <CarouselNext className="right-4" />
                    </>
                  )}
                </Carousel>
              ) : (
                <Card className="border-0 shadow-lg">
                  <CardContent className="flex items-center justify-center p-6 h-96">
                    <CarIcon className="h-32 w-32 text-gray-300" />
                  </CardContent>
                </Card>
              )}
            </div>

            {/* Car Details Sidebar */}
            <div className="space-y-6">
              {/* Car Info Card */}
              <Card className="border-0 shadow-lg">
                <CardContent className="p-6">
                  <div className="space-y-4">
                    <div className="flex items-start justify-between">
                      <Badge
                        variant={car.is_available ? "default" : "secondary"}
                        className={`${
                          car.is_available
                            ? "bg-green-100 text-green-800 hover:bg-green-200"
                            : "bg-red-100 text-red-800"
                        }`}
                      >
                        {car.is_available ? "Available" : "Unavailable"}
                      </Badge>
                    </div>

                    <div>
                      <h1 className="text-3xl font-bold text-gray-900 mb-2">
                        {car.name || `${car.brand} ${car.model}`}
                      </h1>
                      <div className="flex items-center text-gray-600 mb-4">
                        <MapPin className="h-4 w-4 mr-1" />
                        <span>
                          {car.year} • {car.location_city}, {car.location_state}
                        </span>
                      </div>
                    </div>

                    <div className="bg-gradient-to-r from-green-50 to-blue-50 p-4 rounded-lg">
                      <div className="text-3xl font-bold text-green-600 mb-1">
                        {formatPrice(car.price.rental_price_daily)}
                        <span className="text-lg font-normal text-gray-600">
                          /day
                        </span>
                      </div>
                      {car.price.sale_price && car.price.sale_price > 0 && (
                        <div className="text-sm text-gray-600">
                          Sale price: {formatPrice(car.price.sale_price)}
                        </div>
                      )}
                    </div>

                    {car.description && (
                      <div className="bg-gray-50 p-4 rounded-lg">
                        <p className="text-gray-700 text-sm leading-relaxed">
                          {car.description}
                        </p>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>

              {/* Quick Specs Card */}
              <Card className="border-0 shadow-lg">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2 text-lg">
                    <Settings className="h-5 w-5 text-blue-600" />
                    Quick Specs
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Fuel className="h-4 w-4 text-gray-500" />
                      <span className="text-sm text-gray-600">Fuel Type</span>
                    </div>
                    <span className="font-medium text-sm capitalize">
                      {car.fuel_type}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Zap className="h-4 w-4 text-gray-500" />
                      <span className="text-sm text-gray-600">
                        Transmission
                      </span>
                    </div>
                    <span className="font-medium text-sm">
                      {car.engine?.transmission || "N/A"}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Gauge className="h-4 w-4 text-gray-500" />
                      <span className="text-sm text-gray-600">Mileage</span>
                    </div>
                    <span className="font-medium text-sm">
                      {car.mileage.toLocaleString()} mi
                    </span>
                  </div>
                </CardContent>
              </Card>

              {/* Owner Information */}
              {car.owner && (
                <Card className="border-0 shadow-lg">
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2 text-lg">
                      <User className="h-5 w-5 text-blue-600" />
                      Car Owner
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                        <div className="h-10 w-10 bg-blue-100 rounded-full flex items-center justify-center">
                          <User className="h-5 w-5 text-blue-600" />
                        </div>
                        <div>
                          <p className="font-medium text-gray-900">
                            {car.owner.username}
                          </p>
                          <p className="text-sm text-gray-600">Owner</p>
                        </div>
                      </div>
                      <div className="flex items-center gap-3">
                        <Mail className="h-4 w-4 text-gray-500" />
                        <p className="text-sm text-gray-700">
                          {car.owner.email}
                        </p>
                      </div>
                      {car.owner.phone && (
                        <div className="flex items-center gap-3">
                          <Phone className="h-4 w-4 text-gray-500" />
                          <p className="text-sm text-gray-700">
                            {car.owner.phone}
                          </p>
                        </div>
                      )}
                    </div>
                  </CardContent>
                </Card>
              )}

              {/* Book Button */}
              {car.is_available && (
                <Button
                  onClick={handleBooking}
                  size="lg"
                  className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-semibold py-4 text-lg shadow-lg hover:shadow-xl transition-all duration-200"
                >
                  <Calendar className="h-5 w-5 mr-2" />
                  Book This Car
                </Button>
              )}
            </div>
          </div>

          {/* Detailed Specifications - Full Width */}
          <div className="mt-8">
            <Card className="border-0 shadow-lg">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-xl">
                  <Info className="h-6 w-6 text-blue-600" />
                  Detailed Specifications
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {/* Engine Specs */}
                  <div className="space-y-3">
                    <h4 className="font-semibold text-gray-900 flex items-center gap-2">
                      <Zap className="h-4 w-4 text-blue-600" />
                      Engine
                    </h4>
                    <div className="space-y-2 pl-6">
                      <div className="flex justify-between">
                        <span className="text-gray-600 text-sm">Size:</span>
                        <span className="font-medium text-sm">
                          {car.engine?.engine_size || "N/A"}L
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 text-sm">Power:</span>
                        <span className="font-medium text-sm">
                          {car.engine?.horsepower || "N/A"} HP
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 text-sm">Torque:</span>
                        <span className="font-medium text-sm">
                          {car.engine?.torque || "N/A"} lb-ft
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* General Specs */}
                  <div className="space-y-3">
                    <h4 className="font-semibold text-gray-900 flex items-center gap-2">
                      <Settings className="h-4 w-4 text-blue-600" />
                      General
                    </h4>
                    <div className="space-y-2 pl-6">
                      <div className="flex justify-between">
                        <span className="text-gray-600 text-sm">Status:</span>
                        <span className="font-medium text-sm capitalize">
                          {car.status}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 text-sm">Color:</span>
                        <span className="font-medium text-sm">
                          {car.color || "N/A"}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 text-sm">VIN:</span>
                        <span className="font-medium text-xs">
                          {car.vin || "N/A"}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* Pricing */}
                  <div className="space-y-3">
                    <h4 className="font-semibold text-gray-900 flex items-center gap-2">
                      <span className="h-4 w-4 text-green-600">$</span>
                      Pricing
                    </h4>
                    <div className="space-y-2 pl-6">
                      <div className="flex justify-between">
                        <span className="text-gray-600 text-sm">Daily:</span>
                        <span className="font-medium text-sm text-green-600">
                          {formatPrice(car.price.rental_price_daily)}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 text-sm">Weekly:</span>
                        <span className="font-medium text-sm text-green-600">
                          {formatPrice(car.price.rental_price_weekly)}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600 text-sm">Monthly:</span>
                        <span className="font-medium text-sm text-green-600">
                          {formatPrice(car.price.rental_price_monthly)}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}
